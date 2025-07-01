package campaign

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateCampaignInputDTO struct {
	Token           custom_type.Address `json:"token" validate:"required"`
	DebtIssued      *uint256.Int        `json:"debt_issued" validate:"required"`
	MaxInterestRate *uint256.Int        `json:"max_interest_rate" validate:"required"`
	BadgeRouter     custom_type.Address `json:"badge_router" validate:"required"`
	BadgeMinter     custom_type.Address `json:"badge_minter" validate:"required"`
	ClosesAt        int64               `json:"closes_at" validate:"required"`
	MaturityAt      int64               `json:"maturity_at" validate:"required"`
}

type CreateCampaignOutputDTO struct {
	Id                uint                `json:"id"`
	Token             custom_type.Address `json:"token,omitempty"`
	Creator           custom_type.Address `json:"creator,omitempty"`
	CollateralAddress custom_type.Address `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int        `json:"collateral_amount,omitempty"`
	BadgeRouter       custom_type.Address `json:"badge_router,omitempty"`
	BadgeMinter       custom_type.Address `json:"badge_minter,omitempty"`
	DebtIssued        *uint256.Int        `json:"debt_issued"`
	MaxInterestRate   *uint256.Int        `json:"max_interest_rate"`
	State             string              `json:"state"`
	Orders            []*entity.Order     `json:"orders"`
	CreatedAt         int64               `json:"created_at"`
	ClosesAt          int64               `json:"closes_at"`
	MaturityAt        int64               `json:"maturity_at"`
}

type CreateCampaignUseCase struct {
	CampaignRepository repository.CampaignRepository
	UserRepository     repository.UserRepository
}

func NewCreateCampaignUseCase(
	CampaignRepository repository.CampaignRepository,
	UserRepository repository.UserRepository,
) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{
		CampaignRepository: CampaignRepository,
		UserRepository:     UserRepository,
	}
}

func (c *CreateCampaignUseCase) Execute(ctx context.Context, input *CreateCampaignInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateCampaignOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	user, err := c.UserRepository.FindUserByAddress(ctx, custom_type.Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if err := c.Validate(user, input, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	campaigns, err := c.CampaignRepository.FindCampaignsByCreator(ctx, custom_type.Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error retrieving Campaigns: %w", err)
	}
	for _, campaign := range campaigns {
		if campaign.State != entity.CampaignStateSettled && campaign.State != entity.CampaignStateCollateralExecuted {
			return nil, fmt.Errorf("active campaign exists, cannot create a new campaign")
		}
	}

	Campaign, err := entity.NewCampaign(
		input.Token,
		custom_type.Address(erc20Deposit.Sender),
		custom_type.Address(erc20Deposit.Token),
		uint256.MustFromBig(erc20Deposit.Value),
		input.BadgeRouter,
		input.BadgeMinter,
		input.DebtIssued,
		input.MaxInterestRate,
		input.ClosesAt,
		input.MaturityAt,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Campaign: %w", err)
	}

	createdCampaign, err := c.CampaignRepository.CreateCampaign(ctx, Campaign)
	if err != nil {
		return nil, fmt.Errorf("error creating Campaign: %w", err)
	}

	return &CreateCampaignOutputDTO{
		Id:                createdCampaign.Id,
		Token:             createdCampaign.Token,
		Creator:           createdCampaign.Creator,
		CollateralAddress: createdCampaign.CollateralAddress,
		CollateralAmount:  createdCampaign.CollateralAmount,
		BadgeRouter:       createdCampaign.BadgeRouter,
		BadgeMinter:       createdCampaign.BadgeMinter,
		DebtIssued:        createdCampaign.DebtIssued,
		MaxInterestRate:   createdCampaign.MaxInterestRate,
		Orders:            createdCampaign.Orders,
		State:             string(createdCampaign.State),
		ClosesAt:          createdCampaign.ClosesAt,
		MaturityAt:        createdCampaign.MaturityAt,
		CreatedAt:         createdCampaign.CreatedAt,
	}, nil
}

func (c *CreateCampaignUseCase) Validate(
	user *entity.User,
	input *CreateCampaignInputDTO,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if len(user.SocialAccounts) == 0 {
		return fmt.Errorf("%w: user has no social accounts, please verify at least one social account", entity.ErrInvalidCampaign)
	}

	if input.ClosesAt > metadata.BlockTimestamp+180*24*60*60 {
		return fmt.Errorf("%w: close date cannot be greater than 180 days", entity.ErrInvalidCampaign)
	}

	if input.ClosesAt > input.MaturityAt {
		return fmt.Errorf("%w: close date cannot be greater than maturity date", entity.ErrInvalidCampaign)
	}

	if metadata.BlockTimestamp >= input.ClosesAt {
		return fmt.Errorf("%w: creation date cannot be greater than or equal to close date", entity.ErrInvalidCampaign)
	}
	return nil
}
