package gormv2_logrus


// gormlog must match gorm logger.Interface to be compatible with gorm.
// gormlog can be assigned in gorm configuration (see example in README.md)
type gormlog struct {

}

// New create an instance of
func New() *gormlog {
	return &gormlog{}
}


