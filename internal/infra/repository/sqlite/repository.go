package sqlite

import (
	"context"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tribeshq/tribes/internal/domain/entity"
)

type SQLiteRepository struct {
	Db *gorm.DB
}

func (r *SQLiteRepository) Close() error {
	sqlDB, err := r.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

func NewSQLiteRepository(ctx context.Context, conn string) (*SQLiteRepository, error) {
	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             0,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open SQLite database
	db, err := gorm.Open(sqlite.Open(conn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate schema
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Order{},
		&entity.Contract{},
		&entity.Crowdfunding{},
		&entity.SocialAccount{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return &SQLiteRepository{
		Db: db,
	}, nil
}
