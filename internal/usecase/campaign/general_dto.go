package campaign

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindCampaignOutputDTO struct {
	Id                uint                `json:"id"`
	Token             custom_type.Address `json:"token"`
	Creator           custom_type.Address `json:"creator"`
	CollateralAddress custom_type.Address `json:"collateral_address"`
	CollateralAmount  *uint256.Int        `json:"collateral_amount"`
	BadgeRouter       custom_type.Address `json:"badge_router"`
	BadgeMinter       custom_type.Address `json:"badge_minter"`
	DebtIssued        *uint256.Int        `json:"debt_issued"`
	MaxInterestRate   *uint256.Int        `json:"max_interest_rate"`
	TotalObligation   *uint256.Int        `json:"total_obligation"`
	TotalRaised       *uint256.Int        `json:"total_raised"`
	State             string              `json:"state"`
	Orders            []*entity.Order     `json:"orders"`
	CreatedAt         int64               `json:"created_at"`
	ClosesAt          int64               `json:"closes_at"`
	MaturityAt        int64               `json:"maturity_at"`
	UpdatedAt         int64               `json:"updated_at"`
}
