package user

import (
	"github.com/holiman/uint256"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type WithdrawInputDTO struct {
	Token  custom_type.Address `json:"token" validate:"required"`
	Amount *uint256.Int        `json:"amount" validate:"required"`
}

type EmergencyERC20WithdrawInputDTO struct {
	To                       custom_type.Address `json:"to" validate:"required"`
	Token                    custom_type.Address `json:"token" validate:"required"`
	EmergencyWithdrawAddress custom_type.Address `json:"emergency_withdraw_address" validate:"required"`
}

type EmergencyEtherWithdrawInputDTO struct {
	To                       custom_type.Address `json:"to" validate:"required"`
	EmergencyWithdrawAddress custom_type.Address `json:"emergency_withdraw_address" validate:"required"`
}
