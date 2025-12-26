package sqlite

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateIssuance(input *entity.Issuance) (*entity.Issuance, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create issuance: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindIssuanceById(id uint) (*entity.Issuance, error) {
	var Issuance entity.Issuance
	if err := r.Db.
		Preload("Orders").
		First(&Issuance, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrIssuanceNotFound
		}
		return nil, fmt.Errorf("failed to find issuance by id: %w", err)
	}
	return &Issuance, nil
}

func (r *SQLiteRepository) FindAllIssuances() ([]*entity.Issuance, error) {
	var Issuances []*entity.Issuance
	if err := r.Db.
		Preload("Orders").
		Find(&Issuances).Error; err != nil {
		return nil, fmt.Errorf("failed to find all issuances: %w", err)
	}
	return Issuances, nil
}

func (r *SQLiteRepository) FindIssuancesByInvestorAddress(investor types.Address) ([]*entity.Issuance, error) {
	var Issuances []*entity.Issuance
	if err := r.Db.
		Joins("JOIN orders ON orders.issuance_id = issuances.id").
		Where("orders.investor_address = ?", investor).
		Preload("Orders").
		Find(&Issuances).Error; err != nil {
		return nil, fmt.Errorf("failed to find Issuances by investor: %w", err)
	}
	return Issuances, nil
}

func (r *SQLiteRepository) FindIssuancesByCreatorAddress(creator types.Address) ([]*entity.Issuance, error) {
	var Issuances []*entity.Issuance
	if err := r.Db.
		Where("creator_address = ?", creator).
		Preload("Orders").
		Find(&Issuances).Error; err != nil {
		return nil, fmt.Errorf("failed to find issuances by creator: %w", err)
	}
	return Issuances, nil
}

func (r *SQLiteRepository) UpdateIssuance(input *entity.Issuance) (*entity.Issuance, error) {
	if err := r.Db.Save(input).Error; err != nil {
		return nil, fmt.Errorf("failed to update issuance: %w", err)
	}
	issuance, err := r.FindIssuanceById(input.Id)
	if err != nil {
		return nil, err
	}
	return issuance, nil
}
