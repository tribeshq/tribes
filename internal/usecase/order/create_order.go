package order

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CreateOrderInputDTO struct {
	CampaignId   uint         `json:"campaign_id" validate:"required"`
	InterestRate *uint256.Int `json:"interest_rate" validate:"required"`
}

type CreateOrderOutputDTO struct {
	Id           uint                `json:"id"`
	CampaignId   uint                `json:"campaign_id"`
	Investor     *user.UserOutputDTO `json:"investor"`
	Amount       *uint256.Int        `json:"amount"`
	InterestRate *uint256.Int        `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
}

type CreateOrderUseCase struct {
	UserRepository     repository.UserRepository
	OrderRepository    repository.OrderRepository
	campaignRepository repository.CampaignRepository
}

func NewCreateOrderUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
	campaignRepo repository.CampaignRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		UserRepository:     userRepo,
		OrderRepository:    orderRepo,
		campaignRepository: campaignRepo,
	}
}

func (c *CreateOrderUseCase) Execute(input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type provided for order creation: %T", deposit)
	}

	campaign, err := c.campaignRepository.FindCampaignById(input.CampaignId)
	if err != nil {
		return nil, fmt.Errorf("error finding campaign campaigns: %w", err)
	}

	if campaign.ClosesAt < metadata.BlockTimestamp {
		return nil, fmt.Errorf("campaign campaign closed, order cannot be placed")
	}

	if types.Address(erc20Deposit.Token) != campaign.Token {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(campaign.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active Campaign max interest rate")
	}

	order, err := entity.NewOrder(
		campaign.Id,
		types.Address(erc20Deposit.Sender),
		uint256.MustFromBig(erc20Deposit.Value),
		input.InterestRate,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.OrderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	investor, err := c.UserRepository.FindUserByAddress(res.Investor)
	if err != nil {
		return nil, err
	}

	return &CreateOrderOutputDTO{
		Id:         res.Id,
		CampaignId: res.CampaignId,
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
	}, nil
}
