package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/order_usecase"
)

type OrderInspectHandlers struct {
	OrderRepository repository.OrderRepository
}

func NewOrderInspectHandlers(orderRepository repository.OrderRepository) *OrderInspectHandlers {
	return &OrderInspectHandlers{
		OrderRepository: orderRepository,
	}
}

func (h *OrderInspectHandlers) FindOrderById(env rollmelette.EnvInspector, payload []byte) error {
	var input order_usecase.FindOrderByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrderById := order_usecase.NewFindOrderByIdUseCase(h.OrderRepository)
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

func (h *OrderInspectHandlers) FindBisdByCrowdfundingId(env rollmelette.EnvInspector, payload []byte) error {
	var input order_usecase.FindOrdersByCrowdfundingIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrdersByCrowdfundingId := order_usecase.NewFindOrdersByCrowdfundingIdUseCase(h.OrderRepository)
	res, err := findOrdersByCrowdfundingId.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find orders by crowdfunding id: %v", err)
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
	findAllOrders := order_usecase.NewFindAllOrdersUseCase(h.OrderRepository)
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

func (h *OrderInspectHandlers) FindOrdersByInvestor(env rollmelette.EnvInspector, payload []byte) error {
	var input order_usecase.FindOrdersByInvestorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrdersByInvestor := order_usecase.NewFindOrdersByInvestorUseCase(h.OrderRepository)
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
