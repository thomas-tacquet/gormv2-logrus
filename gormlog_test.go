package gormv2_logrus_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWithLogrus(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
}
