package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	. "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type FindOrdersByInvestorAddressInputDTO struct {
	InvestorAddress Address `json:"investor_address" validate:"required"`
}

type FindOrdersByInvestorAddressOutputDTO []*OrderOutputDTO

type FindOrdersByInvestorAddressUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByInvestorAddressUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrdersByInvestorAddressUseCase {
	return &FindOrdersByInvestorAddressUseCase{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (o *FindOrdersByInvestorAddressUseCase) Execute(input *FindOrdersByInvestorAddressInputDTO) (FindOrdersByInvestorAddressOutputDTO, error) {
	res, err := o.OrderRepository.FindOrdersByInvestorAddress(input.InvestorAddress)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByInvestorAddressOutputDTO, len(res))
	for i, order := range res {
		investor, err := o.UserRepository.FindUserByAddress(order.InvestorAddress)
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
	return output, nil
}
