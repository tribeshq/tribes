package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/order"
)

type OrderInspectHandlers struct {
	UserRepository  repository.UserRepository
	OrderRepository repository.OrderRepository
}

func NewOrderInspectHandlers(userRepository repository.UserRepository, orderRepository repository.OrderRepository) *OrderInspectHandlers {
	return &OrderInspectHandlers{
		UserRepository:  userRepository,
		OrderRepository: orderRepository,
	}
}

func (h *OrderInspectHandlers) FindOrderById(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrderByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrderById := order.NewFindOrderByIdUseCase(h.UserRepository, h.OrderRepository)
	res, err := findOrderById.Execute(ctx, &input)
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

	ctx := context.Background()
	findOrdersByCampaignId := order.NewFindOrdersByCampaignIdUseCase(h.UserRepository, h.OrderRepository)
	res, err := findOrdersByCampaignId.Execute(ctx, &input)
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
	ctx := context.Background()
	findAllOrders := order.NewFindAllOrdersUseCase(h.UserRepository, h.OrderRepository)
	res, err := findAllOrders.Execute(ctx)
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

	ctx := context.Background()
	findOrdersByInvestor := order.NewFindOrdersByInvestorAddressUseCase(h.UserRepository, h.OrderRepository)
	res, err := findOrdersByInvestor.Execute(ctx, &input)
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
