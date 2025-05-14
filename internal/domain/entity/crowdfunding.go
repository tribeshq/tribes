package entity

import (
	"errors"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

var (
	ErrCrowdfundingNotFound = errors.New("crowdfunding not found")
	ErrInvalidCrowdfunding  = errors.New("invalid crowdfunding")
)

type CrowdfundingState string

const (
	CrowdfundingStateUnderReview CrowdfundingState = "under_review"
	CrowdfundingStateClosed      CrowdfundingState = "closed"
	CrowdfundingStateOngoing     CrowdfundingState = "ongoing"
	CrowdfundingStateCanceled    CrowdfundingState = "canceled"
	CrowdfundingStateSettled     CrowdfundingState = "settled"
)

type Crowdfunding struct {
	Id                  uint                `json:"id" gorm:"primaryKey"`
	Token               custom_type.Address `json:"token,omitempty" gorm:"custom_type:text;not null"`
	Collateral          *uint256.Int        `json:"collateral,omitempty" gorm:"custom_type:text;not null"`
	Creator             custom_type.Address `json:"creator,omitempty" gorm:"custom_type:text;not null"`
	DebtIssued          *uint256.Int        `json:"debt_issued,omitempty" gorm:"custom_type:text;not null"`
	MaxInterestRate     *uint256.Int        `json:"max_interest_rate,omitempty" gorm:"custom_type:text;not null"`
	TotalObligation     *uint256.Int        `json:"total_obligation,omitempty" gorm:"custom_type:text;not null;default:0"`
	State               CrowdfundingState   `json:"state,omitempty" gorm:"custom_type:text;not null"`
	Orders              []*Order            `json:"orders,omitempty" gorm:"foreignKey:CrowdfundingId;constraint:OnDelete:CASCADE"`
	FundraisingDuration int64               `json:"fundraising_duration,omitempty" gorm:"not null"`
	ClosesAt            int64               `json:"closes_at,omitempty" gorm:"not null"`
	MaturityAt          int64               `json:"maturity_at,omitempty" gorm:"not null"`
	CreatedAt           int64               `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt           int64               `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewCrowdfunding(token custom_type.Address, amount *uint256.Int, creator custom_type.Address, debt_issued *uint256.Int, maxInterestRate *uint256.Int, fundraisingDuration int64, closesAt int64, maturityAt int64, createdAt int64) (*Crowdfunding, error) {
	crowdfunding := &Crowdfunding{
		Token:               token,
		Collateral:          amount,
		Creator:             creator,
		DebtIssued:          debt_issued,
		MaxInterestRate:     maxInterestRate,
		State:               CrowdfundingStateUnderReview,
		FundraisingDuration: fundraisingDuration,
		ClosesAt:            closesAt,
		MaturityAt:          maturityAt,
		CreatedAt:           createdAt,
	}
	if err := crowdfunding.validate(); err != nil {
		return nil, err
	}
	return crowdfunding, nil
}

func (a *Crowdfunding) validate() error {
	if a.Token == (custom_type.Address{}) {
		return fmt.Errorf("%w: invalid token address", ErrInvalidCrowdfunding)
	}
	if a.Collateral.Sign() == 0 {
		return fmt.Errorf("%w: collateral cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.Creator == (custom_type.Address{}) {
		return fmt.Errorf("%w: invalid creator address", ErrInvalidCrowdfunding)
	}
	if a.DebtIssued.Sign() == 0 {
		return fmt.Errorf("%w: debt issued cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.MaxInterestRate.Sign() == 0 {
		return fmt.Errorf("%w: max interest rate cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidCrowdfunding)
	}
	if a.ClosesAt == 0 {
		return fmt.Errorf("%w: close date is missing", ErrInvalidCrowdfunding)
	}
	if a.MaturityAt == 0 {
		return fmt.Errorf("%w: maturity date is missing", ErrInvalidCrowdfunding)
	}
	return nil
}
