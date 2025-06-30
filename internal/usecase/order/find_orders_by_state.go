package order

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindOrdersByStateInputDTO struct {
	CampaignId uint   `json:"campaign_id" validate:"required"`
	State      string `json:"state" validate:"required"`
}

type FindOrdersByStateOutputDTO []*FindOrderOutputDTO

type FindOrdersByStateUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByStateUseCase(orderRepository repository.OrderRepository) *FindOrdersByStateUseCase {
	return &FindOrdersByStateUseCase{
		OrderRepository: orderRepository,
	}
}

func (f *FindOrdersByStateUseCase) Execute(ctx context.Context, input *FindOrdersByStateInputDTO) (FindOrdersByStateOutputDTO, error) {
	res, err := f.OrderRepository.FindOrdersByState(ctx, input.CampaignId, input.State)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByStateOutputDTO, len(res))
	for i, order := range res {
		output[i] = &FindOrderOutputDTO{
			Id:           order.Id,
			CampaignId:   order.CampaignId,
			BadgeChainId: order.BadgeChainId,
			Investor:     order.Investor,
			Amount:       order.Amount,
			InterestRate: order.InterestRate,
			State:        string(order.State),
			CreatedAt:    order.CreatedAt,
			UpdatedAt:    order.UpdatedAt,
		}
	}
	return output, nil
}
