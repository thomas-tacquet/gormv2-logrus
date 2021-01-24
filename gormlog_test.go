package gormv2_logrus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWithLogrus(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:unit_test_01?mode=memory&cache=shared"), &gorm.Config{})

	require.NoError(t, err)
	require.NotNil(t, db)

	sqlDB, err := db.DB()

	require.NoError(t, err)
	require.NotNil(t, sqlDB)

	defer func() {
		assert.NoError(t, sqlDB.Close())
	}()
}
