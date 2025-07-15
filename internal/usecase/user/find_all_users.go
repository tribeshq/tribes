package user

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllUsersOutputDTO []*UserOutputDTO

type FindAllUsersUseCase struct {
	userRepository repository.UserRepository
}

func NewFindAllUsersUseCase(
	userRepo repository.UserRepository,
) *FindAllUsersUseCase {
	return &FindAllUsersUseCase{
		userRepository: userRepo,
	}
}

func (u *FindAllUsersUseCase) Execute() (*FindAllUsersOutputDTO, error) {
	res, err := u.userRepository.FindAllUsers()
	if err != nil {
		return nil, err
	}
	output := make(FindAllUsersOutputDTO, len(res))
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
	return &output, nil
}
