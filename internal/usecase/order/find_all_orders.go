package order

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllOrdersOutputDTO []*OrderOutputDTO

type FindAllOrdersUseCase struct {
	userRepository  repository.UserRepository
	orderRepository repository.OrderRepository
}

func NewFindAllOrdersUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		userRepository:  userRepo,
		orderRepository: orderRepo,
	}
}

func (f *FindAllOrdersUseCase) Execute() (*FindAllOrdersOutputDTO, error) {
	res, err := f.orderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}
	output := make(FindAllOrdersOutputDTO, len(res))
	for i, order := range res {
		investor, err := f.userRepository.FindUserByAddress(order.Investor)
		if err != nil {
			return nil, err
		}
		output[i] = &OrderOutputDTO{
			Id:           order.Id,
			CampaignId:   order.CampaignId,
			Investor:     investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        string(order.State),
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return &output, nil
}
