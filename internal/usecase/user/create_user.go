package user

import (
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateUserInputDTO struct {
	Role    string              `json:"role" validate:"required"`
	Address custom_type.Address `json:"address" validate:"required"`
}

type CreateUserOutputDTO struct {
	Id              uint                    `json:"id"`
	Role            string                  `json:"role"`
	Address         custom_type.Address     `json:"address"`
	SocialAccounts  []*entity.SocialAccount `json:"social_accounts"`
	InvestmentLimit *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	CreatedAt       int64                   `json:"created_at"`
}

type CreateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewCreateUserUseCase(
	userRepo repository.UserRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepo,
	}
}

func (u *CreateUserUseCase) Execute(input *CreateUserInputDTO, metadata rollmelette.Metadata) (*CreateUserOutputDTO, error) {
	user, err := entity.NewUser(input.Role, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}

	res, err := u.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutputDTO{
		Id:             res.Id,
		Role:           string(res.Role),
		Address:        res.Address,
		SocialAccounts: res.SocialAccounts,
		CreatedAt:      res.CreatedAt,
	}, nil
}
