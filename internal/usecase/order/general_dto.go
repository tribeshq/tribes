package order

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/holiman/uint256"
)

type OrderOutputDTO struct {
	Id           uint         `json:"id"`
	CampaignId   uint         `json:"campaign_id"`
	Investor     *entity.User `json:"investor"`
	Amount       *uint256.Int `json:"amount"`
	InterestRate *uint256.Int `json:"interest_rate"`
	State        string       `json:"state"`
	CreatedAt    int64        `json:"created_at"`
	UpdatedAt    int64        `json:"updated_at"`
}
