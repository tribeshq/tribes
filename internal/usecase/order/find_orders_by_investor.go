package order

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type FindOrdersByInvestorInputDTO struct {
	Investor Address `json:"investor" validate:"required"`
}

type FindOrdersByInvestorOutputDTO []*FindOrderOutputDTO

type FindOrdersByInvestorUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByInvestorUseCase(orderRepository repository.OrderRepository) *FindOrdersByInvestorUseCase {
	return &FindOrdersByInvestorUseCase{
		OrderRepository: orderRepository,
	}
}

func (o *FindOrdersByInvestorUseCase) Execute(ctx context.Context, input *FindOrdersByInvestorInputDTO) (FindOrdersByInvestorOutputDTO, error) {
	res, err := o.OrderRepository.FindOrdersByInvestor(ctx, input.Investor)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByInvestorOutputDTO, len(res))
	for i, order := range res {
		output[i] = &FindOrderOutputDTO{
			Id:           order.Id,
			CampaignId:   order.CampaignId,
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
