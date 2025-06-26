package user

import (
	"github.com/holiman/uint256"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type WithdrawInputDTO struct {
	Token  Address      `json:"token" validate:"required"`
	Amount *uint256.Int `json:"amount" validate:"required"`
}

type EmergencyERC20WithdrawInputDTO struct {
	To                       Address `json:"to" validate:"required"`
	Token                    Address `json:"token" validate:"required"`
	EmergencyWithdrawAddress Address `json:"emergency_withdraw_address" validate:"required"`
}

type EmergencyEtherWithdrawInputDTO struct {
	To                       Address `json:"to" validate:"required"`
	EmergencyWithdrawAddress Address `json:"emergency_withdraw_address" validate:"required"`
}
