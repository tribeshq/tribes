package user

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type DeleteUserInputDTO struct {
	Address types.Address `json:"address" validate:"required"`
}

type DeleteUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewDeleteUserUseCase(
	userRepo repository.UserRepository,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: userRepo,
	}
}

func (u *DeleteUserUseCase) Execute(input *DeleteUserInputDTO) error {
	return u.UserRepository.DeleteUser(input.Address)
}
