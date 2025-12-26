package sqlite

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/domain/entity"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateOrder(input *entity.Order) (*entity.Order, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindOrderById(id uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.Db.First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to find order by ID: %w", err)
	}
	return &order, nil
}

func (r *SQLiteRepository) FindOrdersByIssuanceId(id uint) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.Where("issuance_id = ?", id).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by issuance ID: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) FindOrdersByState(issuanceId uint, state string) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.
		Where("issuance_id = ? AND state = ?", issuanceId, state).
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by state: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) FindOrdersByInvestorAddress(investor types.Address) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.Where("investor_address = ?", investor).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by investor: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) FindAllOrders() ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find all orders: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) UpdateOrder(input *entity.Order) (*entity.Order, error) {
	if err := r.Db.Save(input).Error; err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	order, err := r.FindOrderById(input.Id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *SQLiteRepository) DeleteOrder(id uint) error {
	res := r.Db.Delete(&entity.Order{}, id)
	if res.Error != nil {
		return fmt.Errorf("failed to delete order: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.ErrOrderNotFound
	}
	return nil
}
