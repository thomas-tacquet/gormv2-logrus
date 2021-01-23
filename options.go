package gormv2_logrus

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// withLogrus obviously
// option to truncate loooong requests
// option do hide requests containing passwords etc...
// option log request latency
// option to change write format

// options represents all available starting options for gormv2_logrus.
// options are optional parameters that can be passed to init function of gormv2_logrus
// to configure log policy.
type options struct {
	// pointer to your logrus instance
	logrus *logrus.Entry

	// if set tO 0, nothing wil be truncated, else you can set it to the value you want to avoid
	// logging too big SQL queries
	truncateLen uint

	// if a query contains one of bannedKeywords, it will not be logged, it's useful for preventing passwords and secrets
	// for being logged.
	bannedKeywords []string

	// if set to true, it will add latency informations for your queries
	logLatency bool

	SlowThreshold time.Duration

	Colorful bool

	LogLevel logger.LogLevel
}

func defaultOptions() options {
	return options{
		logrus:         nil,
		truncateLen:    0,
		bannedKeywords: nil,
		logLatency:     false,
	}
}

// Option interface is used to configure gormv2_logrus options.
type Option interface {
	apply(*options)
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

// apply is used in mstack init function to read given parameters.
func (fo *funcOption) apply(do *options) {
	fo.f(do)
}

// newGormLogOption is implemented by function that save parameters.
func newGormLogOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithNoWarden Option to specify your logrus instance. If you don't set a logrus isntance
// or if your logrusInstance is nil, gormlog will consider that you want log to be printed on stdout.
// It's useful on developpement purposes when you want to see your logs directly in terminal.
func WithLogrus(logrusInstance *logrus.Entry) Option {
	return newGormLogOption(func(o *options) {
		o.logrus = logrusInstance
	})
}
