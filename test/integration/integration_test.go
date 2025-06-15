package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
)

const TestTimeout = 180 * time.Second // Increased timeout to 3 minutes

// TestTribesIntegrationSuite runs the integration test suite
func TestTribesIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TribesIntegrationSuite))
}

// TribesIntegrationSuite contains the test suite for Tribes integration tests
type TribesIntegrationSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
	group  *errgroup.Group
}

// SetupTest initializes the test environment
func (s *TribesIntegrationSuite) SetupTest() {
	fmt.Println("Starting test setup...")
	s.ctx, s.cancel = context.WithTimeout(context.Background(), TestTimeout)
	s.group, s.ctx = errgroup.WithContext(s.ctx)

	// Stop any running instances first
	fmt.Println("Stopping any running instances...")
	stop := exec.CommandContext(s.ctx, "cartesi", "rollups", "stop")
	if err := stop.Run(); err != nil {
		fmt.Printf("Warning: error stopping rollups: %v\n", err)
	}

	// Start the rollups node
	fmt.Println("Starting rollups node...")
	start := exec.CommandContext(s.ctx, "cartesi", "rollups", "start")
	start.Stdout = os.Stdout
	if err := start.Run(); err != nil {
		s.T().Fatalf("Failed to start rollups: %v", err)
	}

	// Wait for services to be ready
	fmt.Println("Waiting 5 seconds for services to be ready...")
	time.Sleep(5 * time.Second)

	// Check rollups status
	fmt.Println("Checking rollups status...")
	status := exec.CommandContext(s.ctx, "cartesi", "rollups", "status")
	outStatus := NewNotifyWriter(os.Stdout, "cartesi-rollups is running")
	status.Stdout = outStatus
	s.group.Go(func() error { return status.Run() })

	select {
	case <-outStatus.Ready():
		fmt.Println("Rollups node is running")
	case <-s.ctx.Done():
		s.T().Fatal("timeout waiting for rollups to start")
	}

	// Deploy the application
	fmt.Println("Starting deployment...")
	deploy := exec.CommandContext(
		s.ctx,
		"cartesi", "rollups", "deploy",
		"--rpc-url", "http://127.0.0.1:8080/anvil",
		"--mnemonic", "test test test test test test test test test test test junk",
		"--mnemonic-index", "7",
		"--authority-owner", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"--application-owner", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"--name", "tribes",
	)
	deploy.Stdout = os.Stdout
	if err := deploy.Run(); err != nil {
		s.T().Fatalf("Failed to deploy: %v", err)
	}

	// Wait for deployment to be ready
	fmt.Println("Waiting 5 seconds for deployment to be ready...")
	time.Sleep(5 * time.Second)

	// Verify deployment by checking status
	fmt.Println("Verifying deployment...")
	status = exec.CommandContext(s.ctx, "cartesi", "rollups", "status")
	status.Stdout = os.Stdout
	if err := status.Run(); err != nil {
		s.T().Fatalf("Failed to check status: %v", err)
	}

	fmt.Println("Setup completed successfully")
}

// TearDownTest cleans up the test environment
func (s *TribesIntegrationSuite) TearDownTest() {
	fmt.Println("Cleaning up test environment...")
	s.cancel()
	err := exec.Command("cartesi", "rollups", "stop").Run()
	s.NoError(err)
}

// TestAdvance tests the advance functionality
func (s *TribesIntegrationSuite) TestAdvance() {
}

// TestInspect tests the inspect functionality
func (s *TribesIntegrationSuite) TestInspect() {
}
