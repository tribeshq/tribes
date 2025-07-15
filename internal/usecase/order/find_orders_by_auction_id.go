package order

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindOrdersByCampaignIdInputDTO struct {
	CampaignId uint `json:"campaign_id" validate:"required"`
}

type FindOrdersByCampaignIdOutputDTO []*OrderOutputDTO

type FindOrdersByCampaignIdUseCase struct {
	userRepository  repository.UserRepository
	orderRepository repository.OrderRepository
}

func NewFindOrdersByCampaignIdUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrdersByCampaignIdUseCase {
	return &FindOrdersByCampaignIdUseCase{
		userRepository:  userRepo,
		orderRepository: orderRepo,
	}
}

func (c *FindOrdersByCampaignIdUseCase) Execute(input *FindOrdersByCampaignIdInputDTO) (*FindOrdersByCampaignIdOutputDTO, error) {
	res, err := c.orderRepository.FindOrdersByCampaignId(input.CampaignId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByCampaignIdOutputDTO, len(res))
	for i, order := range res {
		investor, err := c.userRepository.FindUserByAddress(order.Investor)
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
