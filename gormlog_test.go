package gormv2logrus_test

import (
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	gormv2logrus "github.com/thomas-tacquet/gormv2-logrus"
)

func TestWithLogrus(t *testing.T) {
	// create a logrusEntry entry to giveit to gormv2_logrus
	logger, hook := test.NewNullLogger()

	// create the gorm compatible logger with logrusEntry instance
	gormLog := gormv2logrus.NewGormlog(gormv2logrus.WithLogrus(logger))

	//
	// open in memory database with previous logger
	db, err := gorm.Open(sqlite.Open(
		"file:unit_test_01?mode=memory&cache=shared"),
		&gorm.Config{Logger: gormLog},
	)

	// check if database correctly created
	require.NoError(t, err)
	require.NotNil(t, db)

	sqlDB, err := db.DB()

	require.NoError(t, err)
	require.NotNil(t, sqlDB)

	defer func() {
		assert.NoError(t, sqlDB.Close())
	}()

	type Placeholder struct{}

	_ = db.Create(&Placeholder{})

	assert.Equal(t, 1, len(hook.Entries))
}
