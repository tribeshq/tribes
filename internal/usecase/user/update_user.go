package user

import (
	"context"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	. "github.com/tribeshq/tribes/pkg/custom_type"
)

type UpdateUserInputDTO struct {
	Role            string       `json:"role" validate:"required"`
	Address         Address      `json:"address" validate:"required"`
	InvestmentLimit *uint256.Int `json:"investment_limit,omitempty" gorm:"type:bigint" validate:"required"`
}

type UpdateUserOutputDTO struct {
	Id              uint                    `json:"id"`
	Role            string                  `json:"role"`
	Address         Address                 `json:"address"`
	SocialAccounts  []*entity.SocialAccount `json:"social_accounts"`
	InvestmentLimit *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	CreatedAt       int64                   `json:"created_at"`
	UpdatedAt       int64                   `json:"updated_at"`
}

type UpdateUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewUpdateUserUseCase(userRepository repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *UpdateUserUseCase) Execute(ctx context.Context, input *UpdateUserInputDTO, metadata rollmelette.Metadata) (*UpdateUserOutputDTO, error) {
	user, err := u.UserRepository.UpdateUser(ctx, &entity.User{
		Role:            entity.UserRole(input.Role),
		Address:         input.Address,
		UpdatedAt:       metadata.BlockTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateUserOutputDTO{
		Id:              user.Id,
		Role:            string(user.Role),
		Address:         user.Address,
		SocialAccounts:  user.SocialAccounts,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}
