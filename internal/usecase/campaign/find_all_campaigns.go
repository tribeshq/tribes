package campaign

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllCampaignsOutputDTO []*FindCampaignOutputDTO

type FindAllCampaignsUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindAllCampaignsUseCase(CampaignRepository repository.CampaignRepository) *FindAllCampaignsUseCase {
	return &FindAllCampaignsUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindAllCampaignsUseCase) Execute(ctx context.Context) (*FindAllCampaignsOutputDTO, error) {
	res, err := f.CampaignRepository.FindAllCampaigns(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllCampaignsOutputDTO, len(res))
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
