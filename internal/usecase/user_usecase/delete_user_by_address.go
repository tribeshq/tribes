package user_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type DeleteUserInputDTO struct {
	Address Address `json:"address"`
}

type DeleteUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewDeleteUserUseCase(userRepository repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *DeleteUserUseCase) Execute(ctx context.Context, input *DeleteUserInputDTO) error {
	return u.UserRepository.DeleteUser(ctx, input.Address)
}
