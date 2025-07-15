package social_account

import (
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindSocialAccountByIdInputDTO struct {
	SocialAccountId uint `json:"social_account_id" validate:"required"`
}

type FindSocialAccountByIdUseCase struct {
	socialAccountRepository repository.SocialAccountRepository
}

func NewFindSocialAccountByIdUseCase(
	socialAccountRepo repository.SocialAccountRepository,
) *FindSocialAccountByIdUseCase {
	return &FindSocialAccountByIdUseCase{
		socialAccountRepository: socialAccountRepo,
	}
}

func (s *FindSocialAccountByIdUseCase) Execute(input *FindSocialAccountByIdInputDTO) (*SocialAccountOutputDTO, error) {
	socialAccount, err := s.socialAccountRepository.FindSocialAccountById(input.SocialAccountId)
	if err != nil {
		return nil, err
	}
	return &SocialAccountOutputDTO{
		Id:        socialAccount.Id,
		UserId:    socialAccount.UserId,
		Username:  socialAccount.Username,
		Platform:  string(socialAccount.Platform),
		CreatedAt: socialAccount.CreatedAt,
		UpdatedAt: socialAccount.UpdatedAt,
	}, nil
}
