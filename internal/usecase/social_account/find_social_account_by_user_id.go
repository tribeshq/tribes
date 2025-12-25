package social_account

import (
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
)

type FindSocialAccountsByUserIdInputDTO struct {
	UserId uint `json:"user_id" validate:"required"`
}

type FindSocialAccountsByUserIdOutputDTO []*SocialAccountOutputDTO

type FindSocialAccountsByUserIdUseCase struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewFindSocialAccountsByUserIdUseCase(
	socialAccountRepo repository.SocialAccountRepository,
) *FindSocialAccountsByUserIdUseCase {
	return &FindSocialAccountsByUserIdUseCase{
		SocialAccountRepository: socialAccountRepo,
	}
}

func (s *FindSocialAccountsByUserIdUseCase) Execute(input *FindSocialAccountsByUserIdInputDTO) (*FindSocialAccountsByUserIdOutputDTO, error) {
	socialAccounts, err := s.SocialAccountRepository.FindSocialAccountsByUserId(input.UserId)
	if err != nil {
		return nil, err
	}
	output := make(FindSocialAccountsByUserIdOutputDTO, len(socialAccounts))
	for i, socialAccount := range socialAccounts {
		output[i] = &SocialAccountOutputDTO{
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
