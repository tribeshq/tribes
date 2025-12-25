package sqlite

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateUser(input *entity.User) (*entity.User, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindUserByAddress(address types.Address) (*entity.User, error) {
	var user entity.User
	if err := r.Db.Preload("SocialAccounts").Where("address = ?", address).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by address: %w", err)
	}
	return &user, nil
}

func (r *SQLiteRepository) FindUsersByRole(role string) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.Db.Preload("SocialAccounts").Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}
	return users, nil
}

func (r *SQLiteRepository) FindAllUsers() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.Db.Preload("SocialAccounts").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}
	return users, nil
}

func (r *SQLiteRepository) DeleteUser(address types.Address) error {
	res := r.Db.Where("address = ?", address).Delete(&entity.User{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}
	return nil
}
