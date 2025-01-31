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
	// TODO: remove this check when update to V2
	appAddress, isSet := env.AppAddress()
	if !isSet {
		return fmt.Errorf("no application address defined yet, contact the Tribes support")
	}
	var input *crowdfunding_usecase.CreateCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	createCrowdfunding := crowdfunding_usecase.NewCreateCrowdfundingUseCase(h.UserRepository, h.ContractRepository, h.SocialAccountRepository, h.CrowdfundingRepository)
	res, err := createCrowdfunding.Execute(ctx, input, deposit, metadata)
	if err != nil {
		return err
	}
	if err := env.ERC20Transfer(
		deposit.(*rollmelette.ERC20Deposit).Token,
		deposit.(*rollmelette.ERC20Deposit).Sender,
		appAddress,
		deposit.(*rollmelette.ERC20Deposit).Amount,
	); err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding created - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) CloseCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	// TODO: remove this check when update to V2
	appAddress, isSet := env.AppAddress()
	if !isSet {
		return fmt.Errorf("no application address defined yet, contact Tribes support")
	}

	var input *crowdfunding_usecase.CloseCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	closeCrowdfunding := crowdfunding_usecase.NewCloseCrowdfundingUseCase(h.CrowdfundingRepository, h.OrderRepository)
	res, err := closeCrowdfunding.Execute(ctx, input, metadata)
	if err != nil && res == nil {
		return err
	}

	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	stablecoin, err := findContractBySymbol.Execute(ctx, &contract_usecase.FindContractBySymbolInputDTO{
		Symbol: "STABLECOIN",
	})
	if err != nil {
		return err
	}

	// Return the funds to investors with rejected orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateRejected {
			if err = env.ERC20Transfer(
				common.Address(stablecoin.Address),
				appAddress,
				common.Address(order.Investor),
				order.Amount.ToBig(),
			); err != nil {
				return err
			}
		} else {
			quotes := new(uint256.Int).Div(new(uint256.Int).Mul(res.Amount, order.Amount), res.DebtIssued)
			if err = env.ERC20Transfer(
				common.Address(res.Token),
				appAddress,
				common.Address(order.Investor),
				quotes.ToBig(),
			); err != nil {
				return err
			}
		}
	}

	if err = env.ERC20Transfer(
		common.Address(stablecoin.Address),
		appAddress,
		common.Address(res.Creator),
		res.DebtIssued.ToBig(),
	); err != nil {
		return err
	}

	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte(fmt.Sprintf("crowdfunding %v - ", res.State)), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) SettleCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *crowdfunding_usecase.SettleCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	settleCrowdfunding := crowdfunding_usecase.NewSettleCrowdfundingUseCase(h.UserRepository, h.CrowdfundingRepository, h.ContractRepository, h.OrderRepository)
	res, err := settleCrowdfunding.Execute(ctx, input, deposit, metadata)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	findContractBySymbol := contract_usecase.NewFindContractBySymbolUseCase(h.ContractRepository)
	contract, err := findContractBySymbol.Execute(ctx, &contract_usecase.FindContractBySymbolInputDTO{
		Symbol: "STABLECOIN",
	})
	if err != nil {
		return err
	}
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettled {
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			if err := env.ERC20Transfer(
				common.Address(contract.Address),
				common.Address(res.Creator),
				common.Address(order.Investor),
				new(uint256.Int).Add(order.Amount, interest).ToBig(),
			); err != nil {
				return err
			}
		}
	}
	env.Notice(append([]byte("crowdfunding settled - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) UpdateCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input crowdfunding_usecase.UpdateCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	updateCrowdfunding := crowdfunding_usecase.NewUpdateCrowdfundingUseCase(h.CrowdfundingRepository)
	res, err := updateCrowdfunding.Execute(ctx, input, metadata)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(res)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding updated - "), crowdfunding...))
	return nil
}

func (h *CrowdfundingAdvanceHandlers) DeleteCrowdfundingHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input *crowdfunding_usecase.DeleteCrowdfundingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	ctx := context.Background()
	deleteCrowdfunding := crowdfunding_usecase.NewDeleteCrowdfundingUseCase(h.CrowdfundingRepository)
	err := deleteCrowdfunding.Execute(ctx, input)
	if err != nil {
		return err
	}
	crowdfunding, err := json.Marshal(input)
	if err != nil {
		return err
	}
	env.Notice(append([]byte("crowdfunding deleted - "), crowdfunding...))
	return nil
}
