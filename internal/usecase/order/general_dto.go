package order

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindOrderOutputDTO struct {
	Id           uint                `json:"id"`
	CampaignId   uint                `json:"campaign_id"`
	BadgeChainId uint64              `json:"badge_chain_id"`
	Investor     custom_type.Address `json:"investor"`
	Amount       *uint256.Int        `json:"amount"`
	InterestRate *uint256.Int        `json:"interest_rate"`
	State        string              `json:"state"`
	CreatedAt    int64               `json:"created_at"`
	UpdatedAt    int64               `json:"updated_at"`
}
