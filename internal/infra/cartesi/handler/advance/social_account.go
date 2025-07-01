package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/social_account"
)

type SocialAccountAdvanceHandlers struct {
	UserRepository          repository.UserRepository
	SocialAccountRepository repository.SocialAccountRepository
}

func NewSocialAccountAdvanceHandlers(userRepository repository.UserRepository, socialAccountRepository repository.SocialAccountRepository) *SocialAccountAdvanceHandlers {
	return &SocialAccountAdvanceHandlers{
		UserRepository:          userRepository,
		SocialAccountRepository: socialAccountRepository,
	}
}

func (s *SocialAccountAdvanceHandlers) CreateSocialAccount(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input social_account.CreateSocialAccountInputDTO
	err := json.Unmarshal(payload, &input)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	createSocialAccount := social_account.NewCreateSocialAccountUseCase(s.UserRepository, s.SocialAccountRepository)
	res, err := createSocialAccount.Execute(ctx, &input, &metadata)
	if err != nil {
		return err
	}
	socialAccount, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("social account created - "), socialAccount...))
	return nil
}

func (s *SocialAccountAdvanceHandlers) DeleteSocialAccount(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input social_account.DeleteSocialAccountInputDTO
	err := json.Unmarshal(payload, &input)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	deleteSocialAccount := social_account.NewDeleteSocialAccountUseCase(s.SocialAccountRepository)
	err = deleteSocialAccount.Execute(ctx, &input)
	if err != nil {
		return err
	}
	socialAccount, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("social account deleted - "), socialAccount...))
	return nil
}
