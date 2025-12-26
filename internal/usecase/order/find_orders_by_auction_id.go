package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
)

type FindOrdersByIssuanceIdInputDTO struct {
	IssuanceId uint `json:"issuance_id" validate:"required"`
}

type FindOrdersByIssuanceIdOutputDTO []*OrderOutputDTO

type FindOrdersByIssuanceIdUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByIssuanceIdUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrdersByIssuanceIdUseCase {
	return &FindOrdersByIssuanceIdUseCase{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (c *FindOrdersByIssuanceIdUseCase) Execute(input *FindOrdersByIssuanceIdInputDTO) (*FindOrdersByIssuanceIdOutputDTO, error) {
	res, err := c.OrderRepository.FindOrdersByIssuanceId(input.IssuanceId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByIssuanceIdOutputDTO, len(res))
	for i, order := range res {
		investor, err := c.UserRepository.FindUserByAddress(order.InvestorAddress)
		if err != nil {
			return nil, err
		}
		output[i] = &OrderOutputDTO{
			Id:         order.Id,
			IssuanceId: order.IssuanceId,
			Investor:     &user.UserOutputDTO{
				Id: investor.Id,
				Role: string(investor.Role),
				Address: investor.Address,
				SocialAccounts: investor.SocialAccounts,
				CreatedAt: investor.CreatedAt,
				UpdatedAt: investor.UpdatedAt,
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
