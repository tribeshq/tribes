package sqlite

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	. "github.com/tribeshq/tribes/pkg/custom_type"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	// Remove sqlite:// prefix if present
	dbPath := strings.TrimPrefix(conn, "sqlite://")

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

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
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

	adminUser := entity.User{
		Role:              entity.UserRoleAdmin,
		Address:           HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		InvestmentLimit:   uint256.NewInt(0),
		DebtIssuanceLimit: uint256.NewInt(0),
		CreatedAt:         time.Now().Unix(),
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &SQLiteRepository{
		Db: db,
	}, nil
}
