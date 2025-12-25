package entity

import (
	"errors"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
)

var (
	ErrInvalidOrder  = errors.New("invalid order")
	ErrOrderNotFound = errors.New("order not found")
)

type OrderState string

const (
	OrderStatePending             OrderState = "pending"
	OrderStateAccepted            OrderState = "accepted"
	OrderStateCancelled           OrderState = "canceled"
	OrderStatePartiallyAccepted   OrderState = "partially_accepted"
	OrderStateRejected            OrderState = "rejected"
	OrderStateSettled             OrderState = "settled"
	OrderStateSettledByCollateral OrderState = "settled_by_collateral"
)

type Order struct {
	Id              uint          `json:"id" gorm:"primaryKey"`
	CampaignId      uint          `json:"campaign_id" gorm:"not null;index"`
	InvestorAddress types.Address `json:"investor_address,omitempty" gorm:"not null"`
	Amount          *uint256.Int  `json:"amount,omitempty" gorm:"types:text;not null"`
	InterestRate    *uint256.Int  `json:"interest_rate,omitempty" gorm:"types:text;not null"`
	State           OrderState    `json:"state,omitempty" gorm:"types:text;not null"`
	CreatedAt       int64         `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt       int64         `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewOrder(campaignId uint, investorAddress types.Address, amount *uint256.Int, interestRate *uint256.Int, createdAt int64) (*Order, error) {
	order := &Order{
		CampaignId:      campaignId,
		InvestorAddress: investorAddress,
		Amount:          amount,
		InterestRate:    interestRate,
		State:           OrderStatePending,
		CreatedAt:       createdAt,
	}
	if err := order.validate(); err != nil {
		return nil, err
	}
	return order, nil
}

func (b *Order) validate() error {
	if b.CampaignId == 0 {
		return fmt.Errorf("%w: Campaign ID cannot be zero", ErrInvalidOrder)
	}
	if b.InvestorAddress == (types.Address{}) {
		return fmt.Errorf("%w: investor address cannot be empty", ErrInvalidOrder)
	}
	if b.Amount.Sign() <= 0 {
		return fmt.Errorf("%w: amount cannot be zero or negative", ErrInvalidOrder)
	}
	if b.InterestRate.Sign() <= 0 {
		return fmt.Errorf("%w: interest rate cannot be zero or negative", ErrInvalidOrder)
	}
	if b.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidOrder)
	}
	return nil
}
