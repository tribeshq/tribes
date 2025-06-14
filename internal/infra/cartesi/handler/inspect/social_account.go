package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/social_account"
)

type SocialAccountInspectHandlers struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewSocialAccountInspectHandlers(socialAccountRepository repository.SocialAccountRepository) *SocialAccountInspectHandlers {
	return &SocialAccountInspectHandlers{
		SocialAccountRepository: socialAccountRepository,
	}
}

func (h *SocialAccountInspectHandlers) FindSocialAccountById(env rollmelette.EnvInspector, payload []byte) error {
	var input social_account.FindSocialAccountByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findSocialAccountById := social_account.NewFindSocialAccountByIdUseCase(h.SocialAccountRepository)
	res, err := findSocialAccountById.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find social account: %w", err)
	}
	socialAccount, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("")
	}
	env.Report(socialAccount)
	return nil
}

func (h *SocialAccountInspectHandlers) FindSocialAccountsByUserId(env rollmelette.EnvInspector, payload []byte) error {
	var input social_account.FindSocialAccountsByUserIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findSocialAccountsByUserId := social_account.NewFindSocialAccountsByUserIdUseCase(h.SocialAccountRepository)
	res, err := findSocialAccountsByUserId.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find social accounts: %w", err)
	}
	socialAccounts, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal social accounts: %w", err)
	}
	env.Report(socialAccounts)
	return nil
}
