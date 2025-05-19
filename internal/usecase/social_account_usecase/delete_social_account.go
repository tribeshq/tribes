package social_account_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type DeleteSocialAccountInputDTO struct {
	SocialAccountId uint `json:"social_account_id"`
}

type DeleteSocialAccountUseCase struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewDeleteSocialAccountUseCase(socialAccountRepository repository.SocialAccountRepository) *DeleteSocialAccountUseCase {
	return &DeleteSocialAccountUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *DeleteSocialAccountUseCase) Execute(ctx context.Context, input *DeleteSocialAccountInputDTO) error {
	return s.SocialAccountRepository.DeleteSocialAccount(ctx, input.SocialAccountId)
}
