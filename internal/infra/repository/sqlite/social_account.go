package sqlite

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateSocialAccount(input *entity.SocialAccount) (*entity.SocialAccount, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to verify social account: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindSocialAccountById(id uint) (*entity.SocialAccount, error) {
	var account entity.SocialAccount
	if err := r.Db.First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("social account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find social account by ID: %w", err)
	}
	return &account, nil
}

func (r *SQLiteRepository) FindSocialAccountsByUserId(userID uint) ([]*entity.SocialAccount, error) {
	var accounts []*entity.SocialAccount
	if err := r.Db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to find social accounts by user ID: %w", err)
	}
	return accounts, nil
}

func (r *SQLiteRepository) DeleteSocialAccount(id uint) error {
	res := r.Db.Delete(&entity.SocialAccount{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete social account: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("social account not found")
	}
	return nil
}
