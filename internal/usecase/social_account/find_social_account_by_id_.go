package social_account

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type FindSocialAccountByIdInputDTO struct {
	SocialAccountId uint `json:"social_account_id" validate:"required"`
}

type FindSocialAccountByIdUseCase struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewFindSocialAccountByIdUseCase(socialAccountRepository repository.SocialAccountRepository) *FindSocialAccountByIdUseCase {
	return &FindSocialAccountByIdUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *FindSocialAccountByIdUseCase) Execute(ctx context.Context, input *FindSocialAccountByIdInputDTO) (*FindSocialAccountOutputDTO, error) {
	socialAccount, err := s.SocialAccountRepository.FindSocialAccountById(ctx, input.SocialAccountId)
	if err != nil {
		return nil, err
	}
	return &FindSocialAccountOutputDTO{
		Id:        socialAccount.Id,
		UserId:    socialAccount.UserId,
		Username:  socialAccount.Username,
		Platform:  string(socialAccount.Platform),
		Proof:     socialAccount.Proof,
		CreatedAt: socialAccount.CreatedAt,
		UpdatedAt: socialAccount.UpdatedAt,
	}, nil
}
