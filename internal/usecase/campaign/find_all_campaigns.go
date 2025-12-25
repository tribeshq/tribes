package campaign

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
)

type FindAllCampaignsOutputDTO []*CampaignOutputDTO

type FindAllCampaignsUseCase struct {
	UserRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository
}

func NewFindAllCampaignsUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository) *FindAllCampaignsUseCase {
	return &FindAllCampaignsUseCase{
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
	}
}

func (f *FindAllCampaignsUseCase) Execute() (*FindAllCampaignsOutputDTO, error) {
	res, err := f.campaignRepository.FindAllCampaigns()
	if err != nil {
		return nil, err
	}
	output := make(FindAllCampaignsOutputDTO, len(res))
	for i, campaign := range res {
		orders := make([]*entity.Order, len(campaign.Orders))
		for j, order := range campaign.Orders {
			orders[j] = &entity.Order{
				Id:           order.Id,
				CampaignId:   order.CampaignId,
				Investor:     order.Investor,
				Amount:       order.Amount,
				InterestRate: order.InterestRate,
				State:        order.State,
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    order.UpdatedAt,
			}
		}
		creator, err := f.UserRepository.FindUserByAddress(campaign.Creator)
		if err != nil {
			return nil, fmt.Errorf("error finding creator: %w", err)
		}
		output[i] = &CampaignOutputDTO{
			Id:          campaign.Id,
			Title:       campaign.Title,
			Description: campaign.Description,
			Promotion:   campaign.Promotion,
			Token:       campaign.Token,
			Creator: &user.UserOutputDTO{
				Id:             creator.Id,
				Role:           string(creator.Role),
				Address:        creator.Address,
				SocialAccounts: creator.SocialAccounts,
				CreatedAt:      creator.CreatedAt,
				UpdatedAt:      creator.UpdatedAt,
			},
			CollateralAddress: campaign.CollateralAddress,
			CollateralAmount:  campaign.CollateralAmount,
			BadgeAddress:      campaign.BadgeAddress,
			DebtIssued:        campaign.DebtIssued,
			MaxInterestRate:   campaign.MaxInterestRate,
			TotalObligation:   campaign.TotalObligation,
			TotalRaised:       campaign.TotalRaised,
			State:             string(campaign.State),
			Orders:            orders,
			CreatedAt:         campaign.CreatedAt,
			ClosesAt:          campaign.ClosesAt,
			MaturityAt:        campaign.MaturityAt,
			UpdatedAt:         campaign.UpdatedAt,
		}
	}
	return &output, nil
}
