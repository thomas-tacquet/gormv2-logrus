package gormv2logrus

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// GormV2Logrus is a GORM v2 logger implementation that uses Logrus for logging.
// It provides options to:
// - Avoid logging sensitive data (e.g., passwords) using banned keywords
// - Truncate long SQL queries
// - Log query execution latency
// - Customize log output format via Logrus
// - Set log levels and slow query thresholds

// Options represents all available configuration options for GormV2Logrus.
// These are optional parameters passed to the logger initialization
// to define logging behavior.
type options struct {
	// bannedKeywords contains a list of keywords that, if found in a query,
	// will prevent the query from being logged. This helps avoid logging
	// sensitive information such as passwords or tokens.
	bannedKeywords []BannedKeyword

	// logrusEntry is a pointer to a Logrus entry instance for structured logging.
	// If not set, logs will be written to stdout using a default Logrus logger.
	logrusEntry *logrus.Entry

	// lr is a pointer to a Logrus logger instance.
	// Used when no entry is provided, or to override the default logger.
	lr *logrus.Logger

	// SlowThreshold defines the duration after which a query is considered "slow".
	// Slow queries can be logged at a different level (e.g., Warn).
	SlowThreshold time.Duration

	// LogLevel sets the verbosity of the logs (e.g., Silent, Error, Warn, Info).
	LogLevel logger.LogLevel

	// truncateLen specifies the maximum length of SQL queries to log.
	// If set to 0, no truncation is performed.
	// Useful to avoid logging excessively long queries.
	truncateLen uint

	// logLatency, when true, includes execution time information in the logs.
	logLatency bool
}

// BannedKeyword defines a rule for filtering out log entries that contain sensitive data.
type BannedKeyword struct {
	// Keyword is the string to search for in the log output (e.g., "password").
	Keyword string

	// CaseMatters determines whether the keyword matching is case-sensitive.
	// If false, matching is case-insensitive.
	CaseMatters bool
}

// GormOptions is a public configuration struct used to customize the behavior of the GormV2Logrus logger.
// It allows setting common logging options such as slow query thresholds, log level, query truncation,
// and whether to include execution latency in logs.
//
// Example usage:
//
//	opts := GormOptions{
//	    SlowThreshold: 200 * time.Millisecond,
//	    LogLevel:      logger.Info,
//	    TruncateLen:   1000,
//	    LogLatency:    true,
//	}
//	logger := New(opts)
type GormOptions struct {
	// LogLevel sets the verbosity level for GORM logs (e.g., Silent, Error, Warn, Info).
	LogLevel logger.LogLevel

	// TruncateLen sets the maximum number of characters to log for SQL queries.
	// If 0, queries are not truncated.
	TruncateLen uint

	// LogLatency controls whether query execution time is included in the log output.
	LogLatency bool

	// SlowThreshold defines the duration beyond which a query is considered slow.
	// Slow queries may be logged at a higher log level (e.g., Warn).
	SlowThreshold time.Duration
}

// defaultOptions returns a new options struct with default values.
// These defaults are used when initializing the logger without custom configuration.
// The returned options include:
// - No Logrus entry or logger (uses default Logrus instance)
// - No banned keywords
// - No truncation (truncateLen = 0)
// - Latency logging disabled
// - SlowThreshold and LogLevel must be explicitly set
func defaultOptions() options {
	return options{
		logrusEntry:    nil,
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

// newGormLogOption creates a new Option that applies the given configuration function.
func newGormLogOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithLogrusEntry sets a custom Logrus Entry for logging.
// This option is not compatible with WithLogrus.
// If no Logrus Entry or Logger is provided, logs will be output to stdout using the default Logrus logger.
// This is particularly useful during development when you want to see logs directly in the terminal.
func WithLogrusEntry(logrusEntry *logrus.Entry) Option {
	return newGormLogOption(func(o *options) {
		o.logrusEntry = logrusEntry
	})
}

// WithLogrus sets a custom Logrus Logger instance for logging.
// This option is not compatible with WithLogrusEntry.
// Use this when you want to use a specific Logrus logger configuration (e.g., custom output, formatter, or level).
func WithLogrus(lr *logrus.Logger) Option {
	return newGormLogOption(func(o *options) {
		o.lr = lr
	})
}

// WithBannedKeyword configures a list of keywords that, if present in a SQL query or log message,
// will cause the log entry to be suppressed. This helps prevent sensitive data like passwords
// or tokens from being inadvertently logged.
// Matching behavior depends on the CaseMatters field of each BannedKeyword.
func WithBannedKeyword(bannedKeywords []BannedKeyword) Option {
	return newGormLogOption(func(o *options) {
		o.bannedKeywords = bannedKeywords
	})
}

// WithGormOptions applies a set of common logging options from a GormOptions struct.
// This includes settings for log level, query truncation, slow query threshold,
// and whether to include execution latency in logs.
// It provides a convenient way to configure multiple options at once.
func WithGormOptions(gormOpt GormOptions) Option {
	return newGormLogOption(func(o *options) {
		o.logLatency = gormOpt.LogLatency
		o.LogLevel = gormOpt.LogLevel
		o.SlowThreshold = gormOpt.SlowThreshold
		o.truncateLen = gormOpt.TruncateLen
	})
}
