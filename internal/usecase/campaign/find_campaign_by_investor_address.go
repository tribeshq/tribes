package campaign

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCampaignsByInvestorAddressInputDTO struct {
	InvestorAddress custom_type.Address `json:"investor_address" validate:"required"`
}

type FindCampaignsByInvestorAddressOutputDTO []*CampaignOutputDTO

type FindCampaignsByInvestorAddressUseCase struct {
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewFindCampaignsByInvestorAddressUseCase(userRepository repository.UserRepository, campaignRepository repository.CampaignRepository) *FindCampaignsByInvestorAddressUseCase {
	return &FindCampaignsByInvestorAddressUseCase{UserRepository: userRepository, CampaignRepository: campaignRepository}
}

func (f *FindCampaignsByInvestorAddressUseCase) Execute(ctx context.Context, input *FindCampaignsByInvestorAddressInputDTO) (*FindCampaignsByInvestorAddressOutputDTO, error) {
	res, err := f.CampaignRepository.FindCampaignsByInvestorAddress(ctx, input.InvestorAddress)
	if err != nil {
		return nil, err
	}
	output := make(FindCampaignsByInvestorAddressOutputDTO, len(res))
	for i, campaign := range res {
		orders := make([]*entity.Order, len(campaign.Orders))
		for j, order := range campaign.Orders {
			orders[j] = &entity.Order{
				Id:                 order.Id,
				CampaignId:         order.CampaignId,
				BadgeChainSelector: order.BadgeChainSelector,
				Investor:           order.Investor,
				Amount:             order.Amount,
				InterestRate:       order.InterestRate,
				State:              order.State,
				CreatedAt:          order.CreatedAt,
				UpdatedAt:          order.UpdatedAt,
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
			BadgeRouter:       campaign.BadgeRouter,
			BadgeMinter:       campaign.BadgeMinter,
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
