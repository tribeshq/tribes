package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
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
		investor, err := f.UserRepository.FindUserByAddress(order.InvestorAddress)
		if err != nil {
			return nil, err
		}
		output[i] = &OrderOutputDTO{
			Id:         order.Id,
			IssuanceId: order.IssuanceId,
			Investor: &user.UserOutputDTO{
				Id:             investor.Id,
				Role:           string(investor.Role),
				Address:        investor.Address,
				SocialAccounts: investor.SocialAccounts,
				CreatedAt:      investor.CreatedAt,
				UpdatedAt:      investor.UpdatedAt,
			},
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        string(order.State),
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return &output, nil
}
