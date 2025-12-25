package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/campaign"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type CampaignInspectHandlers struct {
	UserRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository
}

func NewCampaignInspectHandlers(
	userRepo repository.UserRepository,
	campaignRepo repository.CampaignRepository,
) *CampaignInspectHandlers {
	return &CampaignInspectHandlers{
		UserRepository:     userRepo,
		campaignRepository: campaignRepo,
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

	findCampaignById := campaign.NewFindCampaignByIdUseCase(h.UserRepository, h.campaignRepository)
	res, err := findCampaignById.Execute(&input)
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

	findAllCampaignsUseCase := campaign.NewFindAllCampaignsUseCase(h.UserRepository, h.campaignRepository)
	res, err := findAllCampaignsUseCase.Execute()
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

func (h *CampaignInspectHandlers) FindCampaignsByInvestorAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input campaign.FindCampaignsByInvestorAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	findCampaignsByInvestor := campaign.NewFindCampaignsByInvestorAddressUseCase(h.UserRepository, h.campaignRepository)
	res, err := findCampaignsByInvestor.Execute(&input)
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

func (h *CampaignInspectHandlers) FindCampaignsByCreatorAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input campaign.FindCampaignsByCreatorAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	findCampaignsByCreator := campaign.NewFindCampaignsByCreatorAddressUseCase(h.UserRepository, h.campaignRepository)
	res, err := findCampaignsByCreator.Execute(&input)
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
