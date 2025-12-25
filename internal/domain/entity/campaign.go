package entity

import (
	"errors"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
)

var (
	ErrCampaignNotFound = errors.New("campaign not found")
	ErrInvalidCampaign  = errors.New("invalid campaign")
)

type CampaignState string

const (
	CampaignStateClosed             CampaignState = "closed"
	CampaignStateOngoing            CampaignState = "ongoing"
	CampaignStateCanceled           CampaignState = "canceled"
	CampaignStateSettled            CampaignState = "settled"
	CampaignStateCollateralExecuted CampaignState = "collateral_executed"
)

type Campaign struct {
	Id                uint          `json:"id" gorm:"primaryKey"`
	Title             string        `json:"title,omitempty" gorm:"not null"`
	Description       string        `json:"description,omitempty" gorm:"not null"`
	Promotion         string        `json:"promotion,omitempty" gorm:"not null"`
	Token             types.Address `json:"token,omitempty" gorm:"types:text;not null"`
	CreatorAddress    types.Address `json:"creator_address,omitempty" gorm:"types:text;not null"`
	CollateralAddress types.Address `json:"collateral_address,omitempty" gorm:"types:text;not null"`
	CollateralAmount  *uint256.Int  `json:"collateral_amount,omitempty" gorm:"types:text;not null"`
	BadgeAddress      types.Address `json:"badge_address,omitempty" gorm:"types:text;not null"`
	DebtIssued        *uint256.Int  `json:"debt_issued,omitempty" gorm:"types:text;not null"`
	MaxInterestRate   *uint256.Int  `json:"max_interest_rate,omitempty" gorm:"types:text;not null"`
	TotalObligation   *uint256.Int  `json:"total_obligation,omitempty" gorm:"types:text;not null;default:0"`
	TotalRaised       *uint256.Int  `json:"total_raised,omitempty" gorm:"types:text;not null;default:0"`
	State             CampaignState `json:"state,omitempty" gorm:"types:text;not null"`
	Orders            []*Order      `json:"orders,omitempty" gorm:"foreignKey:CampaignId;constraint:OnDelete:CASCADE"`
	ClosesAt          int64         `json:"closes_at,omitempty" gorm:"not null"`
	MaturityAt        int64         `json:"maturity_at,omitempty" gorm:"not null"`
	CreatedAt         int64         `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt         int64         `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewCampaign(title string, description string, promotion string, token types.Address, creatorAddress types.Address, collateralAddress types.Address, collateralAmount *uint256.Int, badgeAddress types.Address, debtIssued *uint256.Int, maxInterestRate *uint256.Int, closesAt int64, maturityAt int64, createdAt int64) (*Campaign, error) {
	campaign := &Campaign{
		Title:             title,
		Description:       description,
		Promotion:         promotion,
		Token:             token,
		CreatorAddress:    creatorAddress,
		CollateralAddress: collateralAddress,
		CollateralAmount:  collateralAmount,
		BadgeAddress:      badgeAddress,
		DebtIssued:        debtIssued,
		MaxInterestRate:   maxInterestRate,
		State:             CampaignStateOngoing,
		Orders:            []*Order{},
		ClosesAt:          closesAt,
		MaturityAt:        maturityAt,
		CreatedAt:         createdAt,
	}
	if err := campaign.validate(); err != nil {
		return nil, err
	}
	return campaign, nil
}

func (a *Campaign) validate() error {
	if a.Title == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrInvalidCampaign)
	}
	if a.Description == "" {
		return fmt.Errorf("%w: description cannot be empty", ErrInvalidCampaign)
	}
	if a.Promotion == "" {
		return fmt.Errorf("%w: promotion cannot be empty", ErrInvalidCampaign)
	}
	if a.Token == (types.Address{}) {
		return fmt.Errorf("%w: invalid token address", ErrInvalidCampaign)
	}
	if a.CreatorAddress == (types.Address{}) {
		return fmt.Errorf("%w: invalid creator address", ErrInvalidCampaign)
	}
	if a.CollateralAddress == (types.Address{}) {
		return fmt.Errorf("%w: invalid collateral address", ErrInvalidCampaign)
	}
	if a.CollateralAmount.Sign() == 0 {
		return fmt.Errorf("%w: collateral amount cannot be zero", ErrInvalidCampaign)
	}
	if a.BadgeAddress == (types.Address{}) {
		return fmt.Errorf("%w: invalid badge address", ErrInvalidCampaign)
	}
	if a.DebtIssued.Sign() == 0 {
		return fmt.Errorf("%w: debt issued cannot be zero", ErrInvalidCampaign)
	}
	if a.MaxInterestRate.Sign() == 0 {
		return fmt.Errorf("%w: max interest rate cannot be zero", ErrInvalidCampaign)
	}
	if a.CreatedAt == 0 {
		return fmt.Errorf("%w: creation date is missing", ErrInvalidCampaign)
	}
	if a.ClosesAt == 0 {
		return fmt.Errorf("%w: close date is missing", ErrInvalidCampaign)
	}
	if a.MaturityAt == 0 {
		return fmt.Errorf("%w: maturity date is missing", ErrInvalidCampaign)
	}
	return nil
}
