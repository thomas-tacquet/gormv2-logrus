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

// Gormlog implements the gorm logger.Interface, making it compatible with GORM v2.
// It can be assigned in the GORM configuration (see example in README.md).
//
// This logger uses Logrus for logging and supports configurable levels, slow query thresholds,
// and optional source code location tracking.
type Gormlog struct {
	// SkipErrRecordNotFound, when set to true, suppresses logging of gorm.ErrRecordNotFound errors.
	SkipErrRecordNotFound bool

	// SlowThreshold defines the duration beyond which a query is considered "slow"
	// and will be logged as a warning.
	SlowThreshold time.Duration

	// SourceField, if specified, enables logging of the file and line number
	// where the log entry originated. The value will be used as the field name in the log.
	SourceField string

	// LogLevel sets the minimum log level for messages to be logged.
	LogLevel logger.LogLevel

	opts options
}

// NewGormlog creates and returns a new instance of Gormlog with the provided options.
// If no options are given, default values are used.
func NewGormlog(opts ...Option) *Gormlog {
	gl := &Gormlog{
		opts: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&gl.opts)
	}

	return gl
}

// LogMode sets the logging level for the Gormlog instance.
// It returns a new logger interface with the specified log level.
func (gl *Gormlog) LogMode(ll logger.LogLevel) logger.Interface {
	gl.LogLevel = ll

	return gl
}

// Info logs an informational message using the configured Logrus logger.
// It accepts a context, a format string, and optional arguments.
func (gl *Gormlog) Info(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Infof(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Infof(msg, args...)
	}
}

// Warn logs a warning message using the configured Logrus logger.
// It accepts a context, a format string, and optional arguments.
func (gl *Gormlog) Warn(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Warnf(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Warnf(msg, args...)
	}
}

// Error logs an error message using the configured Logrus logger.
// It accepts a context, a format string, and optional arguments.
func (gl *Gormlog) Error(ctx context.Context, msg string, args ...interface{}) {
	if gl.opts.lr != nil {
		gl.opts.lr.WithContext(ctx).Errorf(msg, args...)
	}

	if gl.opts.logrusEntry != nil {
		gl.opts.logrusEntry.WithContext(ctx).Errorf(msg, args...)
	}
}

// Trace implementation of trace log level.
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
	for i := range int(len(gl.opts.bannedKeywords)) {
		keyword := gl.opts.bannedKeywords[i]
		if keyword.CaseMatters && strings.Contains(traceLog, keyword.Keyword) {
			return
		} else if !keyword.CaseMatters &&
			strings.Contains(strings.ToLower(traceLog), strings.ToLower(keyword.Keyword)) {
			return
		}
	}

	// check if we have an error
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && gl.SkipErrRecordNotFound) {
		logrusFields[logrus.ErrorKey] = err

		if gl.opts.lr != nil {
			gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Errorf("%s", traceLog)
		}

		if gl.opts.logrusEntry != nil {
			gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Errorf("%s", traceLog)
		}
	}

	if gl.opts.SlowThreshold != 0 && stopWatch > gl.opts.SlowThreshold && gl.opts.LogLevel >= logger.Warn {
		// instead of adding SLOW SQL to the message, add reason field
		// this can be parsed easily with logs management tools
		logrusFields["reason"] = "SLOW SQL"

		if gl.opts.lr != nil {
			gl.opts.lr.WithContext(ctx).WithFields(logrusFields).Warnf("%s", traceLog)
		}

		if gl.opts.logrusEntry != nil {
			gl.opts.logrusEntry.WithContext(ctx).WithFields(logrusFields).Warnf("%s", traceLog)
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
