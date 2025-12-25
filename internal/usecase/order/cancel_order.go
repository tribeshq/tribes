package order

import (
	"errors"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CancelOrderInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type CancelOrderOutputDTO struct {
	Id           uint                `json:"id"`
	CampaignId   uint                `json:"campaign_id"`
	Token        types.Address       `json:"token"`
	Investor     *user.UserOutputDTO `json:"investor"`
	Amount       *uint256.Int        `json:"amount"`
	InterestRate *uint256.Int        `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
	UpdatedAt    int64               `json:"updated_at"`
}

type CancelOrderUseCase struct {
	UserRepository     repository.UserRepository
	OrderRepository    repository.OrderRepository
	CampaignRepository repository.CampaignRepository
}

func NewCancelOrderUseCase(userRepo repository.UserRepository, orderRepo repository.OrderRepository, campaignRepo repository.CampaignRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		UserRepository:     userRepo,
		OrderRepository:    orderRepo,
		CampaignRepository: campaignRepo,
	}
}

func (c *CancelOrderUseCase) Execute(input *CancelOrderInputDTO, metadata rollmelette.Metadata) (*CancelOrderOutputDTO, error) {
	order, err := c.OrderRepository.FindOrderById(input.Id)
	if err != nil {
		return nil, err
	}
	if order.InvestorAddress != types.Address(metadata.MsgSender) {
		return nil, errors.New("only the investor can cancel the order")
	}
	campaign, err := c.CampaignRepository.FindCampaignById(order.CampaignId)
	if err != nil {
		return nil, err
	}
	if campaign.State == entity.CampaignStateClosed {
		return nil, errors.New("cannot cancel order after Campaign closes")
	}
	order.State = entity.OrderStateCancelled
	res, err := c.OrderRepository.UpdateOrder(order)
	if err != nil {
		return nil, err
	}
	investor, err := c.UserRepository.FindUserByAddress(res.InvestorAddress)
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
