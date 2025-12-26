package advance

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/configs"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type EmergencyAdvanceHandlers struct {
	Config *configs.RollupConfig
}

func NewEmergencyAdvanceHandlers(cfg *configs.RollupConfig) *EmergencyAdvanceHandlers {
	return &EmergencyAdvanceHandlers{
		Config: cfg,
	}
}

func (h *EmergencyAdvanceHandlers) EmergencyERC20Withdraw(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user.EmergencyERC20WithdrawInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	abiJSON := `[{
		"type":"function",
		"name":"emergencyERC20Withdraw",
		"inputs":[
			{"type":"address"},
			{"type":"address"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	delegatecallPayload, err := abiInterface.Pack(
		"emergencyERC20Withdraw",
		input.Token,
		input.To,
	)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %w", err)
	}

	env.DelegateCallVoucher(
		h.Config.EmergencyWithdrawAddress,
		delegatecallPayload,
	)
	return nil
}

func (h *EmergencyAdvanceHandlers) EmergencyEtherWithdraw(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user.EmergencyEtherWithdrawInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	abiJSON := `[{
		"type":"function",
		"name":"emergencyETHWithdraw",
		"inputs":[
			{"type":"address"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	delegatecallPayload, err := abiInterface.Pack(
		"emergencyETHWithdraw",
		input.To,
	)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %w", err)
	}

	env.DelegateCallVoucher(
		h.Config.EmergencyWithdrawAddress,
		delegatecallPayload,
	)
	return nil
}
