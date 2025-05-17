package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type UserAdvanceHandlers struct {
	UserRepository     repository.UserRepository
	ContractRepository repository.ContractRepository
}

func NewUserAdvanceHandlers(userRepository repository.UserRepository, contractRepository repository.ContractRepository) *UserAdvanceHandlers {
	return &UserAdvanceHandlers{
		UserRepository:     userRepository,
		ContractRepository: contractRepository,
	}
}

func (h *UserAdvanceHandlers) CreateUser(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.CreateUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	createUser := user_usecase.NewCreateUserUseCase(h.UserRepository)
	res, err := createUser.Execute(ctx, &input, metadata)
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

func (h *UserAdvanceHandlers) UpdateUser(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.UpdateUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	updateUser := user_usecase.NewUpdateUserUseCase(h.UserRepository)
	res, err := updateUser.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("user updated - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) DeleteUser(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.DeleteUserInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	deleteUserByAddress := user_usecase.NewDeleteUserUseCase(h.UserRepository)
	if err := deleteUserByAddress.Execute(ctx, &input); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	user, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	env.Notice(append([]byte("user deleted - "), user...))
	return nil
}

func (h *UserAdvanceHandlers) Withdraw(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input user_usecase.WithdrawInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findUserByAddress := user_usecase.NewCreateUserUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(ctx, &user_usecase.CreateUserInputDTO{
		Address: Address(metadata.MsgSender),
	}, metadata)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	tokenAddr := common.Address(input.Token)
	amount := input.Amount.ToBig()
	msgSender := metadata.MsgSender

	switch entity.UserRole(res.Role) {
	case entity.UserRoleAdmin:
		// The Admin role can withdraw the entire Application Balance if wanted
		if err := env.ERC20Transfer(
			tokenAddr,
			env.AppAddress(),
			msgSender,
			amount,
		); err != nil {
			return fmt.Errorf("failed to transfer ERC20: %w", err)
		}

		if _, err := env.ERC20Withdraw(
			tokenAddr,
			msgSender,
			amount,
		); err != nil {
			return fmt.Errorf("failed to withdraw ERC20: %w", err)
		}

	default:
		if _, err := env.ERC20Withdraw(
			tokenAddr,
			msgSender,
			amount,
		); err != nil {
			return fmt.Errorf("failed to withdraw ERC20: %w", err)
		}
	}

	return nil
}
