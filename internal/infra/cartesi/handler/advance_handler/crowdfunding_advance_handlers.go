package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/contract_usecase"
	"github.com/tribeshq/tribes/internal/usecase/crowdfunding_usecase"
)

type CrowdfundingAdvanceHandlers struct {
	OrderRepository         entity.OrderRepository
	UserRepository          entity.UserRepository
	SocialAccountRepository entity.SocialAccountRepository
	CrowdfundingRepository  entity.CrowdfundingRepository
	ContractRepository      entity.ContractRepository
}

func NewCrowdfundingAdvanceHandlers(
	orderRepository entity.OrderRepository,
	userRepository entity.UserRepository,
	socialAccountRepository entity.SocialAccountRepository,
	crowdfundingRepository entity.CrowdfundingRepository,
	contractRepository entity.ContractRepository,
) *CrowdfundingAdvanceHandlers {
	return &CrowdfundingAdvanceHandlers{
		OrderRepository:         orderRepository,
		UserRepository:          userRepository,
		SocialAccountRepository: socialAccountRepository,
		CrowdfundingRepository:  crowdfundingRepository,
		ContractRepository:      contractRepository,
	}
}

func (h *CrowdfundingAdvanceHandlers) CreateCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.CreateCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	createCrowdfunding := crowdfunding_usecase.NewCreateCrowdfundingUseCase(
		h.UserRepository,
		h.ContractRepository,
		h.SocialAccountRepository,
		h.CrowdfundingRepository,
	)

	res, err := createCrowdfunding.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create crowdfunding: %w", err)
	}

	erc20Deposit := deposit.(*rollmelette.ERC20Deposit)
	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Amount,
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("crowdfunding created - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) CloseCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.CloseCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	closeCrowdfunding := crowdfunding_usecase.NewCloseCrowdfundingUseCase(h.CrowdfundingRepository, h.OrderRepository)
	res, err := closeCrowdfunding.Execute(ctx, &input, metadata)
	if err != nil && res == nil {
		return fmt.Errorf("failed to close crowdfunding: %w", err)
	}

	// Find stablecoin contract once
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findContractBySymbol.Execute(ctx, &contract_usecase.FindContractBySymbolInputDTO{
		Symbol: "STABLECOIN",
	})
	if err != nil {
		return fmt.Errorf("failed to find stablecoin contract: %w", err)
	}

	// Reuse variables for calculations
	quotes := new(uint256.Int)
	stablecoinAddr := common.Address(stablecoin.Address)
	tokenAddr := common.Address(res.Token)

	// Process orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateRejected {
			if err = env.ERC20Transfer(
				stablecoinAddr,
				env.AppAddress(),
				common.Address(order.Investor),
				order.Amount.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer rejected order: %w", err)
			}
		} else {
			// Calculate quotes
			quotes.Mul(res.Collateral, order.Amount)
			quotes.Div(quotes, res.DebtIssued)

			if err = env.ERC20Transfer(
				tokenAddr,
				env.AppAddress(),
				common.Address(order.Investor),
				quotes.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer accepted order: %w", err)
			}
		}
	}

	// Transfer debt issued to creator
	if err = env.ERC20Transfer(
		stablecoinAddr,
		env.AppAddress(),
		common.Address(res.Creator),
		res.DebtIssued.ToBig(),
	); err != nil {
		return fmt.Errorf("failed to transfer debt issued: %w", err)
	}

	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte(fmt.Sprintf("crowdfunding %v - ", res.State)), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) SettleCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.SettleCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	settleCrowdfunding := crowdfunding_usecase.NewSettleCrowdfundingUseCase(
		h.UserRepository,
		h.CrowdfundingRepository,
		h.ContractRepository,
		h.OrderRepository,
	)

	res, err := settleCrowdfunding.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to settle crowdfunding: %w", err)
	}

	// Find stablecoin contract once
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	contract, err := findContractBySymbol.Execute(ctx, &contract_usecase.FindContractBySymbolInputDTO{
		Symbol: "STABLECOIN",
	})
	if err != nil {
		return fmt.Errorf("failed to find stablecoin contract: %w", err)
	}

	// Reuse variables for calculations
	interest := new(uint256.Int)
	contractAddr := common.Address(contract.Address)
	creatorAddr := common.Address(res.Creator)

	// Process settled orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettled {
			// Calculate interest
			interest.Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))

			// Calculate total payment
			totalPayment := new(uint256.Int).Add(order.Amount, interest)

			if err := env.ERC20Transfer(
				contractAddr,
				creatorAddr,
				common.Address(order.Investor),
				totalPayment.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer settled order: %w", err)
			}
		}
	}

	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("crowdfunding settled - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) UpdateCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.UpdateCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	updateCrowdfunding := crowdfunding_usecase.NewUpdateCrowdfundingUseCase(h.CrowdfundingRepository)
	res, err := updateCrowdfunding.Execute(ctx, input, metadata)
	if err != nil {
		return fmt.Errorf("failed to update crowdfunding: %w", err)
	}

	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("crowdfunding updated - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) DeleteCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.DeleteCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	deleteCrowdfunding := crowdfunding_usecase.NewDeleteCrowdfundingUseCase(h.CrowdfundingRepository)
	if err := deleteCrowdfunding.Execute(ctx, &input); err != nil {
		return fmt.Errorf("failed to delete crowdfunding: %w", err)
	}

	crowdfunding, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal input: %w", err)
	}

	env.Notice(append([]byte("crowdfunding deleted - "), crowdfunding...))
	return nil
}
