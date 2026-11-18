package store

import (
	"fmt"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open opens (creating if needed) the SQLite database, applies the embedded
// goose migrations, and returns the shared *gorm.DB handle.
//
// glebarez/sqlite is a pure-Go driver (modernc.org/sqlite under the hood), so
// CGO is not required and cross-compiled binaries stay static.
func Open(path string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", path)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access sql db: %w", err)
	}
	// SQLite allows only one writer at a time. Serialising connections avoids
	// SQLITE_BUSY errors under concurrent writes.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	if err := migrate(sqlDB); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return db, nil
}
