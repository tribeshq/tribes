package crowdfunding_usecase

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateCrowdfundingInputDTO struct {
	DebtIssued          *uint256.Int `json:"debt_issued"`
	MaxInterestRate     *uint256.Int `json:"max_interest_rate"`
	FundraisingDuration int64        `json:"fundraising_duration"`
	ClosesAt            int64        `json:"closes_at"`
	MaturityAt          int64        `json:"maturity_at"`
	Proof               string       `json:"proof"`
}

type CreateCrowdfundingOutputDTO struct {
	Id                  uint            `json:"id"`
	Token               Address         `json:"token,omitempty"`
	Collateral          *uint256.Int    `json:"collateral,omitempty"`
	Creator             Address         `json:"creator,omitempty"`
	DebtIssued          *uint256.Int    `json:"debt_issued"`
	MaxInterestRate     *uint256.Int    `json:"max_interest_rate"`
	Orders              []*entity.Order `json:"orders"`
	State               string          `json:"state"`
	FundraisingDuration int64           `json:"fundraising_duration"`
	ClosesAt            int64           `json:"closes_at"`
	MaturityAt          int64           `json:"maturity_at"`
	CreatedAt           int64           `json:"created_at"`
}

type CreateCrowdfundingUseCase struct {
	UserRepository          repository.UserRepository
	ContractRepository      repository.ContractRepository
	SocialAccountRepository repository.SocialAccountRepository
	CrowdfundingRepository  repository.CrowdfundingRepository
}

func NewCreateCrowdfundingUseCase(
	userRepository repository.UserRepository,
	contractRepository repository.ContractRepository,
	socialRepository repository.SocialAccountRepository,
	crowdfundingRepository repository.CrowdfundingRepository,
) *CreateCrowdfundingUseCase {
	return &CreateCrowdfundingUseCase{
		UserRepository:          userRepository,
		ContractRepository:      contractRepository,
		SocialAccountRepository: socialRepository,
		CrowdfundingRepository:  crowdfundingRepository,
	}
}

func (c *CreateCrowdfundingUseCase) Execute(ctx context.Context, input *CreateCrowdfundingInputDTO, deposit rollmelette.Deposit, metadata rollmelette.Metadata) (*CreateCrowdfundingOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	if input.DebtIssued.Cmp(uint256.NewInt(15000000)) > 0 {
		return nil, fmt.Errorf("%w: debt issued exceeds the maximum allowed value", entity.ErrInvalidCrowdfunding)
	}
	if input.ClosesAt > metadata.BlockTimestamp+15552000 {
		return nil, fmt.Errorf("%w: close date cannot be greater than 6 months", entity.ErrInvalidCrowdfunding)
	}
	if input.ClosesAt > input.MaturityAt {
		return nil, fmt.Errorf("%w: close date cannot be greater than maturity date", entity.ErrInvalidCrowdfunding)
	}
	if metadata.BlockTimestamp >= input.ClosesAt {
		return nil, fmt.Errorf("%w: creation date cannot be greater than or equal to close date", entity.ErrInvalidCrowdfunding)
	}

	// TODO: Add this when in prod
	// if input.FundraisingDuration < 604800 {
	// 	return nil, fmt.Errorf("%w: fundraising duration must be at least 7 days", entity.ErrInvalidCrowdfunding)
	// }
	// if (metadata.BlockTimestamp-input.FundraisingDuration)-metadata.BlockTimestamp < 604800 {
	// 	return nil, fmt.Errorf("%w: cannot create crowndfunding campaign without at least 7 days for the approval process", entity.ErrInvalidCrowdfunding)
	// }

	creator, err := c.UserRepository.FindUserByAddress(ctx, Address(erc20Deposit.Sender))
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}
	if creator.DebtIssuanceLimit.Cmp(input.DebtIssued) < 0 {
		return nil, fmt.Errorf("creator's debt issuance limit exceeded")
	}

	if _, err = c.ContractRepository.FindContractByAddress(ctx, Address(erc20Deposit.Token)); err != nil {
		return nil, fmt.Errorf("unknown token: %w", err)
	}

	crowdfundings, err := c.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, creator.Address)
	if err != nil {
		return nil, fmt.Errorf("error retrieving crowdfundings: %w", err)
	}
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State != entity.CrowdfundingStateSettled && metadata.BlockTimestamp-crowdfunding.CreatedAt < 120*24*60*60 {
			return nil, fmt.Errorf("active crowdfunding exists within the last 120 days")
		}
	}

	// TODO: Change this when in prod
	mockAccount, err := entity.NewSocialAccount(creator.Id, "vitalik", 1000, "twitter", metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	if _, err = c.SocialAccountRepository.CreateSocialAccount(ctx, mockAccount); err != nil {
		return nil, err
	}

	crowdfunding, err := entity.NewCrowdfunding(
		Address(erc20Deposit.Token),
		uint256.MustFromBig(erc20Deposit.Amount),
		creator.Address,
		input.DebtIssued,
		input.MaxInterestRate,
		input.FundraisingDuration,
		input.ClosesAt,
		input.MaturityAt,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating crowdfunding: %w", err)
	}
	createdCrowdfunding, err := c.CrowdfundingRepository.CreateCrowdfunding(ctx, crowdfunding)
	if err != nil {
		return nil, fmt.Errorf("error creating crowdfunding: %w", err)
	}

	creator.DebtIssuanceLimit.Sub(creator.DebtIssuanceLimit, input.DebtIssued)
	if _, err = c.UserRepository.UpdateUser(ctx, creator); err != nil {
		return nil, fmt.Errorf("error updating creator debt issuance limit: %w", err)
	}

	return &CreateCrowdfundingOutputDTO{
		Id:                  createdCrowdfunding.Id,
		Token:               createdCrowdfunding.Token,
		Collateral:          createdCrowdfunding.Collateral,
		Creator:             createdCrowdfunding.Creator,
		DebtIssued:          createdCrowdfunding.DebtIssued,
		MaxInterestRate:     createdCrowdfunding.MaxInterestRate,
		Orders:              createdCrowdfunding.Orders,
		State:               string(createdCrowdfunding.State),
		FundraisingDuration: createdCrowdfunding.FundraisingDuration,
		ClosesAt:            createdCrowdfunding.ClosesAt,
		MaturityAt:          createdCrowdfunding.MaturityAt,
		CreatedAt:           createdCrowdfunding.CreatedAt,
	}, nil
}
