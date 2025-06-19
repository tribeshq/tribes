package auction

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAuctionByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindAuctionByIdUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewFindAuctionByIdUseCase(AuctionRepository repository.AuctionRepository) *FindAuctionByIdUseCase {
	return &FindAuctionByIdUseCase{AuctionRepository: AuctionRepository}
}

func (f *FindAuctionByIdUseCase) Execute(ctx context.Context, input *FindAuctionByIdInputDTO) (*FindAuctionOutputDTO, error) {
	res, err := f.AuctionRepository.FindAuctionById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	orders := make([]*entity.Order, len(res.Orders))
	for i, order := range res.Orders {
		orders[i] = &entity.Order{
			Id:           order.Id,
			AuctionId:    order.AuctionId,
			Investor:     order.Investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        order.State,
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return &FindAuctionOutputDTO{
		Id:                res.Id,
		Token:             res.Token,
		Creator:           res.Creator,
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		State:             string(res.State),
		Orders:            orders,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
