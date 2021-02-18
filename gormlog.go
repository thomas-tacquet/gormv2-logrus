package gormv2logrus

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// gormlog must match gorm logger.Interface to be compatible with gorm.
// gormlog can be assigned in gorm configuration (see example in README.md)
type gormlog struct {
	// SkipErrRecordNotFound if set to true, errors of type gorm.ErrRecordNotFound will be ignored.
	SkipErrRecordNotFound bool

	// SlowThreshold is used to determine a limit of slow requests, if a request time is above SlowThreshold,
	// it will be logged as warning.
	SlowThreshold time.Duration

	// SourceField if definied, source will appear in log with detailled file context.
	SourceField string

	opts options
}

// NewGormlog create an instance of
func NewGormlog(opts ...Option) *gormlog {
	gl := &gormlog{
		opts: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&gl.opts)
	}

	return gl
}

// LogMod implementation log mode
func (gl *gormlog) LogMode(logger.LogLevel) logger.Interface {
	return gl
}

// Info implementaiton of info log level
func (gl *gormlog) Info(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Infof(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Infof(msg, args...)
	}
}

// Warn implementaiton of warn log level
func (gl *gormlog) Warn(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Warnf(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Warnf(msg, args...)
	}
}

// Error gormlog of error log level
func (gl *gormlog) Error(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Errorf(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Errorf(msg, args...)
	}
}

// Trace implementaiton of trace log level
func (gl *gormlog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// retrive sql string
	traceLog, _ := fc()

	// use begin to compute performances
	stopWatch := time.Since(begin)

	// additional logrus fields
	logrusFields := logrus.Fields{}

	// if source field is definied, we retreive line number informations
	if len(gl.SourceField) > 0 {
		logrusFields[gl.SourceField] = utils.FileWithLineNum()
	}

	// check if we have an error
	if err != nil {
		if !(errors.Is(err, gorm.ErrRecordNotFound) && gl.SkipErrRecordNotFound) {
			logrusFields[logrus.ErrorKey] = err

			if gl.opts.lr != nil {
				gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Errorf("%s [%s]", traceLog, stopWatch)
			}

			if gl.opts.logrusEntry != nil {
				gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Errorf("%s [%s]", traceLog, stopWatch)
			}

			return
		}
	}

	if gl.SlowThreshold != 0 && stopWatch > gl.SlowThreshold {
		if gl.opts.lr != nil {
			gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Warnf("%s [%s]", traceLog, stopWatch)
		}

		if gl.opts.logrusEntry != nil {
			gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Warnf("%s [%s]", traceLog, stopWatch)
		}

		return
	}

	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Debugf("%s [%s]", traceLog, stopWatch)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Debugf("%s [%s]", traceLog, stopWatch)
	}
}
