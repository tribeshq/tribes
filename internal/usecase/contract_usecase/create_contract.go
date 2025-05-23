package contract_usecase

import (
	"context"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateContractInputDTO struct {
	Symbol  string  `json:"symbol"`
	Address Address `json:"address"`
}

type CreateContractOutputDTO struct {
	Id        uint    `json:"id"`
	Symbol    string  `json:"symbol"`
	Address   Address `json:"address"`
	CreatedAt int64   `json:"created_at"`
}

type CreateContractUseCase struct {
	ContractRepository repository.ContractRepository
}

func NewCreateContractUseCase(contractRepository repository.ContractRepository) *CreateContractUseCase {
	return &CreateContractUseCase{
		ContractRepository: contractRepository,
	}
}

func (s *CreateContractUseCase) Execute(ctx context.Context, input *CreateContractInputDTO, metadata rollmelette.Metadata) (*CreateContractOutputDTO, error) {
	contract, err := entity.NewContract(input.Symbol, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}
	res, err := s.ContractRepository.CreateContract(ctx, contract)
	if err != nil {
		return nil, err
	}
	output := &CreateContractOutputDTO{
		Id:        res.Id,
		Symbol:    res.Symbol,
		Address:   res.Address,
		CreatedAt: res.CreatedAt,
	}
	return output, nil
}
