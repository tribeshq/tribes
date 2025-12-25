package social_account

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type DeleteSocialAccountInputDTO struct {
	SocialAccountId uint `json:"social_account_id" validate:"required"`
}

type DeleteSocialAccountUseCase struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewDeleteSocialAccountUseCase(
	socialAccountRepo repository.SocialAccountRepository,
) *DeleteSocialAccountUseCase {
	return &DeleteSocialAccountUseCase{
		SocialAccountRepository: socialAccountRepo,
	}
}

func (s *DeleteSocialAccountUseCase) Execute(input *DeleteSocialAccountInputDTO) error {
	return s.SocialAccountRepository.DeleteSocialAccount(input.SocialAccountId)
}
