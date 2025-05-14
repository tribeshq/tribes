package middleware

import (
	"context"
	"fmt"

	"github.com/rollmelette/rollmelette"
	"github.com/tribeshq/tribes/internal/infra/repository"
	"github.com/tribeshq/tribes/internal/usecase/user_usecase"
	"github.com/tribeshq/tribes/pkg/custom_type"
	"github.com/tribeshq/tribes/pkg/rollups_router"
)

type RBACFactory struct {
	userRepository repository.UserRepository
}

func NewRBACFactory(userRepository repository.UserRepository) *RBACFactory {
	return &RBACFactory{
		userRepository: userRepository,
	}
}

func (f *RBACFactory) Create(roles []string) rollups_router.Middleware {
	return func(handler interface{}) interface{} {
		switch h := handler.(type) {
		case rollups_router.AdvanceHandlerFunc:
			return rollups_router.AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
				var address custom_type.Address
				ctx := context.Background()

				// Get the sender address from either ERC20 deposit or metadata
				erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
				if ok {
					address = custom_type.Address(erc20Deposit.Sender)
				} else {
					address = custom_type.Address(metadata.MsgSender)
				}

				// Find user and check roles
				findUserByAddress := user_usecase.NewFindUserByAddressUseCase(f.userRepository)
				user, err := findUserByAddress.Execute(ctx, &user_usecase.FindUserByAddressInputDTO{
					Address: address,
				})
				if err != nil {
					return err
				}

				// Check if user has any of the required roles
				var hasRole bool
				for _, role := range roles {
					if user.Role == role {
						hasRole = true
						break
					}
				}
				if !hasRole {
					return fmt.Errorf("user with address: %v does not have necessary permissions: %v", user.Address, roles)
				}

				return h(env, metadata, deposit, payload)
			})
		case rollups_router.InspectHandlerFunc:
			return h
		default:
			return handler
		}
	}
}

func (f *RBACFactory) AdminOnly() rollups_router.Middleware {
	return f.Create([]string{"admin"})
}

func (f *RBACFactory) InvestorOnly() rollups_router.Middleware {
	return f.Create([]string{"qualified_investor", "non_qualified_investor"})
}

func (f *RBACFactory) CreatorOnly() rollups_router.Middleware {
	return f.Create([]string{"creator"})
}
