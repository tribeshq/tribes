package social_account

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	types "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/rollmelette/rollmelette"
)

type CreateSocialAccountInputDTO struct {
	Address  types.Address `json:"address" validate:"required"`
	Username string        `json:"username" validate:"required"`
	Platform string        `json:"platform" validate:"required"`
}

type CreateSocialAccountOutputDTO struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
}

type CreateSocialAccountUseCase struct {
	UserRepository          repository.UserRepository
	SocialAccountRepository repository.SocialAccountRepository
}

func NewCreateSocialAccountUseCase(
	userRepo repository.UserRepository,
	socialAccountRepo repository.SocialAccountRepository,
) *CreateSocialAccountUseCase {
	return &CreateSocialAccountUseCase{
		UserRepository:          userRepo,
		SocialAccountRepository: socialAccountRepo,
	}
}

func (s *CreateSocialAccountUseCase) Execute(input *CreateSocialAccountInputDTO, metadata *rollmelette.Metadata) (*CreateSocialAccountOutputDTO, error) {
	user, err := s.UserRepository.FindUserByAddress(input.Address)
	if err != nil {
		return nil, err
	}

	socialAccount, err := entity.NewSocialAccount(
		user.Id,
		input.Username,
		input.Platform,
		metadata.BlockTimestamp,
	)
	if err != nil {
		return nil, err
	}

	socialAccount, err = s.SocialAccountRepository.CreateSocialAccount(socialAccount)
	if err != nil {
		return nil, err
	}
	return &CreateSocialAccountOutputDTO{
		Id:        socialAccount.Id,
		UserId:    socialAccount.UserId,
		Username:  socialAccount.Username,
		Platform:  string(socialAccount.Platform),
		CreatedAt: socialAccount.CreatedAt,
	}, nil
}
