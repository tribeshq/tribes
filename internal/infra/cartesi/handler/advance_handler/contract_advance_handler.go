package advance_handler

import (
	"context"
	"encoding/json"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
)

type ContractAdvanceHandlers struct {
	ContractRepository entity.ContractRepository
}

func NewContractAdvanceHandlers(contractRepository entity.ContractRepository) *ContractAdvanceHandlers {
	return &ContractAdvanceHandlers{
		ContractRepository: contractRepository,
	}
}

func (h *ContractAdvanceHandlers) CreateContractHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.CreateContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	createContract := contract_usecase.NewCreateContractUseCase(h.ContractRepository)
	res, err := createContract.Execute(ctx, &input, metadata)
	if err != nil {
		return err
	}
	contract, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("contract created - "), contract...))
	return nil
}

func (h *ContractAdvanceHandlers) UpdateContractHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.UpdateContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	updateContract := contract_usecase.NewUpdateContractUseCase(h.ContractRepository)
	res, err := updateContract.Execute(ctx, &input, metadata)
	if err != nil {
		return err
	}
	contract, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("updated contract - "), contract...))
	return nil
}

func (h *ContractAdvanceHandlers) DeleteContractHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input contract_usecase.DeleteContractInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	deleteContract := contract_usecase.NewDeleteContractUseCase(h.ContractRepository)
	err := deleteContract.Execute(ctx, &input)
	if err != nil {
		return err
	}
	contract, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("deleted contract with - "), contract...))
	return nil
}
