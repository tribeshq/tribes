package social_account

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindSocialAccountsByUserIdInputDTO struct {
	UserId uint `json:"user_id" validate:"required"`
}

type FindSocialAccountsByUserIdOutputDTO []*FindSocialAccountOutputDTO

type FindSocialAccountsByUserIdUseCase struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewFindSocialAccountsByUserIdUseCase(socialAccountRepository repository.SocialAccountRepository) *FindSocialAccountsByUserIdUseCase {
	return &FindSocialAccountsByUserIdUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *FindSocialAccountsByUserIdUseCase) Execute(ctx context.Context, input *FindSocialAccountsByUserIdInputDTO) (*FindSocialAccountsByUserIdOutputDTO, error) {
	socialAccounts, err := s.SocialAccountRepository.FindSocialAccountsByUserId(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	output := make(FindSocialAccountsByUserIdOutputDTO, len(socialAccounts))
	for i, socialAccount := range socialAccounts {
		output[i] = &FindSocialAccountOutputDTO{
			Id:        socialAccount.Id,
			UserId:    socialAccount.UserId,
			Username:  socialAccount.Username,
			Platform:  string(socialAccount.Platform),
			CreatedAt: socialAccount.CreatedAt,
			UpdatedAt: socialAccount.UpdatedAt,
		}
	}
	return &output, nil
}
