package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/crowdfunding_usecase"
)

type CrowdfundingInspectHandlers struct {
	CrowdfundingRepository repository.CrowdfundingRepository
}

func NewCrowdfundingInspectHandlers(crowdfundingRepository repository.CrowdfundingRepository) *CrowdfundingInspectHandlers {
	return &CrowdfundingInspectHandlers{
		CrowdfundingRepository: crowdfundingRepository,
	}
}

func (h *CrowdfundingInspectHandlers) FindCrowdfundingById(env rollmelette.EnvInspector, payload []byte) error {
	var input crowdfunding_usecase.FindCrowdfundingByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findCrowdfundingById := crowdfunding_usecase.NewFindCrowdfundingByIdUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingById.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find crowdfunding: %w", err)
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal crowdfunding: %w", err)
	}
	env.Report(crowdfunding)
	return nil
}

func (h *CrowdfundingInspectHandlers) FindAllCrowdfundings(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllCrowdfundingsUseCase := crowdfunding_usecase.NewFindAllCrowdfundingsUseCase(h.CrowdfundingRepository)
	res, err := findAllCrowdfundingsUseCase.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all crowdfundings: %w", err)
	}
	allCrowdfundings, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all crowdfundings: %w", err)
	}
	env.Report(allCrowdfundings)
	return nil
}

func (h *CrowdfundingInspectHandlers) FindCrowdfundingsByInvestor(env rollmelette.EnvInspector, payload []byte) error {
	var input crowdfunding_usecase.FindCrowdfundingsByInvestorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findCrowdfundingsByInvestor := crowdfunding_usecase.NewFindCrowdfundingsByInvestorUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingsByInvestor.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find crowdfundings by investor: %w", err)
	}
	crowdfundings, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal crowdfundings: %w", err)
	}
	env.Report(crowdfundings)
	return nil
}

func (h *CrowdfundingInspectHandlers) FindCrowdfundingsByCreator(env rollmelette.EnvInspector, payload []byte) error {
	var input crowdfunding_usecase.FindCrowdfundingsByCreatorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findCrowdfundingsByCreator := crowdfunding_usecase.NewFindCrowdfundingsByCreatorUseCase(h.CrowdfundingRepository)
	res, err := findCrowdfundingsByCreator.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find crowdfundings by creator: %w", err)
	}
	crowdfundings, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal crowdfundings: %w", err)
	}
	env.Report(crowdfundings)
	return nil
}
