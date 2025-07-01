package campaign

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCampaignsByInvestorInputDTO struct {
	Investor custom_type.Address `json:"investor" validate:"required"`
}

type FindCampaignsByInvestorOutputDTO []*FindCampaignOutputDTO

type FindCampaignsByInvestorUseCase struct {
	CampaignRepository repository.CampaignRepository
}

func NewFindCampaignsByInvestorUseCase(CampaignRepository repository.CampaignRepository) *FindCampaignsByInvestorUseCase {
	return &FindCampaignsByInvestorUseCase{CampaignRepository: CampaignRepository}
}

func (f *FindCampaignsByInvestorUseCase) Execute(ctx context.Context, input *FindCampaignsByInvestorInputDTO) (*FindCampaignsByInvestorOutputDTO, error) {
	res, err := f.CampaignRepository.FindCampaignsByInvestor(ctx, input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindCampaignsByInvestorOutputDTO, len(res))
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
