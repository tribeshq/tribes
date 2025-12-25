package user

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type FindAllUsersOutputDTO []*UserOutputDTO

type FindAllUsersUseCase struct {
	UserRepository repository.UserRepository
}

func NewFindAllUsersUseCase(
	userRepo repository.UserRepository,
) *FindAllUsersUseCase {
	return &FindAllUsersUseCase{
		UserRepository: userRepo,
	}
}

func (u *FindAllUsersUseCase) Execute() (*FindAllUsersOutputDTO, error) {
	res, err := u.UserRepository.FindAllUsers()
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
