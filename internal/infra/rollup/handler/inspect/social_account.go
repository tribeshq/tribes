package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/social_account"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type SocialAccountInspectHandlers struct {
	SocialAccountRepository repository.SocialAccountRepository
}

func NewSocialAccountInspectHandlers(
	socialAccountRepo repository.SocialAccountRepository,
) *SocialAccountInspectHandlers {
	return &SocialAccountInspectHandlers{
		SocialAccountRepository: socialAccountRepo,
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

	findSocialAccountById := social_account.NewFindSocialAccountByIdUseCase(h.SocialAccountRepository)
	res, err := findSocialAccountById.Execute(&input)
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

	findSocialAccountsByUserId := social_account.NewFindSocialAccountsByUserIdUseCase(h.SocialAccountRepository)
	res, err := findSocialAccountsByUserId.Execute(&input)
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
