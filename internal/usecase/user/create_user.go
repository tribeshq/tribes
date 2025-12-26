package user

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	. "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CreateUserInputDTO struct {
	Role    string  `json:"role" validate:"required"`
	Address Address `json:"address" validate:"required"`
}

type CreateUserOutputDTO struct {
	Id              uint                    `json:"id"`
	Role            string                  `json:"role"`
	Address         Address                 `json:"address"`
	SocialAccounts  []*entity.SocialAccount `json:"social_accounts"`
	InvestmentLimit *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	CreatedAt       int64                   `json:"created_at"`
}

type CreateUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewCreateUserUseCase(
	userRepo repository.UserRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepo,
	}
}

func (u *CreateUserUseCase) Execute(input *CreateUserInputDTO, metadata rollmelette.Metadata) (*CreateUserOutputDTO, error) {
	user, err := entity.NewUser(input.Role, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}

	res, err := u.UserRepository.CreateUser(user)
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
