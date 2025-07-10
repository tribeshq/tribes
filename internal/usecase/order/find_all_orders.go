package order

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllOrdersOutputDTO []*OrderOutputDTO

type FindAllOrdersUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindAllOrdersUseCase(userRepository repository.UserRepository, orderRepository repository.OrderRepository) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		UserRepository:  userRepository,
		OrderRepository: orderRepository,
	}
}

func (f *FindAllOrdersUseCase) Execute(ctx context.Context) (*FindAllOrdersOutputDTO, error) {
	res, err := f.OrderRepository.FindAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllOrdersOutputDTO, len(res))
	for i, order := range res {
		investor, err := f.UserRepository.FindUserByAddress(ctx, order.Investor)
		if err != nil {
			return nil, err
		}
		output[i] = &OrderOutputDTO{
			Id:                 order.Id,
			CampaignId:         order.CampaignId,
			BadgeChainSelector: order.BadgeChainSelector,
			Investor:           investor,
			Amount:             order.Amount,
			InterestRate:       order.InterestRate,
			State:              string(order.State),
			CreatedAt:          order.CreatedAt,
			UpdatedAt:          order.UpdatedAt,
		}
	}
	return &output, nil
}
