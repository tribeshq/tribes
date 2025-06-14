package user

import (
	"github.com/holiman/uint256"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type WithdrawInputDTO struct {
	Token  Address      `json:"token" validate:"required"`
	Amount *uint256.Int `json:"amount" validate:"required"`
}
