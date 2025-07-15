package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestIntegration runs all integration test suites
// This serves as a convenient way to run all integration tests at once
func TestIntegration(t *testing.T) {
	// Run all test suites
	t.Run("Campaign", func(t *testing.T) {
		suite.Run(t, new(CampaignSuite))
	})
	t.Run("Emergency", func(t *testing.T) {
		suite.Run(t, new(EmergencySuite))
	})
}
