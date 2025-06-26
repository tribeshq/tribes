package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	campaign "github.com/tribeshq/tribes/internal/usecase/campaign"
)

type CampaignInspectHandlers struct {
	CampaignRepository repository.CampaignRepository
}

func NewCampaignInspectHandlers(campaignRepository repository.CampaignRepository) *CampaignInspectHandlers {
	return &CampaignInspectHandlers{
		CampaignRepository: campaignRepository,
	}
}

func (h *CampaignInspectHandlers) FindCampaignById(env rollmelette.EnvInspector, payload []byte) error {
	var input campaign.FindCampaignByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findCampaignById := campaign.NewFindCampaignByIdUseCase(h.CampaignRepository)
	res, err := findCampaignById.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find campaign: %w", err)
	}
	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign: %w", err)
	}
	env.Report(campaign)
	return nil
}

func (h *CampaignInspectHandlers) FindAllCampaigns(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllCampaignsUseCase := campaign.NewFindAllCampaignsUseCase(h.CampaignRepository)
	res, err := findAllCampaignsUseCase.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all campaigns: %w", err)
	}
	allCampaigns, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all campaigns: %w", err)
	}
	env.Report(allCampaigns)
	return nil
}

func (h *CampaignInspectHandlers) FindCampaignsByInvestor(env rollmelette.EnvInspector, payload []byte) error {
	var input campaign.FindCampaignsByInvestorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findCampaignsByInvestor := campaign.NewFindCampaignsByInvestorUseCase(h.CampaignRepository)
	res, err := findCampaignsByInvestor.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find campaigns by investor: %w", err)
	}
	campaigns, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal campaigns: %w", err)
	}
	env.Report(campaigns)
	return nil
}

func (h *CampaignInspectHandlers) FindCampaignsByCreator(env rollmelette.EnvInspector, payload []byte) error {
	var input campaign.FindCampaignsByCreatorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findCampaignsByCreator := campaign.NewFindCampaignsByCreatorUseCase(h.CampaignRepository)
	res, err := findCampaignsByCreator.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find campaigns by creator: %w", err)
	}
	campaigns, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal campaigns: %w", err)
	}
	env.Report(campaigns)
	return nil
}
