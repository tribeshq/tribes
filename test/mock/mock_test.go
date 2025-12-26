package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestIntegration runs all integration test suites
// This serves as a convenient way to run all integration tests at once
func TestIntegration(t *testing.T) {
	// Run all test suites
	t.Run("User", func(t *testing.T) {
		suite.Run(t, new(UserSuite))
	})
	t.Run("SocialAccount", func(t *testing.T) {
		suite.Run(t, new(SocialAccountSuite))
	})
	t.Run("Order", func(t *testing.T) {
		suite.Run(t, new(OrderSuite))
	})
	t.Run("Issuance", func(t *testing.T) {
		suite.Run(t, new(IssuanceSuite))
	})
	t.Run("Emergency", func(t *testing.T) {
		suite.Run(t, new(EmergencySuite))
	})
}
