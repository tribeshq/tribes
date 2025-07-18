package order

import (
	"errors"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CancelOrderInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type CancelOrderOutputDTO struct {
	Id           uint                `json:"id"`
	CampaignId   uint                `json:"campaign_id"`
	Token        custom_type.Address `json:"token"`
	Investor     *user.UserOutputDTO `json:"investor"`
	Amount       *uint256.Int        `json:"amount"`
	InterestRate *uint256.Int        `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
	UpdatedAt    int64               `json:"updated_at"`
}

type CancelOrderUseCase struct {
	userRepository     repository.UserRepository
	orderRepository    repository.OrderRepository
	campaignRepository repository.CampaignRepository
}

func NewCancelOrderUseCase(userRepo repository.UserRepository, orderRepo repository.OrderRepository, campaignRepo repository.CampaignRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		userRepository:     userRepo,
		orderRepository:    orderRepo,
		campaignRepository: campaignRepo,
	}
}

func (c *CancelOrderUseCase) Execute(input *CancelOrderInputDTO, metadata rollmelette.Metadata) (*CancelOrderOutputDTO, error) {
	order, err := c.orderRepository.FindOrderById(input.Id)
	if err != nil {
		return nil, err
	}
	if order.Investor != custom_type.Address(metadata.MsgSender) {
		return nil, errors.New("only the investor can cancel the order")
	}
	campaign, err := c.campaignRepository.FindCampaignById(order.CampaignId)
	if err != nil {
		return nil, err
	}
	if campaign.State == entity.CampaignStateClosed {
		return nil, errors.New("cannot cancel order after Campaign closes")
	}
	order.State = entity.OrderStateCancelled
	res, err := c.orderRepository.UpdateOrder(order)
	if err != nil {
		return nil, err
	}
	investor, err := c.userRepository.FindUserByAddress(res.Investor)
	if err != nil {
		return nil, err
	}
	return &CancelOrderOutputDTO{
		Id:         res.Id,
		CampaignId: res.CampaignId,
		Token:      campaign.Token,
		Investor: &user.UserOutputDTO{
			Id:             investor.Id,
			Role:           string(investor.Role),
			Address:        investor.Address,
			SocialAccounts: investor.SocialAccounts,
			CreatedAt:      investor.CreatedAt,
			UpdatedAt:      investor.UpdatedAt,
		},
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}
