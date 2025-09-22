package sqlite

import (
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
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

func (r *SQLiteRepository) FindCampaignsByInvestorAddress(investor custom_type.Address) ([]*entity.Campaign, error) {
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

func (r *SQLiteRepository) FindCampaignsByCreatorAddress(creator custom_type.Address) ([]*entity.Campaign, error) {
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
	if err := r.Db.Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}
	campaign, err := r.FindCampaignById(input.Id)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}

func (r *SQLiteRepository) FindOngoingCampaignByCreatorAddress(creator custom_type.Address) (*entity.Campaign, error) {
	var campaign entity.Campaign
	if err := r.Db.
		Where("creator = ? AND state = ?", creator, entity.CampaignStateOngoing).
		Preload("Orders").
		First(&campaign).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find ongoing campaign by creator: %w", err)
	}
	return &campaign, nil
}
