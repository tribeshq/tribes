package campaign

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCampaignsByCreatorInputDTO struct {
	Creator custom_type.Address `json:"creator" validate:"required"`
}

type FindCampaignsByCreatorOutputDTO []*FindCampaignOutputDTO

type FindCampaignsByCreatorUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindCampaignsByCreatorUseCase(CampaignRepository repository.CampaignRepository) *FindCampaignsByCreatorUseCase {
	return &FindCampaignsByCreatorUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindCampaignsByCreatorUseCase) Execute(ctx context.Context, input *FindCampaignsByCreatorInputDTO) (*FindCampaignsByCreatorOutputDTO, error) {
	res, err := f.CampaignRepository.FindCampaignsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}
	output := make(FindCampaignsByCreatorOutputDTO, len(res))
	for i, Campaign := range res {
		orders := make([]*entity.Order, len(Campaign.Orders))
		for j, order := range Campaign.Orders {
			orders[j] = &entity.Order{
				Id:           order.Id,
				CampaignId:   order.CampaignId,
				BadgeChainId: order.BadgeChainId,
				Investor:     order.Investor,
				Amount:       order.Amount,
				InterestRate: order.InterestRate,
				State:        order.State,
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    order.UpdatedAt,
			}
		}
		output[i] = &FindCampaignOutputDTO{
			Id:                Campaign.Id,
			Token:             Campaign.Token,
			Creator:           Campaign.Creator,
			CollateralAddress: Campaign.CollateralAddress,
			CollateralAmount:  Campaign.CollateralAmount,
			BadgeRouter:       Campaign.BadgeRouter,
			BadgeMinter:       Campaign.BadgeMinter,
			DebtIssued:        Campaign.DebtIssued,
			MaxInterestRate:   Campaign.MaxInterestRate,
			TotalObligation:   Campaign.TotalObligation,
			TotalRaised:       Campaign.TotalRaised,
			State:             string(Campaign.State),
			Orders:            orders,
			CreatedAt:         Campaign.CreatedAt,
			ClosesAt:          Campaign.ClosesAt,
			MaturityAt:        Campaign.MaturityAt,
			UpdatedAt:         Campaign.UpdatedAt,
		}
	}
	return &output, nil
}
