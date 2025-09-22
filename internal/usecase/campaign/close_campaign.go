package campaign

import (
	"fmt"
	"sort"

	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/domain/entity"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user"
	"github.com/tribeshq/tribes/pkg/custom_type"
)

type CloseCampaignInputDTO struct {
	CreatorAddress custom_type.Address `json:"creator_address" validate:"required"`
}

type CloseCampaignOutputDTO struct {
	Id                uint                `json:"id"`
	Title             string              `json:"title,omitempty"`
	Description       string              `json:"description,omitempty"`
	Promotion         string              `json:"promotion,omitempty"`
	Token             custom_type.Address `json:"token,omitempty"`
	Creator           *user.UserOutputDTO `json:"creator,omitempty"`
	CollateralAddress custom_type.Address `json:"collateral_address,omitempty"`
	CollateralAmount  *uint256.Int        `json:"collateral_amount,omitempty"`
	BadgeAddress      custom_type.Address `json:"badge_address,omitempty"`
	DebtIssued        *uint256.Int        `json:"debt_issued,omitempty"`
	MaxInterestRate   *uint256.Int        `json:"max_interest_rate,omitempty"`
	TotalObligation   *uint256.Int        `json:"total_obligation,omitempty"`
	TotalRaised       *uint256.Int        `json:"total_raised,omitempty"`
	State             string              `json:"state,omitempty"`
	Orders            []*entity.Order     `json:"orders,omitempty"`
	CreatedAt         int64               `json:"created_at,omitempty"`
	ClosesAt          int64               `json:"closes_at,omitempty"`
	MaturityAt        int64               `json:"maturity_at,omitempty"`
	UpdatedAt         int64               `json:"updated_at,omitempty"`
}

type CloseCampaignUseCase struct {
	UserRepository     repository.UserRepository
	OrderRepository    repository.OrderRepository
	CampaignRepository repository.CampaignRepository
}

func NewCloseCampaignUseCase(userRepo repository.UserRepository, campaignRepo repository.CampaignRepository, orderRepo repository.OrderRepository) *CloseCampaignUseCase {
	return &CloseCampaignUseCase{
		UserRepository:     userRepo,
		CampaignRepository: campaignRepo,
		OrderRepository:    orderRepo,
	}
}

func (u *CloseCampaignUseCase) Execute(input *CloseCampaignInputDTO, metadata rollmelette.Metadata) (*CloseCampaignOutputDTO, error) {
	// -------------------------------------------------------------------------
	// 1. Find ongoing campaign for the creator
	// -------------------------------------------------------------------------
	ongoingCampaign, err := u.CampaignRepository.FindOngoingCampaignByCreatorAddress(input.CreatorAddress)
	if err != nil {
		return nil, err
	}
	if ongoingCampaign == nil {
		return nil, fmt.Errorf("no ongoing campaign found, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 2. Validate campaign expiration
	// -------------------------------------------------------------------------
	if metadata.BlockTimestamp < ongoingCampaign.ClosesAt {
		return nil, fmt.Errorf("campaign not expired yet, cannot close it")
	}

	// -------------------------------------------------------------------------
	// 3. Fetch creator info
	// -------------------------------------------------------------------------
	creator, err := u.UserRepository.FindUserByAddress(ongoingCampaign.Creator)
	if err != nil {
		return nil, fmt.Errorf("error finding creator: %w", err)
	}

	// -------------------------------------------------------------------------
	// 4. Fetch and sort campaign orders
	// -------------------------------------------------------------------------
	orders, err := u.OrderRepository.FindOrdersByCampaignId(ongoingCampaign.Id)
	if err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].InterestRate.Cmp(orders[j].InterestRate) == 0 {
			return orders[i].Amount.Cmp(orders[j].Amount) > 0
		}
		return orders[i].InterestRate.Cmp(orders[j].InterestRate) < 0
	})

	// -------------------------------------------------------------------------
	// 5. Select winning orders and calculate obligations
	// -------------------------------------------------------------------------
	debtRemaining := new(uint256.Int).Set(ongoingCampaign.DebtIssued)
	totalCollected := uint256.NewInt(0)
	totalObligation := uint256.NewInt(0)

	var ordersToUpdate []*entity.Order
	var ordersToCreate []*entity.Order

	for _, order := range orders {
		if debtRemaining.IsZero() {
			// Reject surplus orders (batch later)
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			ordersToUpdate = append(ordersToUpdate, order)
			continue
		}

		// Accept full or partial order
		acceptAmount := order.Amount
		if debtRemaining.Lt(order.Amount) {
			acceptAmount = new(uint256.Int).Set(debtRemaining)
		}
		interest := new(uint256.Int).Mul(acceptAmount, order.InterestRate)
		interest.Div(interest, uint256.NewInt(100))

		orderObligation := new(uint256.Int).Add(acceptAmount, interest)
		totalCollected.Add(totalCollected, acceptAmount)
		totalObligation.Add(totalObligation, orderObligation)

		if debtRemaining.Cmp(order.Amount) >= 0 {
			order.State = entity.OrderStateAccepted
			debtRemaining.Sub(debtRemaining, order.Amount)
		} else {
			order.State = entity.OrderStatePartiallyAccepted
			// Queue rejected order for batch creation
			rejectedAmount := new(uint256.Int).Sub(order.Amount, acceptAmount)
			ordersToCreate = append(ordersToCreate, &entity.Order{
				CampaignId:   order.CampaignId,
				Investor:     order.Investor,
				Amount:       rejectedAmount,
				InterestRate: order.InterestRate,
				State:        entity.OrderStateRejected,
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    metadata.BlockTimestamp,
			})
			debtRemaining.Clear()
		}
		order.Amount = acceptAmount
		order.UpdatedAt = metadata.BlockTimestamp
		ordersToUpdate = append(ordersToUpdate, order)
	}

	// Batch update and create orders
	if len(ordersToUpdate) > 0 {
		if _, err := u.OrderRepository.UpdateOrdersBatch(ordersToUpdate); err != nil {
			return nil, fmt.Errorf("error batch updating orders: %w", err)
		}
	}
	if len(ordersToCreate) > 0 {
		if _, err := u.OrderRepository.CreateOrdersBatch(ordersToCreate); err != nil {
			return nil, fmt.Errorf("error batch creating orders: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// 6. Check if minimum funding (2/3) was reached
	// -------------------------------------------------------------------------
	twoThirds := new(uint256.Int).Mul(ongoingCampaign.DebtIssued, uint256.NewInt(2))
	twoThirds.Div(twoThirds, uint256.NewInt(3))
	if totalCollected.Lt(twoThirds) {
		// Cancel campaign and reject all orders (batch operation)
		var canceledOrders []*entity.Order
		for _, order := range orders {
			order.State = entity.OrderStateRejected
			order.UpdatedAt = metadata.BlockTimestamp
			canceledOrders = append(canceledOrders, order)
		}
		if len(canceledOrders) > 0 {
			if _, err := u.OrderRepository.UpdateOrdersBatch(canceledOrders); err != nil {
				return nil, err
			}
		}
		ongoingCampaign.State = entity.CampaignStateCanceled
		ongoingCampaign.UpdatedAt = metadata.BlockTimestamp
		if _, err := u.CampaignRepository.UpdateCampaign(ongoingCampaign); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("campaign canceled due to insufficient funds collected, expected at least 2/3 of the debt issued: %s, got: %s", twoThirds.String(), totalCollected.String())
	}

	// -------------------------------------------------------------------------
	// 7. Close campaign and return result
	// -------------------------------------------------------------------------
	ongoingCampaign.State = entity.CampaignStateClosed
	ongoingCampaign.TotalObligation = totalObligation
	ongoingCampaign.TotalRaised = totalCollected
	ongoingCampaign.UpdatedAt = metadata.BlockTimestamp
	res, err := u.CampaignRepository.UpdateCampaign(ongoingCampaign)
	if err != nil {
		return nil, err
	}


	return &CloseCampaignOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		
		Promotion:   res.Promotion,
		Token:       res.Token,
		Creator: &user.UserOutputDTO{
			Id:             creator.Id,
			Role:           string(creator.Role),
			Address:        creator.Address,
			SocialAccounts: creator.SocialAccounts,
			CreatedAt:      creator.CreatedAt,
			UpdatedAt:      creator.UpdatedAt,
		},
		CollateralAddress: res.CollateralAddress,
		CollateralAmount:  res.CollateralAmount,
		BadgeAddress:      res.BadgeAddress,
		DebtIssued:        res.DebtIssued,
		MaxInterestRate:   res.MaxInterestRate,
		TotalObligation:   res.TotalObligation,
		TotalRaised:       res.TotalRaised,
		Orders:            res.Orders,
		State:             string(res.State),
		ClosesAt:          res.ClosesAt,
		MaturityAt:        res.MaturityAt,
		CreatedAt:         res.CreatedAt,
		UpdatedAt:         res.UpdatedAt,
	}, nil
}
