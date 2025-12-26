package advance

import (
	"encoding/json"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/configs"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	types "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type UserAdvanceHandlers struct {
	Config         *configs.RollupConfig
	UserRepository repository.UserRepository
}

func NewUserAdvanceHandlers(
	cfg *configs.RollupConfig,
	userRepo repository.UserRepository,
) *UserAdvanceHandlers {
	return &UserAdvanceHandlers{
		Config:         cfg,
		UserRepository: userRepo,
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

	createUser := user.NewCreateUserUseCase(h.UserRepository)
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

	deleteUserByAddress := user.NewDeleteUserUseCase(h.UserRepository)
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

	findUserByAddress := user.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(&user.FindUserByAddressInputDTO{
		Address: types.Address(metadata.MsgSender),
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
