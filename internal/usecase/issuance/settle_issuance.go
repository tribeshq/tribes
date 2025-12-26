package issuance

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	. "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type SettleIssuanceInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type SettleIssuanceOutputDTO struct {
	Id                uint                    `json:"id"`
	Title             string                  `json:"title,omitempty"`
	Description       string                  `json:"description,omitempty"`
	Promotion         string                  `json:"promotion,omitempty"`
	Token             Address                 `json:"token"`
	Creator           *user.UserOutputDTO     `json:"creator"`
	CollateralAddress Address                 `json:"collateral"`
	CollateralAmount  *uint256.Int            `json:"collateral_amount"`
	BadgeAddress      Address                 `json:"badge_address"`
	DebtIssued        *uint256.Int            `json:"debt_issued"`
	MaxInterestRate   *uint256.Int            `json:"max_interest_rate"`
	TotalObligation   *uint256.Int            `json:"total_obligation"`
	TotalRaised       *uint256.Int            `json:"total_raised"`
	State             string                  `json:"state"`
	Orders            []*order.OrderOutputDTO `json:"orders"`
	CreatedAt         int64                   `json:"created_at"`
	ClosesAt          int64                   `json:"closes_at"`
	MaturityAt        int64                   `json:"maturity_at"`
	UpdatedAt         int64                   `json:"updated_at"`
}

type SettleIssuanceUseCase struct {
	UserRepository     repository.UserRepository
	IssuanceRepository repository.IssuanceRepository
	OrderRepository    repository.OrderRepository
}

func NewSettleIssuanceUseCase(
	UserRepository repository.UserRepository,
	IssuanceRepository repository.IssuanceRepository,
	OrderRepository repository.OrderRepository,
) *SettleIssuanceUseCase {
	return &SettleIssuanceUseCase{
		UserRepository:     UserRepository,
		IssuanceRepository: IssuanceRepository,
		OrderRepository:    OrderRepository,
	}
}

func (uc *SettleIssuanceUseCase) Execute(
	input *SettleIssuanceInputDTO,
	deposit rollmelette.Deposit,
	metadata rollmelette.Metadata,
) (*SettleIssuanceOutputDTO, error) {
	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return nil, fmt.Errorf("invalid deposit types: %T", deposit)
	}

	issuance, err := uc.IssuanceRepository.FindIssuanceById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("error finding issuance: %w", err)
	}

	if err := uc.Validate(issuance, erc20Deposit, metadata); err != nil {
		return nil, err
	}

	for _, order := range issuance.Orders {
		if order.State == entity.OrderStateAccepted || order.State == entity.OrderStatePartiallyAccepted {
			order.State = entity.OrderStateSettled
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := uc.OrderRepository.UpdateOrder(order); err != nil {
				return nil, fmt.Errorf("error updating order: %w", err)
			}
		}
	}

	issuance.State = entity.IssuanceStateSettled
	issuance.UpdatedAt = metadata.BlockTimestamp
	
	res, err := uc.IssuanceRepository.UpdateIssuance(issuance)
	if err != nil {
		return nil, fmt.Errorf("error updating issuance: %w", err)
	}

	creator, err := uc.UserRepository.FindUserByAddress(res.CreatorAddress)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	orderDTOs := make([]*order.OrderOutputDTO, len(res.Orders))
	for i, o := range res.Orders {
		investor, err := uc.UserRepository.FindUserByAddress(o.InvestorAddress)
		if err != nil {
			return nil, fmt.Errorf("error finding investor: %w", err)
		}
		orderDTOs[i] = &order.OrderOutputDTO{
			Id:         o.Id,
			IssuanceId: o.IssuanceId,
			Investor: &user.UserOutputDTO{
				Id:             investor.Id,
				Role:           string(investor.Role),
				Address:        investor.Address,
				SocialAccounts: investor.SocialAccounts,
				CreatedAt:      investor.CreatedAt,
				UpdatedAt:      investor.UpdatedAt,
			},
			Amount:       o.Amount,
			InterestRate: o.InterestRate,
			State:        string(o.State),
			CreatedAt:    o.CreatedAt,
			UpdatedAt:    o.UpdatedAt,
		}
	}

	return &SettleIssuanceOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Promotion:   res.Promotion,
		Token:       res.Token,
		Creator: &user.UserOutputDTO{
			Id:             creator.Id,
			Role:           string(creator.Role),
			Address:        creator.Address,
			SocialAccounts: creator.SocialAccounts,
			CreatedAt:      creator.CreatedAt,
			UpdatedAt:      creator.UpdatedAt,
		},
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		BadgeAddress:      res.BadgeAddress,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		State:             string(res.State),
		Orders:            orderDTOs,
		CreatedAt:         res.CreatedAt,
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}

func (uc *SettleIssuanceUseCase) Validate(
	Issuance *entity.Issuance,
	deposit *rollmelette.ERC20Deposit,
	metadata rollmelette.Metadata,
) error {
	if metadata.BlockTimestamp > Issuance.MaturityAt {
		return fmt.Errorf("the maturity date of the issuance issuance has passed")
	}

	if Issuance.State == entity.IssuanceStateSettled {
		return fmt.Errorf("issuance issuance already settled")
	}

	if Issuance.State != entity.IssuanceStateClosed {
		return fmt.Errorf("issuance issuance not closed")
	}

	if deposit.Value.Cmp(Issuance.TotalObligation.ToBig()) < 0 {
		return fmt.Errorf("deposit amount is lower than the total obligation")
	}

	if Issuance.CreatorAddress != Address(deposit.Sender) {
		return fmt.Errorf("only the issuance creator can settle the issuance")
	}
	return nil
}
