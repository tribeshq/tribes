package social_account

import (
	"context"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
)

type CreateSocialAccountInputDTO struct {
	UserId    uint   `json:"user_id" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Followers uint   `json:"followers" validate:"required"`
	Platform  string `json:"platform" validate:"required"`
	CreatedAt int64  `json:"created_at" validate:"required"`
}

type CreateSocialAccountOutputDTO struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Followers uint   `json:"followers"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
}

type CreateSocialAccountUseCase struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewCreateSocialAccountUseCase(socialAccountRepository repository.SocialAccountRepository) *CreateSocialAccountUseCase {
	return &CreateSocialAccountUseCase{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *CreateSocialAccountUseCase) Execute(ctx context.Context, input *CreateSocialAccountInputDTO) (*CreateSocialAccountOutputDTO, error) {
	socialAccount, err := entity.NewSocialAccount(input.UserId, input.Username, input.Followers, input.Platform, int64(input.CreatedAt))
	if err != nil {
		return nil, err
	}
	socialAccount, err = s.SocialAccountRepository.CreateSocialAccount(ctx, socialAccount)
	if err != nil {
		return nil, err
	}
	return &CreateSocialAccountOutputDTO{
		Id:        socialAccount.Id,
		UserId:    socialAccount.UserId,
		Username:  socialAccount.Username,
		Followers: socialAccount.Followers,
		Platform:  string(socialAccount.Platform),
		CreatedAt: socialAccount.CreatedAt,
	}, nil
}
