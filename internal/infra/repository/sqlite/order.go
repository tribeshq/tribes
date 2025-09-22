package sqlite

import (
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/pkg/custom_type"
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

func (r *SQLiteRepository) FindOrdersByCampaignId(id uint) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.Where("campaign_id = ?", id).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by campaign ID: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) FindOrdersByState(campaignId uint, state string) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.
		Where("campaign_id = ? AND state = ?", campaignId, state).
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to find orders by state: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) FindOrdersByInvestorAddress(investor custom_type.Address) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.Db.Where("investor = ?", investor).Find(&orders).Error; err != nil {
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
	if err := r.Db.Updates(&input).Error; err != nil {
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

func (r *SQLiteRepository) CreateOrdersBatch(orders []*entity.Order) ([]*entity.Order, error) {
	if err := r.Db.CreateInBatches(orders, 100).Error; err != nil {
		return nil, fmt.Errorf("failed to create orders in batch: %w", err)
	}
	return orders, nil
}

func (r *SQLiteRepository) UpdateOrdersBatch(orders []*entity.Order) ([]*entity.Order, error) {
	tx := r.Db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	for _, order := range orders {
		if err := tx.Updates(order).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order in batch: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit batch update: %w", err)
	}

	return orders, nil
}
