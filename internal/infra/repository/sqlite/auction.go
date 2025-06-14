package sqlite

import (
	"context"
	"fmt"

	"github.com/tribeshq/tribes/internal/domain/entity"
	. "github.com/tribeshq/tribes/pkg/custom_type"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateAuction(ctx context.Context, input *entity.Auction) (*entity.Auction, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create auction: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindAuctionById(ctx context.Context, id uint) (*entity.Auction, error) {
	var Auction entity.Auction
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		First(&Auction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrAuctionNotFound
		}
		return nil, fmt.Errorf("failed to find auction by id: %w", err)
	}
	return &Auction, nil
}

func (r *SQLiteRepository) FindAllAuctions(ctx context.Context) ([]*entity.Auction, error) {
	var Auctions []*entity.Auction
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		Find(&Auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to find all auctions: %w", err)
	}
	return Auctions, nil
}

func (r *SQLiteRepository) FindAuctionsByInvestor(ctx context.Context, investor Address) ([]*entity.Auction, error) {
	var Auctions []*entity.Auction
	if err := r.Db.WithContext(ctx).
		Joins("JOIN orders ON orders.auction_id = auctions.id").
		Where("orders.investor = ?", investor).
		Preload("Orders").
		Find(&Auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to find Auctions by investor: %w", err)
	}
	return Auctions, nil
}

func (r *SQLiteRepository) FindAuctionsByCreator(ctx context.Context, creator Address) ([]*entity.Auction, error) {
	var Auctions []*entity.Auction
	if err := r.Db.WithContext(ctx).
		Where("creator = ?", creator).
		Preload("Orders").
		Find(&Auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to find auctions by creator: %w", err)
	}
	return Auctions, nil
}

func (r *SQLiteRepository) UpdateAuction(ctx context.Context, input *entity.Auction) (*entity.Auction, error) {
	if err := r.Db.WithContext(ctx).Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update auction: %w", err)
	}
	Auction, err := r.FindAuctionById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return Auction, nil
}
