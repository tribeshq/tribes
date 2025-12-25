package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type FindOrderByIdInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type FindOrderByIdUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrderByIdUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrderByIdUseCase {
	return &FindOrderByIdUseCase{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (c *FindOrderByIdUseCase) Execute(input *FindOrderByIdInputDTO) (*OrderOutputDTO, error) {
	res, err := c.OrderRepository.FindOrderById(input.Id)
	if err != nil {
		return nil, err
	}
	investor, err := c.UserRepository.FindUserByAddress(res.Investor)
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
