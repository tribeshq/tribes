package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/usecase/order_usecase"
)

type OrderAdvanceHandlers struct {
	UserRepository         entity.UserRepository
	OrderRepository        entity.OrderRepository
	CrowdfundingRepository entity.CrowdfundingRepository
	ContractRepository     entity.ContractRepository
}

func NewOrderAdvanceHandlers(
	userRepository entity.UserRepository,
	orderRepository entity.OrderRepository,
	contractRepository entity.ContractRepository,
	crowdfundingRepository entity.CrowdfundingRepository,
) *OrderAdvanceHandlers {
	return &OrderAdvanceHandlers{
		UserRepository:         userRepository,
		OrderRepository:        orderRepository,
		CrowdfundingRepository: crowdfundingRepository,
		ContractRepository:     contractRepository,
	}
}

func (h *OrderAdvanceHandlers) CreateOrderHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order_usecase.CreateOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w, payload: %s", err, string(payload))
	}

	ctx := context.Background()
	createOrder := order_usecase.NewCreateOrderUseCase(
		h.UserRepository,
		h.OrderRepository,
		h.ContractRepository,
		h.CrowdfundingRepository,
	)

	res, err := createOrder.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return fmt.Errorf("invalid deposit type, expected ERC20Deposit")
	}

	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Amount,
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("order created - "), order...))
	return nil
}

func (h *OrderAdvanceHandlers) CancelOrderHandler(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order_usecase.CancelOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	cancelOrder := order_usecase.NewCancelOrderUseCase(
		h.UserRepository,
		h.OrderRepository,
		h.CrowdfundingRepository,
	)

	res, err := cancelOrder.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	contract, err := h.ContractRepository.FindContractBySymbol(ctx, "STABLECOIN")
	if err != nil {
		return fmt.Errorf("failed to find stablecoin contract: %w", err)
	}

	if err := env.ERC20Transfer(
		common.Address(contract.Address),
		env.AppAddress(),
		metadata.MsgSender,
		res.Amount.ToBig(),
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("order canceled - "), order...))
	return nil
}
