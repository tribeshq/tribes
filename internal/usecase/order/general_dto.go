package order

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
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
