package sqlite

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	. "github.com/tribeshq/tribes/pkg/custom_type"

	"github.com/holiman/uint256"
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

func NewSQLiteRepository(conn string) (*SQLiteRepository, error) {
	dbPath := strings.TrimPrefix(conn, "sqlite://")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&entity.Auction{},
		&entity.Order{},
		&entity.User{},
		&entity.SocialAccount{},
	)
	if err != nil {
		return nil, err
	}

	adminUser := entity.User{
		Role:              entity.UserRoleAdmin,
		Address:           HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9"),
		InvestmentLimit:   uint256.NewInt(0),
		CreatedAt:         time.Now().Unix(),
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &SQLiteRepository{Db: db}, nil
}
