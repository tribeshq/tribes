package crowdfunding_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindCrowdfundingByIdInputDTO struct {
	Id uint `json:"id"`
}

type FindCrowdfundingByIdUseCase struct {
	CrowdfundingRepository repository.CrowdfundingRepository
}

func NewFindCrowdfundingByIdUseCase(crowdfundingRepository repository.CrowdfundingRepository) *FindCrowdfundingByIdUseCase {
	return &FindCrowdfundingByIdUseCase{CrowdfundingRepository: crowdfundingRepository}
}

func (f *FindCrowdfundingByIdUseCase) Execute(ctx context.Context, input *FindCrowdfundingByIdInputDTO) (*FindCrowdfundingOutputDTO, error) {
	res, err := f.CrowdfundingRepository.FindCrowdfundingById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	var orders []*entity.Order
	for _, order := range res.Orders {
		orders = append(orders, &entity.Order{
			Id:             order.Id,
			CrowdfundingId: order.CrowdfundingId,
			Investor:       order.Investor,
			Amount:         order.Amount,
			InterestRate:   order.InterestRate,
			State:          order.State,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		})
	}
	return &FindCrowdfundingOutputDTO{
		Id:                  res.Id,
		Token:               res.Token,
		Collateral:          res.Collateral,
		Creator:             res.Creator,
		DebtIssued:          res.DebtIssued,
		MaxInterestRate:     res.MaxInterestRate,
		TotalObligation:     res.TotalObligation,
		Orders:              orders,
		State:               string(res.State),
		FundraisingDuration: res.FundraisingDuration,
		ClosesAt:            res.ClosesAt,
		MaturityAt:          res.MaturityAt,
		CreatedAt:           res.CreatedAt,
		UpdatedAt:           res.UpdatedAt,
	}, nil
}
