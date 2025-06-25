package user

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindAllUsersOutputDTO []*FindUserOutputDTO

type FindAllUsersUseCase struct {
	UserRepository repository.UserRepository
}

func NewFindAllUsersUseCase(userRepository repository.UserRepository) *FindAllUsersUseCase {
	return &FindAllUsersUseCase{
		UserRepository: userRepository,
	}
}

func (u *FindAllUsersUseCase) Execute(ctx context.Context) (*FindAllUsersOutputDTO, error) {
	res, err := u.UserRepository.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	output := make(FindAllUsersOutputDTO, len(res))
	for i, user := range res {
		output[i] = &FindUserOutputDTO{
			Id:              user.Id,
			Role:            string(user.Role),
			Address:         user.Address,
			SocialAccounts:  user.SocialAccounts,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		}
	}
	return &output, nil
}
