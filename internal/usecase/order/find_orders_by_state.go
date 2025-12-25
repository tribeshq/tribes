package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type FindOrdersByStateInputDTO struct {
	CampaignId uint   `json:"campaign_id" validate:"required"`
	State      string `json:"state" validate:"required"`
}

type FindOrdersByStateOutputDTO []*OrderOutputDTO

type FindOrdersByStateUseCase struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewFindOrdersByStateUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *FindOrdersByStateUseCase {
	return &FindOrdersByStateUseCase{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (f *FindOrdersByStateUseCase) Execute(input *FindOrdersByStateInputDTO) (FindOrdersByStateOutputDTO, error) {
	res, err := f.OrderRepository.FindOrdersByState(input.CampaignId, input.State)
	if err != nil {
		return nil, err
	}
	output := make(FindOrdersByStateOutputDTO, len(res))
	for i, order := range res {
		investor, err := f.UserRepository.FindUserByAddress(order.Investor)
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
	return output, nil
}
