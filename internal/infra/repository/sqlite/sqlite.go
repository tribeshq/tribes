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

	if err := db.AutoMigrate(
		&entity.Campaign{},
		&entity.Order{},
		&entity.User{},
		&entity.SocialAccount{},
	); err != nil {
		return nil, err
	}

	configs.SetDefaults()

	isMemory := dbPath == ":memory:"
	var (
		adminAddr, verifierAddr, deployerAddr custom_type.Address
	)

	if isMemory {
		if a, err := configs.GetTribesAdminAddressTest(); err != nil {
			return nil, fmt.Errorf("failed to get TRIBES_ADMIN_ADDRESS_TEST: %w", err)
		} else {
			adminAddr = custom_type.HexToAddress(a.Hex())
		}

		if v, err := configs.GetTribesVerifierAddressTest(); err != nil {
			return nil, fmt.Errorf("failed to get TRIBES_VERIFIER_ADDRESS_TEST: %w", err)
		} else {
			verifierAddr = custom_type.HexToAddress(v.Hex())
		}

		if d, err := configs.GetTribesDeployerAddressTest(); err != nil {
			return nil, fmt.Errorf("failed to get TRIBES_DEPLOYER_ADDRESS_TEST: %w", err)
		} else {
			deployerAddr = custom_type.HexToAddress(d.Hex())
		}
	} else {
		if a, err := configs.GetTribesAdminAddress(); err != nil {
			return nil, fmt.Errorf("failed to get TRIBES_ADMIN_ADDRESS: %w", err)
		} else {
			adminAddr = custom_type.HexToAddress(a.Hex())
		}

		if v, err := configs.GetTribesVerifierAddress(); err != nil {
			return nil, fmt.Errorf("failed to get TRIBES_VERIFIER_ADDRESS: %w", err)
		} else {
			verifierAddr = custom_type.HexToAddress(v.Hex())
		}

		if d, err := configs.GetTribesDeployerAddress(); err != nil {
			return nil, fmt.Errorf("failed to get TRIBES_DEPLOYER_ADDRESS: %w", err)
		} else {
			deployerAddr = custom_type.HexToAddress(d.Hex())
		}
	}

	now := time.Now().Unix()
	users := []entity.User{
		{
			Role:      entity.UserRoleAdmin,
			Address:   adminAddr,
			CreatedAt: now,
		},
		{
			Role:      entity.UserRoleVerifier,
			Address:   verifierAddr,
			CreatedAt: now,
		},
		{
			Role:      entity.UserRoleDeployer,
			Address:   deployerAddr,
			CreatedAt: now,
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user %v: %w", user.Role, err)
		}
	}

	return &SQLiteRepository{Db: db}, nil
}
