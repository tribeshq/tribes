package campaign

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
)

type FindAllCampaignsOutputDTO []*CampaignOutputDTO

type FindAllCampaignsUseCase struct {
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewFindAllCampaignsUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository) *FindAllCampaignsUseCase {
	return &FindAllCampaignsUseCase{
		UserRepository:     userRepo,
		CampaignRepository: campaignRepo,
	}
}

func (f *FindAllCampaignsUseCase) Execute() (*FindAllCampaignsOutputDTO, error) {
	res, err := f.CampaignRepository.FindAllCampaigns()
	if err != nil {
		return nil, err
	}
	output := make(FindAllCampaignsOutputDTO, len(res))
	for i, campaign := range res {
		orders := make([]*order.OrderOutputDTO, len(campaign.Orders))
		for j, o := range campaign.Orders {
			investor, err := f.UserRepository.FindUserByAddress(o.InvestorAddress)
			if err != nil {
				return nil, fmt.Errorf("error finding investor: %w", err)
			}
			orders[j] = &order.OrderOutputDTO{
				Id:         o.Id,
				CampaignId: o.CampaignId,
				Investor: &user.UserOutputDTO{
					Id:             investor.Id,
					Role:           string(investor.Role),
					Address:        investor.Address,
					SocialAccounts: investor.SocialAccounts,
					CreatedAt:      investor.CreatedAt,
					UpdatedAt:      investor.UpdatedAt,
				},
				Amount:       o.Amount,
				InterestRate: o.InterestRate,
				State:        string(o.State),
				CreatedAt:    o.CreatedAt,
				UpdatedAt:    o.UpdatedAt,
			}
		}
		creator, err := f.UserRepository.FindUserByAddress(campaign.CreatorAddress)
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
