package crowdfunding_usecase

import (
	"context"
	"fmt"
	"sort"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type CloseCrowdfundingInputDTO struct {
	Creator Address `json:"creator"`
}

type CloseCrowdfundingOutputDTO struct {
	Id                  uint            `json:"id"`
	Token               Address         `json:"token,omitempty"`
	Collateral          *uint256.Int    `json:"collateral,omitempty"`
	Creator             Address         `json:"creator,omitempty"`
	DebtIssued          *uint256.Int    `json:"debt_issued,omitempty"`
	MaxInterestRate     *uint256.Int    `json:"max_interest_rate,omitempty"`
	TotalObligation     *uint256.Int    `json:"total_obligation,omitempty"`
	Orders              []*entity.Order `json:"orders,omitempty"`
	State               string          `json:"state,omitempty"`
	FundraisingDuration int64           `json:"fundraising_duration,omitempty"`
	ClosesAt            int64           `json:"closes_at,omitempty"`
	MaturityAt          int64           `json:"maturity_at,omitempty"`
	CreatedAt           int64           `json:"created_at,omitempty"`
	UpdatedAt           int64           `json:"updated_at,omitempty"`
}

type CloseCrowdfundingUseCase struct {
	OrderRepository        repository.OrderRepository
	CrowdfundingRepository repository.CrowdfundingRepository
}

func NewCloseCrowdfundingUseCase(crowdfundingRepository repository.CrowdfundingRepository, orderRepository repository.OrderRepository) *CloseCrowdfundingUseCase {
	return &CloseCrowdfundingUseCase{
		OrderRepository:        orderRepository,
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (u *CloseCrowdfundingUseCase) Execute(ctx context.Context, input *CloseCrowdfundingInputDTO, metadata rollmelette.Metadata) (*CloseCrowdfundingOutputDTO, error) {
	crowdfundings, err := u.CrowdfundingRepository.FindCrowdfundingsByCreator(ctx, input.Creator)
	if err != nil {
		return nil, err
	}

	var ongoingCrowdfunding *entity.Crowdfunding
	for _, crowdfunding := range crowdfundings {
		if crowdfunding.State == entity.CrowdfundingStateOngoing {
			ongoingCrowdfunding = crowdfunding
			break
		}
	}

	if ongoingCrowdfunding == nil {
		return nil, fmt.Errorf("no ongoing crowdfunding found, cannot close it")
	}

	// Ensure crowdfunding has expired before closing
	if metadata.BlockTimestamp < ongoingCrowdfunding.ClosesAt {
		return nil, fmt.Errorf("crowdfunding not expired yet, you can't close it")
	}

	// Retrieve all orders related to the crowdfunding
	orders, err := u.OrderRepository.FindOrdersByCrowdfundingId(ctx, ongoingCrowdfunding.Id)
	if err != nil {
		return nil, err
	}

	// Sort orders by InterestRate ascending, Amount descending
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].InterestRate.Cmp(orders[j].InterestRate) == 0 {
			return orders[i].Amount.Cmp(orders[j].Amount) > 0 // larger amounts first
		}
		return orders[i].InterestRate.Cmp(orders[j].InterestRate) < 0
	})

	// Reuse variables to reduce allocations
	debtIssuedRemaining := new(uint256.Int).Set(ongoingCrowdfunding.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)
	twoThirdsTarget := uint256.NewInt(0)
	acceptedAmount := uint256.NewInt(0)
	rejectedAmount := uint256.NewInt(0)

	for _, order := range orders {
		if debtIssuedRemaining.IsZero() {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}
			continue
		}

		if debtIssuedRemaining.Gt(order.Amount) || debtIssuedRemaining.Eq(order.Amount) {
			order.State = entity.OrderStateAccepted
			order.UpdatedAt = metadata.BlockTimestamp
			totalCollected.Add(totalCollected, order.Amount)

			// Calculate interest and obligation for this order
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			orderObligation := new(uint256.Int).Add(order.Amount, interest)
			totalObligation.Add(totalObligation, orderObligation)

			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.Sub(debtIssuedRemaining, order.Amount)
		} else {
			// Partially accept order
			acceptedAmount.Set(debtIssuedRemaining)
			rejectedAmount.Sub(order.Amount, acceptedAmount)

			order.Amount = acceptedAmount
			order.State = entity.OrderStatePartiallyAccepted
			order.UpdatedAt = metadata.BlockTimestamp
			totalCollected.Add(totalCollected, acceptedAmount)

			// Calculate interest and obligation for the accepted portion
			interest := new(uint256.Int).Mul(acceptedAmount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			orderObligation := new(uint256.Int).Add(acceptedAmount, interest)
			totalObligation.Add(totalObligation, orderObligation)

			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}

			// Create rejected part
			_, err = u.OrderRepository.CreateOrder(ctx, &entity.Order{
				CrowdfundingId: order.CrowdfundingId,
				Investor:       order.Investor,
				Amount:         rejectedAmount,
				InterestRate:   order.InterestRate,
				State:          entity.OrderStateRejected,
				CreatedAt:      metadata.BlockTimestamp,
				UpdatedAt:      metadata.BlockTimestamp,
			})
			if err != nil {
				return nil, err
			}

			debtIssuedRemaining.Clear()
		}
	}

	// Check funding threshold
	twoThirdsTarget.Mul(ongoingCrowdfunding.DebtIssued, uint256.NewInt(2))
	twoThirdsTarget.Div(twoThirdsTarget, uint256.NewInt(3))
	if totalCollected.Lt(twoThirdsTarget) {
		// Cancel crowdfunding and mark all orders as rejected
		for _, order := range orders {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			_, err := u.OrderRepository.UpdateOrder(ctx, order)
			if err != nil {
				return nil, err
			}
		}

		ongoingCrowdfunding.State = entity.CrowdfundingStateCanceled
		ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
		_, err := u.CrowdfundingRepository.UpdateCrowdfunding(ctx, ongoingCrowdfunding)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("crowdfunding canceled due to insufficient funds collected")
	}

	// Update crowdfunding state
	ongoingCrowdfunding.State = entity.CrowdfundingStateClosed
	ongoingCrowdfunding.TotalObligation = totalObligation
	ongoingCrowdfunding.UpdatedAt = metadata.BlockTimestamp
	res, err := u.CrowdfundingRepository.UpdateCrowdfunding(ctx, ongoingCrowdfunding)
	if err != nil {
		return nil, err
	}

	return &CloseCrowdfundingOutputDTO{
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
