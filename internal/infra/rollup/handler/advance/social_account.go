package advance

import (
	"encoding/json"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/social_account"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type SocialAccountAdvanceHandlers struct {
	UserRepository          repository.UserRepository
	SocialAccountRepository repository.SocialAccountRepository
}

func NewSocialAccountAdvanceHandlers(
	userRepo repository.UserRepository,
	socialAccountRepo repository.SocialAccountRepository,
) *SocialAccountAdvanceHandlers {
	return &SocialAccountAdvanceHandlers{
		UserRepository:          userRepo,
		SocialAccountRepository: socialAccountRepo,
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

	createSocialAccount := social_account.NewCreateSocialAccountUseCase(s.UserRepository, s.SocialAccountRepository)
	res, err := createSocialAccount.Execute(&input, &metadata)
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

	deleteSocialAccount := social_account.NewDeleteSocialAccountUseCase(s.SocialAccountRepository)
	err = deleteSocialAccount.Execute(&input)
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
