package crowdfunding_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllCrowdfundingsOutputDTO []*FindCrowdfundingOutputDTO

type FindAllCrowdfundingsUseCase struct {
	CrowdfundingRepository repository.CrowdfundingRepository
}

func NewFindAllCrowdfundingsUseCase(crowdfundingRepository repository.CrowdfundingRepository) *FindAllCrowdfundingsUseCase {
	return &FindAllCrowdfundingsUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindAllCrowdfundingsUseCase) Execute(ctx context.Context) (*FindAllCrowdfundingsOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindAllCrowdfundings(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllCrowdfundingsOutputDTO, len(res))
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
