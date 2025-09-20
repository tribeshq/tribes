package user

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type DeleteUserInputDTO struct {
	Address custom_type.Address `json:"address" validate:"required"`
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
