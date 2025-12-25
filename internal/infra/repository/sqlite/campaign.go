package sqlite

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateCampaign(input *entity.Campaign) (*entity.Campaign, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindCampaignById(id uint) (*entity.Campaign, error) {
	var Campaign entity.Campaign
	if err := r.Db.
		Preload("Orders").
		First(&Campaign, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("failed to find campaign by id: %w", err)
	}
	return &Campaign, nil
}

func (r *SQLiteRepository) FindAllCampaigns() ([]*entity.Campaign, error) {
	var Campaigns []*entity.Campaign
	if err := r.Db.
		Preload("Orders").
		Find(&Campaigns).Error; err != nil {
		return nil, fmt.Errorf("failed to find all campaigns: %w", err)
	}
	return Campaigns, nil
}

func (r *SQLiteRepository) FindCampaignsByInvestorAddress(investor types.Address) ([]*entity.Campaign, error) {
	var Campaigns []*entity.Campaign
	if err := r.Db.
		Joins("JOIN orders ON orders.campaign_id = campaigns.id").
		Where("orders.investor = ?", investor).
		Preload("Orders").
		Find(&Campaigns).Error; err != nil {
		return nil, fmt.Errorf("failed to find Campaigns by investor: %w", err)
	}
	return Campaigns, nil
}

func (r *SQLiteRepository) FindCampaignsByCreatorAddress(creator types.Address) ([]*entity.Campaign, error) {
	var Campaigns []*entity.Campaign
	if err := r.Db.
		Where("creator = ?", creator).
		Preload("Orders").
		Find(&Campaigns).Error; err != nil {
		return nil, fmt.Errorf("failed to find campaigns by creator: %w", err)
	}
	return Campaigns, nil
}

func (r *SQLiteRepository) UpdateCampaign(input *entity.Campaign) (*entity.Campaign, error) {
	if err := r.Db.Save(input).Error; err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}
	campaign, err := r.FindCampaignById(input.Id)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}
