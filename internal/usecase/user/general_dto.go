package user

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/internal/domain/entity"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type BalanceOfInputDTO struct {
	Token   Address `json:"token"`
	Address Address `json:"address" validate:"required"`
}

type FindUserOutputDTO struct {
	Id              uint                    `json:"id"`
	Role            string                  `json:"role"`
	Address         Address                 `json:"address"`
	SocialAccounts  []*entity.SocialAccount `json:"social_accounts"`
	InvestmentLimit *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	CreatedAt       int64                   `json:"created_at"`
	UpdatedAt       int64                   `json:"updated_at"`
}
