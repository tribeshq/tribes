package user

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindUserByRoleInputDTO struct {
	Role string `json:"role" validate:"required"`
}

type FindUserByRoleOutputDTO []*UserOutputDTO

type FindUsersByRoleUseCase struct {
	userRepository repository.UserRepository
}

func NewFindUsersByRoleUseCase(
	userRepo repository.UserRepository,
) *FindUsersByRoleUseCase {
	return &FindUsersByRoleUseCase{
		userRepository: userRepo,
	}
}

func (u *FindUsersByRoleUseCase) Execute(input *FindUserByRoleInputDTO) ([]*UserOutputDTO, error) {
	res, err := u.userRepository.FindUsersByRole(input.Role)
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
