package social_account

import (
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CreateSocialAccountInputDTO struct {
	Address  custom_type.Address `json:"address" validate:"required"`
	Username string              `json:"username" validate:"required"`
	Platform string              `json:"platform" validate:"required"`
}

type CreateSocialAccountOutputDTO struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
}

type CreateSocialAccountUseCase struct {
	userRepository          repository.UserRepository
	socialAccountRepository repository.SocialAccountRepository
}

func NewCreateSocialAccountUseCase(
	userRepo repository.UserRepository,
	socialAccountRepo repository.SocialAccountRepository,
) *CreateSocialAccountUseCase {
	return &CreateSocialAccountUseCase{
		userRepository:          userRepo,
		socialAccountRepository: socialAccountRepo,
	}
}

func (s *CreateSocialAccountUseCase) Execute(input *CreateSocialAccountInputDTO, metadata *rollmelette.Metadata) (*CreateSocialAccountOutputDTO, error) {
	user, err := s.userRepository.FindUserByAddress(input.Address)
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

	socialAccount, err = s.socialAccountRepository.CreateSocialAccount(socialAccount)
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
