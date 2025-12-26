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
	IssuanceId   uint         `json:"issuance_id" validate:"required"`
	InterestRate *uint256.Int `json:"interest_rate" validate:"required"`
}

type CreateOrderOutputDTO struct {
	Id           uint                `json:"id"`
	IssuanceId   uint                `json:"issuance_id"`
	Investor     *user.UserOutputDTO `json:"investor"`
	Amount       *uint256.Int        `json:"amount"`
	InterestRate *uint256.Int        `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
}

type CreateOrderUseCase struct {
	UserRepository     repository.UserRepository
	OrderRepository    repository.OrderRepository
	IssuanceRepository repository.IssuanceRepository
}

func NewCreateOrderUseCase(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
	issuanceRepo repository.IssuanceRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		UserRepository:     userRepo,
		OrderRepository:    orderRepo,
		IssuanceRepository: issuanceRepo,
	}
}

func (c *CreateOrderUseCase) Execute(input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit type provided for order creation: %T", deposit)
	}

	issuance, err := c.IssuanceRepository.FindIssuanceById(input.IssuanceId)
	if err != nil {
		return nil, fmt.Errorf("error finding issuance issuances: %w", err)
	}

	if issuance.ClosesAt < metadata.BlockTimestamp {
		return nil, fmt.Errorf("issuance issuance closed, order cannot be placed")
	}

	if types.Address(erc20Deposit.Token) != issuance.Token {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(issuance.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active Issuance max interest rate")
	}

	order, err := entity.NewOrder(
		issuance.Id,
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

	investor, err := c.UserRepository.FindUserByAddress(res.InvestorAddress)
	if err != nil {
		return nil, err
	}

	return &CreateOrderOutputDTO{
		Id:         res.Id,
		IssuanceId: res.IssuanceId,
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
