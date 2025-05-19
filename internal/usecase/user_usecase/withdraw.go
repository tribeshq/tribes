package user_usecase

import (
	"github.com/holiman/uint256"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type WithdrawInputDTO struct {
	Token  Address      `json:"token"`
	Amount *uint256.Int `json:"amount"`
}
