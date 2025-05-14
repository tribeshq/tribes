package contract_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindContractByAddressInputDTO struct {
	Address custom_type.Address `json:"address"`
}

type FindContractByAddressUseCase struct {
	ContractRepository repository.ContractRepository
}

func NewFindContractByAddressUseCase(contractRepository repository.ContractRepository) *FindContractByAddressUseCase {
	return &FindContractByAddressUseCase{
		ContractRepository: contractRepository,
	}
}

func (s *FindContractByAddressUseCase) Execute(ctx context.Context, input *FindContractByAddressInputDTO) (*FindContractOutputDTO, error) {
	res, err := s.ContractRepository.FindContractByAddress(ctx, input.Address)
	if err != nil {
		return nil, err
	}
	return &FindContractOutputDTO{
		Id:        res.Id,
		Symbol:    res.Symbol,
		Address:   res.Address,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
