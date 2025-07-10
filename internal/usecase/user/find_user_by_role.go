package user

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindUserByRoleInputDTO struct {
	Role string `json:"role" validate:"required"`
}

type FindUserByRoleOutputDTO []*UserOutputDTO

type FindUserByRoleUseCase struct {
	userRepository repository.UserRepository
}

func NewFindUserByRoleUseCase(userRepository repository.UserRepository) *FindUserByRoleUseCase {
	return &FindUserByRoleUseCase{
		userRepository: userRepository,
	}
}

func (u *FindUserByRoleUseCase) Execute(ctx context.Context, input *FindUserByRoleInputDTO) ([]*UserOutputDTO, error) {
	res, err := u.userRepository.FindUsersByRole(ctx, input.Role)
	if err != nil {
		return nil, err
	}
	output := make(FindUserByRoleOutputDTO, len(res))
	for i, user := range res {
		output[i] = &UserOutputDTO{
			Id:             user.Id,
			Role:           string(user.Role),
			Address:        user.Address,
			SocialAccounts: user.SocialAccounts,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		}
	}
	return output, nil
}
