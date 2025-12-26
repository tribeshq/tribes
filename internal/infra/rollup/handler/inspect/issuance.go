package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/issuance"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

type IssuanceInspectHandlers struct {
	UserRepository     repository.UserRepository
	IssuanceRepository repository.IssuanceRepository
}

func NewIssuanceInspectHandlers(
	userRepo repository.UserRepository,
	issuanceRepo repository.IssuanceRepository,
) *IssuanceInspectHandlers {
	return &IssuanceInspectHandlers{
		UserRepository:     userRepo,
		IssuanceRepository: issuanceRepo,
	}
}

func (h *IssuanceInspectHandlers) FindIssuanceById(env rollmelette.EnvInspector, payload []byte) error {
	var input issuance.FindIssuanceByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	findIssuanceById := issuance.NewFindIssuanceByIdUseCase(h.UserRepository, h.IssuanceRepository)
	res, err := findIssuanceById.Execute(&input)
	if err != nil {
		return fmt.Errorf("failed to find issuance: %w", err)
	}
	issuance, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal issuance: %w", err)
	}
	env.Report(issuance)
	return nil
}

func (h *IssuanceInspectHandlers) FindAllIssuances(env rollmelette.EnvInspector, payload []byte) error {

	findAllIssuancesUseCase := issuance.NewFindAllIssuancesUseCase(h.UserRepository, h.IssuanceRepository)
	res, err := findAllIssuancesUseCase.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all issuances: %w", err)
	}
	allIssuances, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all issuances: %w", err)
	}
	env.Report(allIssuances)
	return nil
}

func (h *IssuanceInspectHandlers) FindIssuancesByInvestorAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input issuance.FindIssuancesByInvestorAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	findIssuancesByInvestor := issuance.NewFindIssuancesByInvestorAddressUseCase(h.UserRepository, h.IssuanceRepository)
	res, err := findIssuancesByInvestor.Execute(&input)
	if err != nil {
		return fmt.Errorf("failed to find issuances by investor: %w", err)
	}
	issuances, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal issuances: %w", err)
	}
	env.Report(issuances)
	return nil
}

func (h *IssuanceInspectHandlers) FindIssuancesByCreatorAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input issuance.FindIssuancesByCreatorAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	findIssuancesByCreator := issuance.NewFindIssuancesByCreatorAddressUseCase(h.UserRepository, h.IssuanceRepository)
	res, err := findIssuancesByCreator.Execute(&input)
	if err != nil {
		return fmt.Errorf("failed to find issuances by creator: %w", err)
	}
	issuances, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal issuances: %w", err)
	}
	env.Report(issuances)
	return nil
}
