package campaign

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user"
)

type FindAllCampaignsOutputDTO []*CampaignOutputDTO

type FindAllCampaignsUseCase struct {
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewFindAllCampaignsUseCase(userRepository repository.UserRepository, campaignRepository repository.CampaignRepository) *FindAllCampaignsUseCase {
	return &FindAllCampaignsUseCase{UserRepository: userRepository, CampaignRepository: campaignRepository}
}

func (f *FindAllCampaignsUseCase) Execute(ctx context.Context) (*FindAllCampaignsOutputDTO, error) {
	res, err := f.CampaignRepository.FindAllCampaigns(ctx)
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
		creator, err := f.UserRepository.FindUserByAddress(ctx, campaign.Creator)
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
