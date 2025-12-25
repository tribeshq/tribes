package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type FindOrdersByCampaignIdInputDTO struct {
	CampaignId uint `json:"campaign_id" validate:"required"`
}

type FindOrdersByCampaignIdOutputDTO []*OrderOutputDTO

type FindOrdersByCampaignIdUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByCampaignIdUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrdersByCampaignIdUseCase {
	return &FindOrdersByCampaignIdUseCase{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (c *FindOrdersByCampaignIdUseCase) Execute(input *FindOrdersByCampaignIdInputDTO) (*FindOrdersByCampaignIdOutputDTO, error) {
	res, err := c.OrderRepository.FindOrdersByCampaignId(input.CampaignId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByCampaignIdOutputDTO, len(res))
	for i, order := range res {
		investor, err := c.UserRepository.FindUserByAddress(order.Investor)
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
	return &output, nil
}
