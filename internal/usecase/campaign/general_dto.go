package campaign

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
)

type CampaignOutputDTO struct {
	Id                uint                    `json:"id"`
	Title             string                  `json:"title,omitempty"`
	Description       string                  `json:"description,omitempty"`
	Promotion         string                  `json:"promotion,omitempty"`
	Token             types.Address           `json:"token"`
	Creator           *user.UserOutputDTO     `json:"creator"`
	CollateralAddress types.Address           `json:"collateral"`
	CollateralAmount  *uint256.Int            `json:"collateral_amount"`
	BadgeAddress      types.Address           `json:"badge_address"`
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
