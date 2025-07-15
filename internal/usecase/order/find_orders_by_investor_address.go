package order

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindOrdersByInvestorAddressInputDTO struct {
	InvestorAddress custom_type.Address `json:"investor_address" validate:"required"`
}

type FindOrdersByInvestorAddressOutputDTO []*OrderOutputDTO

type FindOrdersByInvestorAddressUseCase struct {
	userRepository  repository.UserRepository
	orderRepository repository.OrderRepository
}

func NewFindOrdersByInvestorAddressUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrdersByInvestorAddressUseCase {
	return &FindOrdersByInvestorAddressUseCase{
		userRepository:  userRepo,
		orderRepository: orderRepo,
	}
}

func (o *FindOrdersByInvestorAddressUseCase) Execute(input *FindOrdersByInvestorAddressInputDTO) (FindOrdersByInvestorAddressOutputDTO, error) {
	res, err := o.orderRepository.FindOrdersByInvestorAddress(input.InvestorAddress)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByInvestorAddressOutputDTO, len(res))
	for i, order := range res {
		investor, err := o.userRepository.FindUserByAddress(order.Investor)
		if err != nil {
			return nil, err
		}
		output[i] = &OrderOutputDTO{
			Id:         order.Id,
			CampaignId: order.CampaignId,

			Investor:     investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        string(order.State),
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return output, nil
}
