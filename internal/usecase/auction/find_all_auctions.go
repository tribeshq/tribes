package auction

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllAuctionsOutputDTO []*FindAuctionOutputDTO

type FindAllAuctionsUseCase struct {
	AuctionRepository repository.AuctionRepository
}

func NewFindAllAuctionsUseCase(AuctionRepository repository.AuctionRepository) *FindAllAuctionsUseCase {
	return &FindAllAuctionsUseCase{AuctionRepository: AuctionRepository}
}

func (f *FindAllAuctionsUseCase) Execute(ctx context.Context) (*FindAllAuctionsOutputDTO, error) {
	res, err := f.AuctionRepository.FindAllAuctions(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllAuctionsOutputDTO, len(res))
	for i, Auction := range res {
		orders := make([]*entity.Order, len(Auction.Orders))
		for j, order := range Auction.Orders {
			orders[j] = &entity.Order{
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
		output[i] = &FindAuctionOutputDTO{
			Id:                Auction.Id,
			Token:             Auction.Token,
			Creator:           Auction.Creator,
			CollateralAddress: Auction.CollateralAddress,
			CollateralAmount:  Auction.CollateralAmount,
			DebtIssued:        Auction.DebtIssued,
			MaxInterestRate:   Auction.MaxInterestRate,
			TotalObligation:   Auction.TotalObligation,
			TotalRaised:       Auction.TotalRaised,
			State:             string(Auction.State),
			Orders:            orders,
			CreatedAt:         Auction.CreatedAt,
			ClosesAt:          Auction.ClosesAt,
			MaturityAt:        Auction.MaturityAt,
			UpdatedAt:         Auction.UpdatedAt,
		}
	}
	return &output, nil
}
