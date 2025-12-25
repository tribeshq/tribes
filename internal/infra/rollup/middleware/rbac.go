package middleware

import (
	"fmt"

	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/infra/repository"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/internal/usecase/user"
	"github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/router"
	types "github.com/2025-2A-T20-G91-INTERNO/src/rollup/pkg/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
)

type RBACFactory struct {
	UserRepository repository.UserRepository
}

func NewRBACFactory(
	userRepo repository.UserRepository,
) *RBACFactory {
	return &RBACFactory{
		UserRepository: userRepo,
	}
}

func (f *RBACFactory) Create(roles []string) router.Middleware {
	return func(handler any) any {
		switch h := handler.(type) {
		case router.AdvanceHandlerFunc:
			return router.AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
				var address types.Address

				// Get the sender address from either ERC20 deposit or metadata
				erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
				if ok {
					address = types.Address(erc20Deposit.Sender)
				} else {
					address = types.Address(metadata.MsgSender)
				}

				// Find user and check roles
				findUserByAddress := user.NewFindUserByAddressUseCase(f.UserRepository)
				user, err := findUserByAddress.Execute(&user.FindUserByAddressInputDTO{
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
					return fmt.Errorf("user %s lacks required permissions: %v", common.Address(user.Address), roles)
				}

				return h(env, metadata, deposit, payload)
			})
		case router.InspectHandlerFunc:
			return h
		default:
			return handler
		}
	}
}

func (f *RBACFactory) AdminOnly() router.Middleware {
	return f.Create([]string{"admin"})
}

func (f *RBACFactory) VerifierOnly() router.Middleware {
	return f.Create([]string{"verifier"})
}

func (f *RBACFactory) InvestorOnly() router.Middleware {
	return f.Create([]string{"investor"})
}

func (f *RBACFactory) CreatorOnly() router.Middleware {
	return f.Create([]string{"creator"})
}
