package user

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindUserByAddressInputDTO struct {
	Address custom_type.Address `json:"address" validate:"required"`
}

type FindUserByAddressUseCase struct {
	userRepository repository.UserRepository
}

func NewFindUserByAddressUseCase(
	userRepo repository.UserRepository,
) *FindUserByAddressUseCase {
	return &FindUserByAddressUseCase{
		userRepository: userRepo,
	}
}

func (u *FindUserByAddressUseCase) Execute(input *FindUserByAddressInputDTO) (*UserOutputDTO, error) {
	res, err := u.userRepository.FindUserByAddress(input.Address)
	if err != nil {
		return nil, err
	}
	return &UserOutputDTO{
		Id:             res.Id,
		Role:           string(res.Role),
		Address:        res.Address,
		SocialAccounts: res.SocialAccounts,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}, nil
}
