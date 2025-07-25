package advance

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/configs"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type UserAdvanceHandlers struct {
	cfg            *configs.RollupConfig
	userRepository repository.UserRepository
}

func NewUserAdvanceHandlers(
	cfg *configs.RollupConfig,
	userRepo repository.UserRepository,
) *UserAdvanceHandlers {
	return &UserAdvanceHandlers{
		cfg:            cfg,
		userRepository: userRepo,
	}
}

func (h *UserAdvanceHandlers) CreateUser(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user.CreateUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	createUser := user.NewCreateUserUseCase(h.userRepository)
	res, err := createUser.Execute(&input, metadata)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("user created - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) DeleteUser(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user.DeleteUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	deleteUserByAddress := user.NewDeleteUserUseCase(h.userRepository)
	if err := deleteUserByAddress.Execute(&input); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	user, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	env.Notice(append([]byte("user deleted - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) ERC20Withdraw(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user.WithdrawInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	findUserByAddress := user.NewFindUserByAddressUseCase(h.userRepository)
	res, err := findUserByAddress.Execute(&user.FindUserByAddressInputDTO{
		Address: custom_type.Address(metadata.MsgSender),
	})
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// For admin, transfer from app address to admin first, then withdraw
	if entity.UserRole(res.Role) == entity.UserRoleAdmin {
		if err := env.ERC20Transfer(
			common.Address(input.Token),
			env.AppAddress(),
			metadata.MsgSender,
			input.Amount.ToBig(),
		); err != nil {
			return fmt.Errorf("failed to transfer ERC20 from app to admin: %w", err)
		}
	}

	// Withdraw tokens to the user's address
	if _, err := env.ERC20Withdraw(
		common.Address(input.Token),
		metadata.MsgSender,
		input.Amount.ToBig(),
	); err != nil {
		return fmt.Errorf("failed to withdraw ERC20: %w", err)
	}

	env.Notice([]byte(
		fmt.Sprintf(
			"ERC20 withdrawn - token: %s, amount: %s, user: %s", common.Address(input.Token), input.Amount.ToBig(), metadata.MsgSender,
		),
	))
	return nil
}

func (h *UserAdvanceHandlers) EmergencyERC20Withdraw(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
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
		input.To,
		metadata.MsgSender,
		input.Token,
	)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %w", err)
	}

	env.DelegateCallVoucher(
		h.cfg.EmergencyWithdrawAddress,
		delegatecallPayload,
	)
	return nil
}

func (h *UserAdvanceHandlers) EmergencyEtherWithdraw(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
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
			{"type":"address"},
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
		metadata.MsgSender,
	)
	if err != nil {
		return fmt.Errorf("failed to pack ABI: %w", err)
	}

	env.DelegateCallVoucher(
		h.cfg.EmergencyWithdrawAddress,
		delegatecallPayload,
	)
	return nil
}
