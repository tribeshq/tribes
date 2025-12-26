package issuance

import (
	"fmt"
	"sort"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CloseIssuanceInputDTO struct {
	CreatorAddress types.Address `json:"creator_address" validate:"required"`
}

type CloseIssuanceOutputDTO struct {
	Id                uint                    `json:"id"`
	Title             string                  `json:"title,omitempty"`
	Description       string                  `json:"description,omitempty"`
	Promotion         string                  `json:"promotion,omitempty"`
	Token             types.Address           `json:"token,omitempty"`
	Creator           *user.UserOutputDTO     `json:"creator,omitempty"`
	CollateralAddress types.Address           `json:"collateral,omitempty"`
	CollateralAmount  *uint256.Int            `json:"collateral_amount,omitempty"`
	BadgeAddress      types.Address           `json:"badge_address,omitempty"`
	DebtIssued        *uint256.Int            `json:"debt_issued,omitempty"`
	MaxInterestRate   *uint256.Int            `json:"max_interest_rate,omitempty"`
	TotalObligation   *uint256.Int            `json:"total_obligation,omitempty"`
	TotalRaised       *uint256.Int            `json:"total_raised,omitempty"`
	State             string                  `json:"state,omitempty"`
	Orders            []*order.OrderOutputDTO `json:"orders,omitempty"`
	CreatedAt         int64                   `json:"created_at,omitempty"`
	ClosesAt          int64                   `json:"closes_at,omitempty"`
	MaturityAt        int64                   `json:"maturity_at,omitempty"`
	UpdatedAt         int64                   `json:"updated_at,omitempty"`
}

type CloseIssuanceUseCase struct {
	UserRepository     repository.UserRepository
	OrderRepository    repository.OrderRepository
	IssuanceRepository repository.IssuanceRepository
}

func NewCloseIssuanceUseCase(userRepo repository.UserRepository, issuanceRepo repository.IssuanceRepository, orderRepo repository.OrderRepository) *CloseIssuanceUseCase {
	return &CloseIssuanceUseCase{
		UserRepository:     userRepo,
		IssuanceRepository: issuanceRepo,
		OrderRepository:    orderRepo,
	}
}

func (u *CloseIssuanceUseCase) Execute(input *CloseIssuanceInputDTO, metadata rollmelette.Metadata) (*CloseIssuanceOutputDTO, error) {
	// -------------------------------------------------------------------------
	// 1. Find ongoing issuance for the creator
	// -------------------------------------------------------------------------
	issuances, err := u.IssuanceRepository.FindIssuancesByCreatorAddress(input.CreatorAddress)
	if err != nil {
		return nil, err
	}
	var ongoingIssuance *entity.Issuance
	for _, issuance := range issuances {
		if issuance.State == entity.IssuanceStateOngoing {
			ongoingIssuance = issuance
			break
		}
	}
	if ongoingIssuance == nil {
		return nil, fmt.Errorf("no ongoing issuance found, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 2. Validate issuance expiration
	// -------------------------------------------------------------------------
	if metadata.BlockTimestamp < ongoingIssuance.ClosesAt {
		return nil, fmt.Errorf("issuance not expired yet, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 3. Fetch and sort issuance orders
	// -------------------------------------------------------------------------
	orders, err := u.OrderRepository.FindOrdersByIssuanceId(ongoingIssuance.Id)
	if err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		cmp := orders[i].InterestRate.Cmp(orders[j].InterestRate)
		if cmp == 0 {
			return orders[i].Amount.Cmp(orders[j].Amount) > 0
		}
		return cmp < 0
	})

	// -------------------------------------------------------------------------
	// 4. Select winning orders and calculate obligations
	// -------------------------------------------------------------------------
	debtRemaining := new(uint256.Int).Set(ongoingIssuance.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)

	for _, order := range orders {
		if debtRemaining.IsZero() {
			// Reject surplus orders
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(order); err != nil {
				return nil, err
			}
			continue
		}

		// Accept full or partial order
		acceptAmount := order.Amount
		if debtRemaining.Lt(order.Amount) {
			acceptAmount = new(uint256.Int).Set(debtRemaining)
		}
		interest := new(uint256.Int).Mul(acceptAmount, order.InterestRate)
		interest.Div(interest, uint256.NewInt(100))

		orderObligation := new(uint256.Int).Add(acceptAmount, interest)
		totalCollected.Add(totalCollected, acceptAmount)
		totalObligation.Add(totalObligation, orderObligation)

		if debtRemaining.Cmp(order.Amount) >= 0 {
			order.State = entity.OrderStateAccepted
			debtRemaining.Sub(debtRemaining, order.Amount)
		} else {
			order.State = entity.OrderStatePartiallyAccepted
			// Create rejected order for the surplus
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptAmount)
			_, err := u.OrderRepository.CreateOrder(&entity.Order{
				IssuanceId:      order.IssuanceId,
				InvestorAddress: order.InvestorAddress,
				Amount:          rejectedAmount,
				InterestRate:    order.InterestRate,
				State:           entity.OrderStateRejected,
				CreatedAt:       order.CreatedAt,
				UpdatedAt:       metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}
			debtRemaining.Clear()
		}
		order.Amount = acceptAmount
		order.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.OrderRepository.UpdateOrder(order); err != nil {
			return nil, err
		}
	}

	// -------------------------------------------------------------------------
	// 5. Check if minimum funding (2/3) was reached
	// -------------------------------------------------------------------------
	twoThirds := new(uint256.Int).Mul(ongoingIssuance.DebtIssued, uint256.NewInt(2))
	twoThirds.Div(twoThirds, uint256.NewInt(3))
	if totalCollected.Lt(twoThirds) {
		// Cancel issuance and reject all orders
		for _, order := range orders {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			if _, err := u.OrderRepository.UpdateOrder(order); err != nil {
				return nil, err
			}
		}
		ongoingIssuance.State = entity.IssuanceStateCanceled
		ongoingIssuance.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.IssuanceRepository.UpdateIssuance(ongoingIssuance); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("issuance canceled due to insufficient funds collected, expected at least 2/3 of the debt issued: %s, got: %s", twoThirds.String(), totalCollected.String())
	}

	// -------------------------------------------------------------------------
	// 6. Close issuance and return result
	// -------------------------------------------------------------------------
	ongoingIssuance.State = entity.IssuanceStateClosed
	ongoingIssuance.TotalObligation = totalObligation
	ongoingIssuance.TotalRaised = totalCollected
	ongoingIssuance.UpdatedAt = metadata.BlockTimestamp
	res, err := u.IssuanceRepository.UpdateIssuance(ongoingIssuance)
	if err != nil {
		return nil, err
	}

	creator, err := u.UserRepository.FindUserByAddress(res.CreatorAddress)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	orderDTOs := make([]*order.OrderOutputDTO, len(res.Orders))
	for i, o := range res.Orders {
		investor, err := u.UserRepository.FindUserByAddress(o.InvestorAddress)
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

	return &CloseIssuanceOutputDTO{
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
		Orders:            orderDTOs,
		State:             string(res.State),
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
