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

type SettleCrowdfundingInputDTO struct {
	CrowdfundingId uint `json:"crowdfunding_id"`
}

type SettleCrowdfundingOutputDTO struct {
	Id                  uint            `json:"id"`
	Token               Address         `json:"token"`
	Collateral          *uint256.Int    `json:"collateral"`
	Creator             Address         `json:"creator"`
	DebtIssued          *uint256.Int    `json:"debt_issued"`
	MaxInterestRate     *uint256.Int    `json:"max_interest_rate"`
	TotalObligation     *uint256.Int    `json:"total_obligation"`
	Orders              []*entity.Order `json:"orders"`
	State               string          `json:"state"`
	FundraisingDuration int64           `json:"fundraising_duration"`
	ClosesAt            int64           `json:"closes_at"`
	MaturityAt          int64           `json:"maturity_at"`
	CreatedAt           int64           `json:"created_at"`
	UpdatedAt           int64           `json:"updated_at"`
}

type SettleCrowdfundingUseCase struct {
	UserRepository         repository.UserRepository
	ContractRepository     repository.ContractRepository
	CrowdfundingRepository repository.CrowdfundingRepository
	OrderRepository        repository.OrderRepository
}

func NewSettleCrowdfundingUseCase(
	userRepository repository.UserRepository,
	crowdfundingRepository repository.CrowdfundingRepository,
	contractRepository repository.ContractRepository,
	orderRepository repository.OrderRepository,
) *SettleCrowdfundingUseCase {
	return &SettleCrowdfundingUseCase{
		UserRepository:         userRepository,
		ContractRepository:     contractRepository,
		CrowdfundingRepository: crowdfundingRepository,
		OrderRepository:        orderRepository,
	}
}

func (uc *SettleCrowdfundingUseCase) Execute(
	ctx context.Context,
	input *SettleCrowdfundingInputDTO,
	deposit rollmelette.Deposit,
	metadata rollmelette.Metadata,
) (*SettleCrowdfundingOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit custom_type: %T", deposit)
	}

	stablecoin, err := uc.ContractRepository.FindContractBySymbol(ctx, "STABLECOIN")
	if err != nil {
		return nil, fmt.Errorf("error finding stablecoin contract: %w", err)
	}

	if Address(erc20Deposit.Token) != stablecoin.Address {
		return nil, fmt.Errorf("token deposit is not the stablecoin %v", stablecoin.Address)
	}

	crowdfunding, err := uc.CrowdfundingRepository.FindCrowdfundingById(ctx, input.CrowdfundingId)
	if err != nil {
		return nil, fmt.Errorf("error finding crowdfunding campaign: %w", err)
	}

	if err := uc.validate(crowdfunding, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	// Update orders
	for _, order := range crowdfunding.Orders {
		if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
			order.State = entity.OrderStateSettled
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := uc.OrderRepository.UpdateOrder(ctx, order); err != nil {
				return nil, fmt.Errorf("error updating order: %w", err)
			}
		}
	}

	// Update crowdfunding
	crowdfunding.State = entity.CrowdfundingStateSettled
	crowdfunding.UpdatedAt = metadata.BlockTimestamp
	res, err := uc.CrowdfundingRepository.UpdateCrowdfunding(ctx, crowdfunding)
	if err != nil {
		return nil, fmt.Errorf("error updating crowdfunding: %w", err)
	}

	// Update creator
	creator, err := uc.UserRepository.FindUserByAddress(ctx, crowdfunding.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	creator.DebtIssuanceLimit = new(uint256.Int).Sub(creator.DebtIssuanceLimit, crowdfunding.DebtIssued)
	if _, err := uc.UserRepository.UpdateUser(ctx, creator); err != nil {
		return nil, fmt.Errorf("error updating creator debt limit: %w", err)
	}

	return &SettleCrowdfundingOutputDTO{
		Id:                  res.Id,
		Token:               res.Token,
		Collateral:          res.Collateral,
		Creator:             res.Creator,
		DebtIssued:          res.DebtIssued,
		MaxInterestRate:     res.MaxInterestRate,
		TotalObligation:     res.TotalObligation,
		Orders:              res.Orders,
		State:               string(res.State),
		FundraisingDuration: res.FundraisingDuration,
		ClosesAt:            res.ClosesAt,
		MaturityAt:          res.MaturityAt,
		CreatedAt:           res.CreatedAt,
		UpdatedAt:           res.UpdatedAt,
	}, nil
}

func (uc *SettleCrowdfundingUseCase) validate(
	crowdfunding *entity.Crowdfunding,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if metadata.BlockTimestamp > crowdfunding.MaturityAt {
		return fmt.Errorf("the maturity date of the crowdfunding campaign has passed")
	}

	if crowdfunding.State == entity.CrowdfundingStateSettled {
		return fmt.Errorf("crowdfunding campaign already settled")
	}

	if crowdfunding.State != entity.CrowdfundingStateClosed {
		return fmt.Errorf("crowdfunding campaign not closed")
	}

	if deposit.Amount.Cmp(crowdfunding.TotalObligation.ToBig()) < 0 {
		return fmt.Errorf("deposit amount is lower than the total obligation")
	}

	return nil
}
