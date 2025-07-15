package social_account

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindSocialAccountsByUserIdInputDTO struct {
	UserId uint `json:"user_id" validate:"required"`
}

type FindSocialAccountsByUserIdOutputDTO []*SocialAccountOutputDTO

type FindSocialAccountsByUserIdUseCase struct {
	socialAccountRepository repository.SocialAccountRepository
}

func NewFindSocialAccountsByUserIdUseCase(
	socialAccountRepo repository.SocialAccountRepository,
) *FindSocialAccountsByUserIdUseCase {
	return &FindSocialAccountsByUserIdUseCase{
		socialAccountRepository: socialAccountRepo,
	}
}

func (s *FindSocialAccountsByUserIdUseCase) Execute(input *FindSocialAccountsByUserIdInputDTO) (*FindSocialAccountsByUserIdOutputDTO, error) {
	socialAccounts, err := s.socialAccountRepository.FindSocialAccountsByUserId(input.UserId)
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
