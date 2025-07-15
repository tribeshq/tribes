package social_account

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type DeleteSocialAccountInputDTO struct {
	SocialAccountId uint `json:"social_account_id" validate:"required"`
}

type DeleteSocialAccountUseCase struct {
	socialAccountRepository repository.SocialAccountRepository
}

func NewDeleteSocialAccountUseCase(
	socialAccountRepo repository.SocialAccountRepository,
) *DeleteSocialAccountUseCase {
	return &DeleteSocialAccountUseCase{
		socialAccountRepository: socialAccountRepo,
	}
}

func (s *DeleteSocialAccountUseCase) Execute(input *DeleteSocialAccountInputDTO) error {
	return s.socialAccountRepository.DeleteSocialAccount(input.SocialAccountId)
}
