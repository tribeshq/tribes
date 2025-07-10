package user

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type FindUserByAddressInputDTO struct {
	Address custom_type.Address `json:"address" validate:"required"`
}

type FindUserByAddressUseCase struct {
	UserRepository repository.UserRepository
}

func NewFindUserByAddressUseCase(userRepository repository.UserRepository) *FindUserByAddressUseCase {
	return &FindUserByAddressUseCase{
		UserRepository: userRepository,
	}
}

func (u *FindUserByAddressUseCase) Execute(ctx context.Context, input *FindUserByAddressInputDTO) (*UserOutputDTO, error) {
	res, err := u.UserRepository.FindUserByAddress(ctx, input.Address)
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
