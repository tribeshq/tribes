package order

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindOrderByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindOrderByIdUseCase struct {
	userRepository  repository.UserRepository
	orderRepository repository.OrderRepository
}

func NewFindOrderByIdUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrderByIdUseCase {
	return &FindOrderByIdUseCase{
		userRepository:  userRepo,
		orderRepository: orderRepo,
	}
}

func (c *FindOrderByIdUseCase) Execute(input *FindOrderByIdInputDTO) (*OrderOutputDTO, error) {
	res, err := c.orderRepository.FindOrderById(input.Id)
	if err != nil {
		return nil, err
	}
	investor, err := c.userRepository.FindUserByAddress(res.Investor)
	if err != nil {
		return nil, err
	}
	return &OrderOutputDTO{
		Id:           res.Id,
		CampaignId:   res.CampaignId,
		Investor:     investor,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}
