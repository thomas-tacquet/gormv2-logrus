package gormv2_logrus

import (
	"context"
	"time"
)

// gormlog must match gorm logger.Interface to be compatible with gorm.
// gormlog can be assigned in gorm configuration (see example in README.md)
type gormlog struct {
	opts options
}

// New create an instance of
func New(opts ...Option) *gormlog {
	gl := &gormlog{}

	for _, opt := range opts {
		opt.apply(&gl.opts)
	}

	return &gormlog{}
}

// Info implementaiton of info log level
func (gl *gormlog) Info(ctx context.Context, msg string, args ...interface{}) {
	gl.opts.logrus.WithContext(ctx).Infof(msg, args)
}

// Warn implementaiton of warn log level
func (gl *gormlog) Warn(ctx context.Context, msg string, args ...interface{}) {
	gl.opts.logrus.WithContext(ctx).Warnf(msg, args)
}

// Error gormlog of error log level
func (gl *gormlog) Error(ctx context.Context, msg string, args ...interface{}) {
	gl.opts.logrus.WithContext(ctx).Errorf(msg, args)
}

// Trace implementaiton of trace log level
func (gl *gormlog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	gl.opts.logrus.WithContext(ctx).Trace("trace test TODO")
}
