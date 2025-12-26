package user

import (
	. "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
)

type WithdrawInputDTO struct {
	Token  Address      `json:"token" validate:"required"`
	Amount *uint256.Int `json:"amount" validate:"required"`
}

type EmergencyERC20WithdrawInputDTO struct {
	To    Address `json:"to" validate:"required"`
	Token Address `json:"token" validate:"required"`
}

type EmergencyEtherWithdrawInputDTO struct {
	To Address `json:"to" validate:"required"`
}
