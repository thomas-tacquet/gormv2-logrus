package gormv2_logrus

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
)

// gormlog must match gorm logger.Interface to be compatible with gorm.
// gormlog can be assigned in gorm configuration (see example in README.md)
type gormlog struct {
	SkipErrRecordNotFound bool
	SlowThreshold         time.Duration
	SourceField           string

	opts options
}

// New create an instance of
func New(opts ...Option) *gormlog {
	gl := &gormlog{}

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
	gl.opts.logrusEntry.WithContext(ctx).Warnf(msg, args...)
}

// Error gormlog of error log level
func (gl *gormlog) Error(ctx context.Context, msg string, args ...interface{}) {
	gl.opts.logrusEntry.WithContext(ctx).Errorf(msg, args...)
}

// Trace implementaiton of trace log level
func (gl *gormlog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.opts.lr != nil {
		gl.opts.lr.Error("trace test TODO")
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Trace("trace test TODO")
	}
}
