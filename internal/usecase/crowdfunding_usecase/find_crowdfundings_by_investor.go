package crowdfunding_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCrowdfundingsByInvestorInputDTO struct {
	Investor Address `json:"investor"`
}

type FindCrowdfundingsByInvestorOutputDTO []*FindCrowdfundingOutputDTO

type FindCrowdfundingsByInvestorUseCase struct {
	CrowdfundingRepository repository.CrowdfundingRepository
}

func NewFindCrowdfundingsByInvestorUseCase(crowdfundingRepository repository.CrowdfundingRepository) *FindCrowdfundingsByInvestorUseCase {
	return &FindCrowdfundingsByInvestorUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindCrowdfundingsByInvestorUseCase) Execute(ctx context.Context, input *FindCrowdfundingsByInvestorInputDTO) (*FindCrowdfundingsByInvestorOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindCrowdfundingsByInvestor(ctx, input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindCrowdfundingsByInvestorOutputDTO, len(res))
	for i, crowdfunding := range res {
		orders := make([]*entity.Order, len(crowdfunding.Orders))
		for j, order := range crowdfunding.Orders {
			orders[j] = &entity.Order{
				Id:             order.Id,
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         order.Amount,
				InterestRate:   order.InterestRate,
				State:          order.State,
				CreatedAt:      order.CreatedAt,
				UpdatedAt:      order.UpdatedAt,
			}
		}
		output[i] = &FindCrowdfundingOutputDTO{
			Id:                  crowdfunding.Id,
			Token:               crowdfunding.Token,
			Collateral:          crowdfunding.Collateral,
			Creator:             crowdfunding.Creator,
			DebtIssued:          crowdfunding.DebtIssued,
			MaxInterestRate:     crowdfunding.MaxInterestRate,
			TotalObligation:     crowdfunding.TotalObligation,
			Orders:              orders,
			State:               string(crowdfunding.State),
			FundraisingDuration: crowdfunding.FundraisingDuration,
			ClosesAt:            crowdfunding.ClosesAt,
			MaturityAt:          crowdfunding.MaturityAt,
			CreatedAt:           crowdfunding.CreatedAt,
			UpdatedAt:           crowdfunding.UpdatedAt,
		}
	}
	return &output, nil
}
