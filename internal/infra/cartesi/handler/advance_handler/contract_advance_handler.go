package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
)

type ContractAdvanceHandlers struct {
	ContractRepository repository.ContractRepository
}

func NewContractAdvanceHandlers(contractRepository repository.ContractRepository) *ContractAdvanceHandlers {
	return &ContractAdvanceHandlers{
		ContractRepository: contractRepository,
	}
}

func (h *ContractAdvanceHandlers) CreateContract(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.CreateContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	createContract := contract_usecase.NewCreateContractUseCase(h.ContractRepository)
	res, err := createContract.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	contract, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("contract created - "), contract...))
	return nil
}

func (h *ContractAdvanceHandlers) UpdateContract(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.UpdateContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	updateContract := contract_usecase.NewUpdateContractUseCase(h.ContractRepository)
	res, err := updateContract.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	contract, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("updated contract - "), contract...))
	return nil
}

func (h *ContractAdvanceHandlers) DeleteContract(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.DeleteContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	deleteContract := contract_usecase.NewDeleteContractUseCase(h.ContractRepository)
	if err := deleteContract.Execute(ctx, &input); err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	contract, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	env.Notice(append([]byte("deleted contract with - "), contract...))
	return nil
}
