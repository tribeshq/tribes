package user

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type DeleteUserInputDTO struct {
	Address custom_type.Address `json:"address" validate:"required"`
}

type DeleteUserUseCase struct {
	userRepository repository.UserRepository
}

func NewDeleteUserUseCase(
	userRepo repository.UserRepository,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepo,
	}
}

func (u *DeleteUserUseCase) Execute(input *DeleteUserInputDTO) error {
	return u.userRepository.DeleteUser(input.Address)
}
