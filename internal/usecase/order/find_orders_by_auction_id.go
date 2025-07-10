package order

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindOrdersByCampaignIdInputDTO struct {
	CampaignId uint `json:"campaign_id" validate:"required"`
}

type FindOrdersByCampaignIdOutputDTO []*OrderOutputDTO

type FindOrdersByCampaignIdUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByCampaignIdUseCase(userRepository repository.UserRepository, orderRepository repository.OrderRepository) *FindOrdersByCampaignIdUseCase {
	return &FindOrdersByCampaignIdUseCase{
		UserRepository:  userRepository,
		OrderRepository: orderRepository,
	}
}

func (c *FindOrdersByCampaignIdUseCase) Execute(ctx context.Context, input *FindOrdersByCampaignIdInputDTO) (*FindOrdersByCampaignIdOutputDTO, error) {
	res, err := c.OrderRepository.FindOrdersByCampaignId(ctx, input.CampaignId)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByCampaignIdOutputDTO, len(res))
	for i, order := range res {
		investor, err := c.UserRepository.FindUserByAddress(ctx, order.Investor)
		if err != nil {
			return nil, err
		}
		output[i] = &OrderOutputDTO{
			Id:                 order.Id,
			CampaignId:         order.CampaignId,
			BadgeChainSelector: order.BadgeChainSelector,
			Investor:           investor,
			Amount:             order.Amount,
			InterestRate:       order.InterestRate,
			State:              string(order.State),
			CreatedAt:          order.CreatedAt,
			UpdatedAt:          order.UpdatedAt,
		}
	}
	return &output, nil
}
