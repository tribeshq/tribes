package retry

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tribeshq/tribes/pkg/service"
)

type RetrySuite struct {
	suite.Suite
}

func TestRetrySuite(t *testing.T) {
	suite.Run(t, new(RetrySuite))
}

func (s *RetrySuite) SetupSuite()    {}
func (s *RetrySuite) TearDownSuite() {}

func (s *RetrySuite) TestRetry() {
	simpleMock := &SimpleMock{}

	simpleMock.On(
		"execute",
		mock.Anything).
		Once().
		Return(0, fmt.Errorf("An error"))

	simpleMock.On(
		"execute",
		mock.Anything).
		Return(0, nil)

	logger := service.NewLogger(slog.LevelDebug, true)
	_, err := CallFunctionWithRetryPolicy(simpleMock.execute, 0, logger, 3, 1*time.Millisecond, "TEST")
	s.Require().Nil(err)

	simpleMock.AssertNumberOfCalls(s.T(), "execute", 2)

}

func (s *RetrySuite) TestRetryMaxRetries() {

	simpleMock := &SimpleMock{}
	simpleMock.On(
		"execute",
		mock.Anything).
		Return(0, fmt.Errorf("An error"))

	logger := service.NewLogger(slog.LevelDebug, true)
	_, err := CallFunctionWithRetryPolicy(simpleMock.execute, 0, logger, 3, 1*time.Millisecond, "TEST")
	s.Require().NotNil(err)

	simpleMock.AssertNumberOfCalls(s.T(), "execute", 4)

}

type SimpleMock struct {
	mock.Mock
}

func (m *SimpleMock) execute(
	arg int,
) (int, error) {
	args := m.Called(arg)
	return args.Get(0).(int), args.Error(1)
}
