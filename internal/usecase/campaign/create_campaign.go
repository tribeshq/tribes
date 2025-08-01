package campaign

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/assets"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateCampaignInputDTO struct {
	Title           string              `json:"title" validate:"required,min=3,max=100"`
	Description     string              `json:"description" validate:"required,min=10,max=1000"`
	Promotion       string              `json:"promotion" validate:"required,min=5,max=500"`
	Token           custom_type.Address `json:"token" validate:"required"`
	BadgeName       string              `json:"badge_name" validate:"required,min=3,max=100"`
	BadgeSymbol     string              `json:"badge_symbol" validate:"required,min=3,max=100"`
	DebtIssued      *uint256.Int        `json:"debt_issued" validate:"required"`
	MaxInterestRate *uint256.Int        `json:"max_interest_rate" validate:"required"`
	ClosesAt        int64               `json:"closes_at" validate:"required"`
	MaturityAt      int64               `json:"maturity_at" validate:"required"`
}

type CreateCampaignOutputDTO struct {
	Id                uint                `json:"id"`
	Title             string              `json:"title,omitempty"`
	Description       string              `json:"description,omitempty"`
	Promotion         string              `json:"promotion,omitempty"`
	Token             custom_type.Address `json:"token,omitempty"`
	Creator           *user.UserOutputDTO `json:"creator,omitempty"`
	CollateralAddress custom_type.Address `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int        `json:"collateral_amount,omitempty"`
	BadgeAddress      custom_type.Address `json:"badge_address,omitempty"`
	DebtIssued        *uint256.Int        `json:"debt_issued"`
	MaxInterestRate   *uint256.Int        `json:"max_interest_rate"`
	State             string              `json:"state"`
	Orders            []*entity.Order     `json:"orders"`
	CreatedAt         int64               `json:"created_at"`
	ClosesAt          int64               `json:"closes_at"`
	MaturityAt        int64               `json:"maturity_at"`
}

type CreateCampaignUseCase struct {
	cfg                *configs.RollupConfig
	campaignRepository repository.CampaignRepository
	userRepository     repository.UserRepository
}

func NewCreateCampaignUseCase(
	cfg *configs.RollupConfig,
	campaignRepo repository.CampaignRepository,
	userRepo repository.UserRepository,
) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{
		cfg:                cfg,
		campaignRepository: campaignRepo,
		userRepository:     userRepo,
	}
}

func (c *CreateCampaignUseCase) Execute(input *CreateCampaignInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateCampaignOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	creator, err := c.userRepository.FindUserByAddress(custom_type.Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if err := c.Validate(creator, input, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	bytecode, err := assets.GetBadgeBytecode()
	if err != nil {
		return nil, fmt.Errorf("error getting badge bytecode: %w", err)
	}

	campaigns, err := c.campaignRepository.FindCampaignsByCreatorAddress(custom_type.Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error retrieving Campaigns: %w", err)
	}
	for _, campaign := range campaigns {
		if campaign.State != entity.CampaignStateSettled && campaign.State != entity.CampaignStateCollateralExecuted {
			return nil, fmt.Errorf("active campaign exists, cannot create a new campaign")
		}
	}

	stringType, _ := abi.NewType("string", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
		{Type: stringType},
		{Type: stringType},
	}.Pack(metadata.AppContract, input.BadgeName, input.BadgeSymbol)
	if err != nil {
		return nil, fmt.Errorf("error encoding constructor args: %w", err)
	}

	badgeAddress := crypto.CreateAddress2(
		c.cfg.BadgeFactoryAddress,
		common.HexToHash(strconv.Itoa(metadata.Index)),
		crypto.Keccak256(append(bytecode, constructorArgs...)),
	)

	campaign, err := entity.NewCampaign(
		input.Title,
		input.Description,
		input.Promotion,
		input.Token,
		custom_type.Address(erc20Deposit.Sender),
		custom_type.Address(erc20Deposit.Token),
		uint256.MustFromBig(erc20Deposit.Value),
		custom_type.Address(badgeAddress),
		input.DebtIssued,
		input.MaxInterestRate,
		input.ClosesAt,
		input.MaturityAt,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating Campaign: %w", err)
	}

	createdCampaign, err := c.campaignRepository.CreateCampaign(campaign)
	if err != nil {
		return nil, fmt.Errorf("error creating Campaign: %w", err)
	}

	return &CreateCampaignOutputDTO{
		Id:          createdCampaign.Id,
		Title:       createdCampaign.Title,
		Description: createdCampaign.Description,
		Promotion:   createdCampaign.Promotion,
		Token:       createdCampaign.Token,
		Creator: &user.UserOutputDTO{
			Id:             creator.Id,
			Role:           string(creator.Role),
			Address:        creator.Address,
			SocialAccounts: creator.SocialAccounts,
			CreatedAt:      creator.CreatedAt,
			UpdatedAt:      creator.UpdatedAt,
		},
		CollateralAddress: createdCampaign.CollateralAddress,
		CollateralAmount:  createdCampaign.CollateralAmount,
		BadgeAddress:      createdCampaign.BadgeAddress,
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
