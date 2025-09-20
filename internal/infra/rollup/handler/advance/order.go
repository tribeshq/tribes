package advance

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/order"
)

type OrderAdvanceHandlers struct {
	OrderRepository    repository.OrderRepository
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewOrderAdvanceHandlers(
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	campaignRepo repository.CampaignRepository,
) *OrderAdvanceHandlers {
	return &OrderAdvanceHandlers{
		OrderRepository:    orderRepo,
		UserRepository:     userRepo,
		CampaignRepository: campaignRepo,
	}
}

func (h *OrderAdvanceHandlers) CreateOrder(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order.CreateOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w, payload: %s", err, string(payload))
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	createOrder := order.NewCreateOrderUseCase(
		h.UserRepository,
		h.OrderRepository,
		h.CampaignRepository,
	)

	res, err := createOrder.Execute(&input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return fmt.Errorf("invalid deposit custom_type, expected ERC20Deposit")
	}

	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Value,
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

func (h *OrderAdvanceHandlers) CancelOrder(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order.CancelOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	cancelOrder := order.NewCancelOrderUseCase(
		h.UserRepository,
		h.OrderRepository,
		h.CampaignRepository,
	)

	res, err := cancelOrder.Execute(&input, metadata)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if err := env.ERC20Transfer(
		common.Address(res.Token),
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
