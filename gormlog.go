package gormv2_logrus

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
