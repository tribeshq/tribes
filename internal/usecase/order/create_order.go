package order

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateOrderInputDTO struct {
	AuctionId    uint         `json:"auction_id" validate:"required"`
	InterestRate *uint256.Int `json:"interest_rate" validate:"required"`
}

type CreateOrderOutputDTO struct {
	Id           uint         `json:"id"`
	AuctionId    uint         `json:"auction_id"`
	Investor     Address      `json:"investor"`
	Amount       *uint256.Int `json:"amount"`
	InterestRate *uint256.Int `json:"interest_rate"`
	State        string       `json:"state"`
	CreatedAt    int64        `json:"created_at"`
}

type CreateOrderUseCase struct {
	OrderRepository   repository.OrderRepository
	UserRepository    repository.UserRepository
	AuctionRepository repository.AuctionRepository
}

func NewCreateOrderUseCase(orderRepository repository.OrderRepository, userRepository repository.UserRepository, auctionRepository repository.AuctionRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:   orderRepository,
		UserRepository:    userRepository,
		AuctionRepository: auctionRepository,
	}
}

func (c *CreateOrderUseCase) Execute(ctx context.Context, input *CreateOrderInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateOrderOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type provided for order creation: %T", deposit)
	}

	// check if user has reached the investment limit
	user, err := c.UserRepository.FindUserByAddress(ctx, Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user.InvestmentLimit.Lt(uint256.MustFromBig(erc20Deposit.Value)) {
		return nil, fmt.Errorf("investor has reached the investment limit")
	}

	auction, err := c.AuctionRepository.FindAuctionById(ctx, input.AuctionId)
	if err != nil {
		return nil, fmt.Errorf("error finding auction campaigns: %w", err)
	}

	if auction.ClosesAt < metadata.BlockTimestamp {
		return nil, fmt.Errorf("auction campaign closed, order cannot be placed")
	}

	if Address(erc20Deposit.Token) != auction.Token {
		return nil, fmt.Errorf("invalid contract address provided for order creation: %v", erc20Deposit.Token)
	}

	if input.InterestRate.Gt(auction.MaxInterestRate) {
		return nil, fmt.Errorf("order interest rate exceeds active Auction max interest rate")
	}

	order, err := entity.NewOrder(
		auction.Id,
		Address(erc20Deposit.Sender),
		uint256.MustFromBig(erc20Deposit.Value),
		input.InterestRate,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, err
	}

	res, err := c.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	user.InvestmentLimit.Sub(user.InvestmentLimit, uint256.MustFromBig(erc20Deposit.Value))
	_, err = c.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &CreateOrderOutputDTO{
		Id:           res.Id,
		AuctionId:    res.AuctionId,
		Investor:     res.Investor,
		Amount:       res.Amount,
		InterestRate: res.InterestRate,
		State:        string(res.State),
		CreatedAt:    res.CreatedAt,
	}, nil
}
