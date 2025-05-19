package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/social_account_usecase"
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
	var input social_account_usecase.FindSocialAccountByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findSocialAccountById := social_account_usecase.NewFindSocialAccountByIdUseCase(h.SocialAccountRepository)
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
	var input social_account_usecase.FindSocialAccountsByUserIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findSocialAccountsByUserId := social_account_usecase.NewFindSocialAccountsByUserIdUseCase(h.SocialAccountRepository)
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
