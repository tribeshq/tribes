package auction

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateAuctionInputDTO struct {
	Token           Address      `json:"token" validate:"required"`
	DebtIssued      *uint256.Int `json:"debt_issued" validate:"required"`
	MaxInterestRate *uint256.Int `json:"max_interest_rate" validate:"required"`
	ClosesAt        int64        `json:"closes_at" validate:"required"`
	MaturityAt      int64        `json:"maturity_at" validate:"required"`
}

type CreateAuctionOutputDTO struct {
	Id                uint            `json:"id"`
	Token             Address         `json:"token,omitempty"`
	Creator           Address         `json:"creator,omitempty"`
	CollateralAddress Address         `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int    `json:"collateral_amount,omitempty"`
	DebtIssued        *uint256.Int    `json:"debt_issued"`
	MaxInterestRate   *uint256.Int    `json:"max_interest_rate"`
	State             string          `json:"state"`
	Orders            []*entity.Order `json:"orders"`
	CreatedAt         int64           `json:"created_at"`
	ClosesAt          int64           `json:"closes_at"`
	MaturityAt        int64           `json:"maturity_at"`
}

type CreateAuctionUseCase struct {
	AuctionRepository repository.AuctionRepository
	UserRepository    repository.UserRepository
}

func NewCreateAuctionUseCase(
	AuctionRepository repository.AuctionRepository,
	UserRepository repository.UserRepository,
) *CreateAuctionUseCase {
	return &CreateAuctionUseCase{
		AuctionRepository: AuctionRepository,
		UserRepository:    UserRepository,
	}
}

func (c *CreateAuctionUseCase) Execute(ctx context.Context, input *CreateAuctionInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateAuctionOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	user, err := c.UserRepository.FindUserByAddress(ctx, Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if err := c.Validate(user, input, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	auctions, err := c.AuctionRepository.FindAuctionsByCreator(ctx, Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error retrieving Auctions: %w", err)
	}
	for _, auction := range auctions {
		if auction.State != entity.AuctionStateSettled && auction.State != entity.AuctionStateCollateralExecuted {
			return nil, fmt.Errorf("active auction exists, cannot create a new auction")
		}
	}

	Auction, err := entity.NewAuction(
		input.Token,
		Address(erc20Deposit.Sender),
		Address(erc20Deposit.Token),
		uint256.MustFromBig(erc20Deposit.Value),
		input.DebtIssued,
		input.MaxInterestRate,
		input.ClosesAt,
		input.MaturityAt,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Auction: %w", err)
	}

	createdAuction, err := c.AuctionRepository.CreateAuction(ctx, Auction)
	if err != nil {
		return nil, fmt.Errorf("error creating Auction: %w", err)
	}

	return &CreateAuctionOutputDTO{
		Id:                createdAuction.Id,
		Token:             createdAuction.Token,
		Creator:           createdAuction.Creator,
		CollateralAddress: createdAuction.CollateralAddress,
		CollateralAmount:  createdAuction.CollateralAmount,
		DebtIssued:        createdAuction.DebtIssued,
		MaxInterestRate:   createdAuction.MaxInterestRate,
		Orders:            createdAuction.Orders,
		State:             string(createdAuction.State),
		ClosesAt:          createdAuction.ClosesAt,
		MaturityAt:        createdAuction.MaturityAt,
		CreatedAt:         createdAuction.CreatedAt,
	}, nil
}

func (c *CreateAuctionUseCase) Validate(
	user *entity.User,
	input *CreateAuctionInputDTO,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if len(user.SocialAccounts) == 0 {
		return fmt.Errorf("%w: user has no social accounts, please verify at least one social account", entity.ErrInvalidAuction)
	}

	if input.DebtIssued.Cmp(uint256.NewInt(15000000)) > 0 {
		return fmt.Errorf("%w: debt issued exceeds the maximum allowed value", entity.ErrInvalidAuction)
	}

	if input.ClosesAt > metadata.BlockTimestamp+180*24*60*60 {
		return fmt.Errorf("%w: close date cannot be greater than 180 days", entity.ErrInvalidAuction)
	}

	if input.ClosesAt > input.MaturityAt {
		return fmt.Errorf("%w: close date cannot be greater than maturity date", entity.ErrInvalidAuction)
	}

	if metadata.BlockTimestamp >= input.ClosesAt {
		return fmt.Errorf("%w: creation date cannot be greater than or equal to close date", entity.ErrInvalidAuction)
	}
	return nil
}
