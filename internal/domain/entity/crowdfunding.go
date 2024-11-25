package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrExpired              = errors.New("crowdfunding expired")
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

type CrowdfundingRepository interface {
	CreateCrowdfunding(ctx context.Context, crowdfunding *Crowdfunding) (*Crowdfunding, error)
	FindCrowdfundingsByCreator(ctx context.Context, creator common.Address) ([]*Crowdfunding, error)
	FindCrowdfundingsByInvestor(ctx context.Context, investor common.Address) ([]*Crowdfunding, error)
	FindCrowdfundingById(ctx context.Context, id uint) (*Crowdfunding, error)
	FindAllCrowdfundings(ctx context.Context) ([]*Crowdfunding, error)
	UpdateCrowdfunding(ctx context.Context, crowdfunding *Crowdfunding) (*Crowdfunding, error)
	DeleteCrowdfunding(ctx context.Context, id uint) error
}

type Crowdfunding struct {
	Id              uint              `json:"id" gorm:"primaryKey"`
	Creator         common.Address    `json:"creator,omitempty" gorm:"type:text;not null"`
	DebtIssued      *uint256.Int      `json:"debt_issued,omitempty" gorm:"type:text;not null"`
	MaxInterestRate *uint256.Int      `json:"max_interest_rate,omitempty" gorm:"type:text;not null"`
	TotalObligation *uint256.Int      `json:"total_obligation,omitempty" gorm:"type:text;not null;default:0"`
	State           CrowdfundingState `json:"state,omitempty" gorm:"type:text;not null"`
	Orders          []*Order          `json:"orders,omitempty" gorm:"foreignKey:CrowdfundingId;constraint:OnDelete:CASCADE"`
	ClosesAt        int64             `json:"closes_at,omitempty" gorm:"not null"`
	MaturityAt      int64             `json:"maturity_at,omitempty" gorm:"not null"`
	CreatedAt       int64             `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt       int64             `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewCrowdfunding(creator common.Address, debt_issued *uint256.Int, maxInterestRate *uint256.Int, closesAt int64, maturityAt int64, createdAt int64) (*Crowdfunding, error) {
	crowdfunding := &Crowdfunding{
		Creator:         creator,
		DebtIssued:      debt_issued,
		MaxInterestRate: maxInterestRate,
		State:           CrowdfundingStateUnderReview,
		ClosesAt:        closesAt,
		MaturityAt:      maturityAt,
		CreatedAt:       createdAt,
	}
	if err := crowdfunding.Validate(); err != nil {
		return nil, err
	}
	return crowdfunding, nil
}

func (a *Crowdfunding) Validate() error {
	if a.Creator == (common.Address{}) {
		return fmt.Errorf("%w: invalid creator address", ErrInvalidCrowdfunding)
	}
	if a.DebtIssued.Sign() == 0 {
		return fmt.Errorf("%w: debt issued cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.DebtIssued.Cmp(uint256.NewInt(15000000000)) > 0 {
		return fmt.Errorf("%w: debt issued exceeds the maximum allowed value", ErrInvalidCrowdfunding)
	}
	if a.MaxInterestRate.Sign() == 0 {
		return fmt.Errorf("%w: max interest rate cannot be zero", ErrInvalidCrowdfunding)
	}
	if a.ClosesAt == 0 {
		return fmt.Errorf("%w: expiration date is missing", ErrInvalidCrowdfunding)
	}
	if a.ClosesAt > a.CreatedAt+15552000 {
		return fmt.Errorf("%w: expiration date cannot be greater than 6 months", ErrInvalidCrowdfunding)
	}
	if a.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidCrowdfunding)
	}
	if a.CreatedAt >= a.ClosesAt {
		return fmt.Errorf("%w: creation date cannot be greater than or equal to expiration date", ErrInvalidCrowdfunding)
	}
	if a.MaturityAt == 0 {
		return fmt.Errorf("%w: maturity date is missing", ErrInvalidCrowdfunding)
	}
	return nil
}
