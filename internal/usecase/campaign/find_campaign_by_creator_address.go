package campaign

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type FindCampaignsByCreatorAddressInputDTO struct {
	CreatorAddress types.Address `json:"creator" validate:"required"`
}

type FindCampaignsByCreatorAddressOutputDTO []*CampaignOutputDTO

type FindCampaignsByCreatorAddressUseCase struct {
	UserRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository
}

func NewFindCampaignsByCreatorAddressUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository) *FindCampaignsByCreatorAddressUseCase {
	return &FindCampaignsByCreatorAddressUseCase{
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
	}
}

func (f *FindCampaignsByCreatorAddressUseCase) Execute(input *FindCampaignsByCreatorAddressInputDTO) (*FindCampaignsByCreatorAddressOutputDTO, error) {
	res, err := f.campaignRepository.FindCampaignsByCreatorAddress(input.CreatorAddress)
	if err != nil {
		return nil, err
	}
	output := make(FindCampaignsByCreatorAddressOutputDTO, len(res))
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
