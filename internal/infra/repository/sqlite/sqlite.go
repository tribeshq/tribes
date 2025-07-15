package sqlite

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
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
	dbPath := strings.TrimPrefix(conn, "sqlite://")
	isMemory := dbPath == ":memory:"

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
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Add context to DB
	db = db.WithContext(ctx)

	// Optional: check DB connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite: %w", err)
	}

	if err := db.AutoMigrate(
		&entity.Campaign{},
		&entity.Order{},
		&entity.User{},
		&entity.SocialAccount{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	configs.SetDefaults()

	var adminAddr, verifierAddr custom_type.Address

	if isMemory {
		a, err := configs.GetAdminAddressTest()
		if err != nil && err != configs.ErrNotDefined {
			return nil, fmt.Errorf("failed to get TRIBES_ADMIN_ADDRESS_TEST: %w", err)
		} else if err == configs.ErrNotDefined {
			return nil, fmt.Errorf("TRIBES_ADMIN_ADDRESS_TEST is required for the rollup service: %w", err)
		}
		adminAddr = custom_type.HexToAddress(a.Hex())

		v, err := configs.GetVerifierAddressTest()
		if err != nil && err != configs.ErrNotDefined {
			return nil, fmt.Errorf("failed to get TRIBES_VERIFIER_ADDRESS_TEST: %w", err)
		} else if err == configs.ErrNotDefined {
			return nil, fmt.Errorf("TRIBES_VERIFIER_ADDRESS_TEST is required for the rollup service: %w", err)
		}
		verifierAddr = custom_type.HexToAddress(v.Hex())
	} else {
		a, err := configs.GetAdminAddress()
		if err != nil && err != configs.ErrNotDefined {
			return nil, fmt.Errorf("failed to get TRIBES_ADMIN_ADDRESS: %w", err)
		} else if err == configs.ErrNotDefined {
			return nil, fmt.Errorf("TRIBES_ADMIN_ADDRESS is required for the rollup service: %w", err)
		}
		adminAddr = custom_type.HexToAddress(a.Hex())

		v, err := configs.GetVerifierAddress()
		if err != nil && err != configs.ErrNotDefined {
			return nil, fmt.Errorf("failed to get TRIBES_VERIFIER_ADDRESS: %w", err)
		} else if err == configs.ErrNotDefined {
			return nil, fmt.Errorf("TRIBES_VERIFIER_ADDRESS is required for the rollup service: %w", err)
		}
		verifierAddr = custom_type.HexToAddress(v.Hex())
	}

	baseTime := time.Now().Unix()
	users := []entity.User{
		{
			Role:      entity.UserRoleAdmin,
			Address:   adminAddr,
			CreatedAt: baseTime,
		},
		{
			Role:      entity.UserRoleVerifier,
			Address:   verifierAddr,
			CreatedAt: baseTime,
		},
	}

	for _, user := range users {
		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user %v: %w", user.Role, err)
		}
	}

	return &SQLiteRepository{Db: db}, nil
}