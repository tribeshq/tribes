package order

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindOrderByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindOrderByIdUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrderByIdUseCase(userRepository repository.UserRepository, orderRepository repository.OrderRepository) *FindOrderByIdUseCase {
	return &FindOrderByIdUseCase{
		UserRepository:  userRepository,
		OrderRepository: orderRepository,
	}
}

func (c *FindOrderByIdUseCase) Execute(ctx context.Context, input *FindOrderByIdInputDTO) (*OrderOutputDTO, error) {
	res, err := c.OrderRepository.FindOrderById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	investor, err := c.UserRepository.FindUserByAddress(ctx, res.Investor)
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
