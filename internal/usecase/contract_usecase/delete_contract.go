package contract_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type DeleteContractInputDTO struct {
	Symbol string
}

type DeleteContractUseCase struct {
	ContractReposiotry repository.ContractRepository
}

func NewDeleteContractUseCase(contractRepository repository.ContractRepository) *DeleteContractUseCase {
	return &DeleteContractUseCase{
		ContractReposiotry: contractRepository,
	}
}

func (s *DeleteContractUseCase) Execute(ctx context.Context, input *DeleteContractInputDTO) error {
	return s.ContractReposiotry.DeleteContract(ctx, input.Symbol)
}
