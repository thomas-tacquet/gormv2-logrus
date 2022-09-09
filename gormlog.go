package gormv2logrus

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Gormlog must match gorm logger.Interface to be compatible with gorm.
// Gormlog can be assigned in gorm configuration (see example in README.md).
type Gormlog struct {
	// SkipErrRecordNotFound if set to true, errors of type gorm.ErrRecordNotFound will be ignored.
	SkipErrRecordNotFound bool

	// SlowThreshold is used to determine a limit of slow requests, if a request time is above SlowThreshold,
	// it will be logged as warning.
	SlowThreshold time.Duration

	// SourceField if definied, source will appear in log with detailed file context.
	SourceField string

	LogLevel logger.LogLevel

	opts options
}

// NewGormlog create an instance of.
func NewGormlog(opts ...Option) *Gormlog {

	gl := &Gormlog{
		opts: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&gl.opts)
	}

	return gl
}

// LogMode implementation log mode.
func (gl *Gormlog) LogMode(ll logger.LogLevel) logger.Interface {
	gl.LogLevel = ll

	return gl
}

// Info implementation of info log level
func (gl *Gormlog) Info(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Infof(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Infof(msg, args...)
	}
}

// Warn implementation of warn log level
func (gl *Gormlog) Warn(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Warnf(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Warnf(msg, args...)
	}
}

// Error Gormlog of error log level
func (gl *Gormlog) Error(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Errorf(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Errorf(msg, args...)
	}
}

// Trace implementation of trace log level
func (gl *Gormlog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	// retrieve sql string and affected rows
	traceLog, rows := fc()

	// use begin to compute performances
	stopWatch := time.Since(begin)

	// additional logrus fields
	logrusFields := logrus.Fields{}

	// if logLatency is true, add stopWatch information
	if gl.opts.logLatency {
		logrusFields["duration"] = stopWatch
	}

	// add number of affected rows as logrus parameter
	logrusFields["rows"] = rows

	// if source field is definied, we retrieve line number information
	if len(gl.SourceField) > 0 {
		logrusFields[gl.SourceField] = utils.FileWithLineNum()
	}

	// scanning for banned keywords
	for idx := 0; idx < len(gl.opts.bannedKeywords); idx++ {
		if gl.opts.bannedKeywords[idx].CaseMatters &&
			strings.Contains(traceLog, gl.opts.bannedKeywords[idx].Keyword) {
			return
		} else if !gl.opts.bannedKeywords[idx].CaseMatters &&
			strings.Contains(
				strings.ToLower(traceLog),
				strings.ToLower(gl.opts.bannedKeywords[idx].Keyword),
			) {
			return
		}
	}

	// check if we have an error
	if err != nil {
		if !(errors.Is(err, gorm.ErrRecordNotFound) && gl.SkipErrRecordNotFound) {
			logrusFields[logrus.ErrorKey] = err

			if gl.opts.lr != nil {
				gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Errorf("%s", traceLog)
			}

			if gl.opts.logrusEntry != nil {
				gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Errorf("%s", traceLog)
			}

			return
		}
	}

	if gl.SlowThreshold != 0 && stopWatch > gl.SlowThreshold {
		if gl.opts.lr != nil {
			gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Warnf("SLOW SQL %s", traceLog)
		}

		if gl.opts.logrusEntry != nil {
			gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Warnf("SLOW SQL %s", traceLog)
		}

		return
	}

	// Use directly with logrus entry
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Debugf("%s", traceLog)
	}

	// Use with logrusEntry
	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Debugf("%s", traceLog)
	}
}
