package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type FindAllOrdersOutputDTO []*OrderOutputDTO

type FindAllOrdersUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindAllOrdersUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindAllOrdersUseCase {
	return &FindAllOrdersUseCase{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (f *FindAllOrdersUseCase) Execute() (*FindAllOrdersOutputDTO, error) {
	res, err := f.OrderRepository.FindAllOrders()
	if err != nil {
		return nil, err
	}
	output := make(FindAllOrdersOutputDTO, len(res))
	for i, order := range res {
		investor, err := f.UserRepository.FindUserByAddress(order.Investor)
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
