package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/order"
	"github.com/rollmelette/rollmelette"
)

type OrderInspectHandlers struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewOrderInspectHandlers(
	userRepo repository.UserRepository,
	orderRepo repository.OrderRepository,
) *OrderInspectHandlers {
	return &OrderInspectHandlers{
		UserRepository:  userRepo,
		OrderRepository: orderRepo,
	}
}

func (h *OrderInspectHandlers) FindOrderById(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrderByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	findOrderById := order.NewFindOrderByIdUseCase(h.UserRepository, h.OrderRepository)
	res, err := findOrderById.Execute(&input)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}
	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	env.Report(order)
	return nil
}

func (h *OrderInspectHandlers) FindBidsByCampaignId(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrdersByCampaignIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	findOrdersByCampaignId := order.NewFindOrdersByCampaignIdUseCase(h.UserRepository, h.OrderRepository)
	res, err := findOrdersByCampaignId.Execute(&input)
	if err != nil {
		return fmt.Errorf("failed to find orders by campaign id: %v", err)
	}
	orders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
	env.Report(orders)
	return nil
}

func (h *OrderInspectHandlers) FindAllOrders(env rollmelette.EnvInspector, payload []byte) error {

	findAllOrders := order.NewFindAllOrdersUseCase(h.UserRepository, h.OrderRepository)
	res, err := findAllOrders.Execute()
	if err != nil {
		return fmt.Errorf("failed to find all orders: %w", err)
	}
	allOrders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all orders: %w", err)
	}
	env.Report(allOrders)
	return nil
}

func (h *OrderInspectHandlers) FindOrdersByInvestorAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrdersByInvestorAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	findOrdersByInvestor := order.NewFindOrdersByInvestorAddressUseCase(h.UserRepository, h.OrderRepository)
	res, err := findOrdersByInvestor.Execute(&input)
	if err != nil {
		return fmt.Errorf("failed to find orders by investor: %w", err)
	}
	orders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
	env.Report(orders)
	return nil
}
