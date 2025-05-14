package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
)

type ContractInspectHandlers struct {
	ContractRepository repository.ContractRepository
}

func NewContractInspectHandlers(contractRepository repository.ContractRepository) *ContractInspectHandlers {
	return &ContractInspectHandlers{
		ContractRepository: contractRepository,
	}
}

func (h *ContractInspectHandlers) FindAllContracts(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllContracts := contract_usecase.NewFindAllContractsUseCase(h.ContractRepository)
	contracts, err := findAllContracts.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all contracts: %w", err)
	}
	contractsBytes, err := json.Marshal(contracts)
	if err != nil {
		return fmt.Errorf("failed to marshal contracts: %w", err)
	}
	env.Report(contractsBytes)
	return nil
}

func (h *ContractInspectHandlers) FindContractBySymbol(env rollmelette.EnvInspector, payload []byte) error {
	var input contract_usecase.FindContractBySymbolInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrderBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	contract, err := findOrderBySymbol.Execute(ctx, &input)
	if err != nil {
		return err
	}
	contractBytes, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}
	env.Report(contractBytes)
	return nil
}

func (h *ContractInspectHandlers) FindContractByAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input contract_usecase.FindContractByAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findContractByAddress := contract_usecase.NewFindContractByAddressUseCase(h.ContractRepository)
	contract, err := findContractByAddress.Execute(ctx, &input)
	if err != nil {
		return err
	}
	contractBytes, err := json.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}
	env.Report(contractBytes)
	return nil
}
