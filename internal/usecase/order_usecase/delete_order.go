package order_usecase

import (
	"context"

	"github.com/tribeshq/tribes/internal/infra/repository"
)

type DeleteOrderInputDTO struct {
	Id uint `json:"id"`
}

type DeleteOrderUseCase struct {
	OrderRepository repository.OrderRepository
}

func NewDeleteOrderUseCase(orderRepository repository.OrderRepository) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *DeleteOrderUseCase) Execute(ctx context.Context, input *DeleteOrderInputDTO) error {
	return c.OrderRepository.DeleteOrder(ctx, input.Id)
}
