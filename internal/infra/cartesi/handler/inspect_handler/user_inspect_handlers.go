package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
)

type UserInspectHandlers struct {
	UserRepository     repository.UserRepository
	ContractRepository repository.ContractRepository
}

func NewUserInspectHandlers(userRepository repository.UserRepository, crowdfundingRepository repository.ContractRepository) *UserInspectHandlers {
	return &UserInspectHandlers{
		UserRepository:     userRepository,
		ContractRepository: crowdfundingRepository,
	}
}

func (h *UserInspectHandlers) FindUserByAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input user_usecase.FindUserByAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findUserByAddress := user_usecase.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find User: %w", err)
	}
	User, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal User: %w", err)
	}
	env.Report(User)
	return nil
}

func (h *UserInspectHandlers) FindAllUsers(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllUsers := user_usecase.NewFindAllUsersUseCase(h.UserRepository)
	res, err := findAllUsers.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all Users: %w", err)
	}
	allUsers, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all Users: %w", err)
	}
	env.Report(allUsers)
	return nil
}

func (h *UserInspectHandlers) Balance(env rollmelette.EnvInspector, payload []byte) error {
	var input user_usecase.FindUserByAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findAllContracts := contract_usecase.NewFindAllContractsUseCase(h.ContractRepository)
	contracts, err := findAllContracts.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all contracts: %w", err)
	}

	findUserbyAddress := user_usecase.NewFindUserByAddressUseCase(h.UserRepository)
	user, err := findUserbyAddress.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	switch user.Role {
	case string(entity.UserRoleAdmin):
		appAddress := env.AppAddress()
		balances := make(map[string]string)
		for _, contract := range contracts {
			balances[contract.Symbol] = env.ERC20BalanceOf(common.Address(contract.Address), appAddress).String()
		}
		balanceBytes, err := json.Marshal(balances)
		if err != nil {
			return fmt.Errorf("failed to marshal balances: %w", err)
		}
		env.Report(balanceBytes)
		return nil
	default:
		balances := make(map[string]string)
		for _, contract := range contracts {
			balances[contract.Symbol] = env.ERC20BalanceOf(common.Address(contract.Address), common.Address(user.Address)).String()
		}
		balanceBytes, err := json.Marshal(balances)
		if err != nil {
			return fmt.Errorf("failed to marshal balances: %w", err)
		}
		env.Report(balanceBytes)
		return nil
	}
}
