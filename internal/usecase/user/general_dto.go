package user

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type BalanceOfInputDTO struct {
	Token   types.Address `json:"token"`
	Address types.Address `json:"address" validate:"required"`
}

type UserOutputDTO struct {
	Id             uint                    `json:"id"`
	Role           string                  `json:"role"`
	Address        types.Address           `json:"address"`
	SocialAccounts []*entity.SocialAccount `json:"social_accounts"`
	CreatedAt      int64                   `json:"created_at"`
	UpdatedAt      int64                   `json:"updated_at"`
}
