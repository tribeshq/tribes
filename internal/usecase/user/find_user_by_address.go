package user

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
)

type FindUserByAddressInputDTO struct {
	Address types.Address `json:"address" validate:"required"`
}

type FindUserByAddressUseCase struct {
	UserRepository repository.UserRepository
}

func NewFindUserByAddressUseCase(
	userRepo repository.UserRepository,
) *FindUserByAddressUseCase {
	return &FindUserByAddressUseCase{
		UserRepository: userRepo,
	}
}

func (u *FindUserByAddressUseCase) Execute(input *FindUserByAddressInputDTO) (*UserOutputDTO, error) {
	res, err := u.UserRepository.FindUserByAddress(input.Address)
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
