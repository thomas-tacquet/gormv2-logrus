package gormv2logrus_test

import (
	"fmt"
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

	db, err := gorm.Open(sqlite.Open(
		generateTestingSqliteString(t)),
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

	// NotExistingTable is a simple empty struct that does not exist in current database,
	// so if we try to create a new entry of this struct, gorm must return an error
	// telling us that this table does not exists
	type NotExistingTable struct{}

	errCreate := db.Create(&NotExistingTable{}).Error
	t.Log(errCreate.Error())

	// testing gorm is not a purprose of this test, but to ensure consistency we
	// must check if errCreate is not empty
	require.NotEmpty(t, errCreate)
	require.Contains(t, errCreate.Error(), "no such table")
	require.Contains(t, errCreate.Error(), "not_existing_tables")

	assert.Equal(t, 1, len(hook.Entries))
	require.NotNil(t, hook.LastEntry())

	lastLogEntry, err := hook.LastEntry().String()
	require.NoError(t, err)

	assert.Contains(t, lastLogEntry, errCreate.Error())
}

func generateTestingSqliteString(t *testing.T) string {
	t.Helper()

	const sqliteConnString = "file:%s?mode=memory&cache=shared"

	return fmt.Sprintf(sqliteConnString, t.Name())
}
