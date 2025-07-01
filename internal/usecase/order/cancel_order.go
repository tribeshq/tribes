package order

import (
	"context"
	"errors"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CancelOrderInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type CancelOrderOutputDTO struct {
	Id           uint
	CampaignId   uint
	BadgeChainId uint64
	Token        custom_type.Address
	Investor     custom_type.Address
	Amount       *uint256.Int
	InterestRate *uint256.Int
	State        string
	CreatedAt    int64
	UpdatedAt    int64
}

type CancelOrderUseCase struct {
	OrderRepository    repository.OrderRepository
	CampaignRepository repository.CampaignRepository
}

func NewCancelOrderUseCase(orderRepository repository.OrderRepository, campaignRepository repository.CampaignRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		OrderRepository:    orderRepository,
		CampaignRepository: campaignRepository,
	}
}

func (c *CancelOrderUseCase) Execute(ctx context.Context, input *CancelOrderInputDTO, metadata rollmelette.Metadata) (*CancelOrderOutputDTO, error) {
	order, err := c.OrderRepository.FindOrderById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	if order.Investor != custom_type.Address(metadata.MsgSender) {
		return nil, errors.New("only the investor can cancel the order")
	}
	campaign, err := c.CampaignRepository.FindCampaignById(ctx, order.CampaignId)
	if err != nil {
		return nil, err
	}
	if campaign.State == entity.CampaignStateClosed {
		return nil, errors.New("cannot cancel order after Campaign closes")
	}
	order.State = entity.OrderStateCancelled
	res, err := c.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &CancelOrderOutputDTO{
		Id:           res.Id,
		CampaignId:   res.CampaignId,
		BadgeChainId: res.BadgeChainId,
		Token:        campaign.Token,
		Investor:     res.Investor,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}
