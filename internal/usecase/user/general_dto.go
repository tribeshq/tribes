package user

import (
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type BalanceOfInputDTO struct {
	Token   custom_type.Address `json:"token"`
	Address custom_type.Address `json:"address" validate:"required"`
}

type UserOutputDTO struct {
	Id             uint                    `json:"id"`
	Role           string                  `json:"role"`
	Address        custom_type.Address     `json:"address"`
	SocialAccounts []*entity.SocialAccount `json:"social_accounts"`
	CreatedAt      int64                   `json:"created_at"`
	UpdatedAt      int64                   `json:"updated_at"`
}
