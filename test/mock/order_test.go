package integration

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}

type OrderSuite struct {
	DCMRollupSuite
}

func (s *OrderSuite) TestCreateOrder() {
	admin, token, creator, _, verifier, collateral, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// create creator user
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	// verify social account
	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	// create campaign
	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)

	// create investor user
	createInvestorInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createInvestorInput)

	// create order
	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	createOrderOutput := s.Tester.DepositERC20(token, investor01, big.NewInt(10000), createOrderInput)
	s.Len(createOrderOutput.Notices, 1)

	expectedCreateOrderOutput := fmt.Sprintf(`order created - {"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":%d}`,
		investor01,
		baseTime,
		baseTime)
	s.Equal(expectedCreateOrderOutput, string(createOrderOutput.Notices[0].Payload))
}

func (s *OrderSuite) TestFindAllOrders() {
	admin, token, creator, _, verifier, collateral, _, _ := s.setupCommonAddresses()
	investor01, investor02, _, _, _ := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// Setup
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)

	// Create investors
	createInvestorInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createInvestorInput)

	createInvestorInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	s.Tester.Advance(admin, createInvestorInput)

	// Create orders
	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	s.Tester.DepositERC20(token, investor01, big.NewInt(10000), createOrderInput)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	s.Tester.DepositERC20(token, investor02, big.NewInt(20000), createOrderInput)

	// Find all orders
	findAllOrdersInput := []byte(`{"path":"order"}`)
	findAllOrdersOutput := s.Tester.Inspect(findAllOrdersInput)
	s.Len(findAllOrdersOutput.Reports, 1)

	expectedFindAllOrdersOutput := fmt.Sprintf(`[{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":%d,"updated_at":0},{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"20000","interest_rate":"8","state":"pending","created_at":%d,"updated_at":0}]`,
		investor01, baseTime, baseTime,
		investor02, baseTime, baseTime)
	s.Equal(expectedFindAllOrdersOutput, string(findAllOrdersOutput.Reports[0].Payload))
}

func (s *OrderSuite) TestFindOrderById() {
	admin, token, creator, _, verifier, collateral, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// Setup
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)

	createInvestorInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createInvestorInput)

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	s.Tester.DepositERC20(token, investor01, big.NewInt(10000), createOrderInput)

	// Find order by id
	findOrderByIdInput := []byte(`{"path":"order/id","data":{"id":1}}`)
	findOrderByIdOutput := s.Tester.Inspect(findOrderByIdInput)
	s.Len(findOrderByIdOutput.Reports, 1)

	expectedFindOrderByIdOutput := fmt.Sprintf(`{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":%d,"updated_at":0}`,
		investor01, baseTime, baseTime)
	s.Equal(expectedFindOrderByIdOutput, string(findOrderByIdOutput.Reports[0].Payload))
}

func (s *OrderSuite) TestFindOrdersByCampaignId() {
	admin, token, creator, _, verifier, collateral, _, _ := s.setupCommonAddresses()
	investor01, investor02, _, _, _ := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// Setup
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)

	createInvestorInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createInvestorInput)

	createInvestorInput = []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor02))
	s.Tester.Advance(admin, createInvestorInput)

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	s.Tester.DepositERC20(token, investor01, big.NewInt(10000), createOrderInput)

	createOrderInput = []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"8"}}`)
	s.Tester.DepositERC20(token, investor02, big.NewInt(20000), createOrderInput)

	// Find orders by campaign id
	findOrdersByCampaignIdInput := []byte(`{"path":"order/campaign","data":{"campaign_id":1}}`)
	findOrdersByCampaignIdOutput := s.Tester.Inspect(findOrdersByCampaignIdInput)
	s.Len(findOrdersByCampaignIdOutput.Reports, 1)

	expectedFindOrdersByCampaignIdOutput := fmt.Sprintf(`[{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":%d,"updated_at":0},{"id":2,"campaign_id":1,"investor":{"id":5,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"20000","interest_rate":"8","state":"pending","created_at":%d,"updated_at":0}]`,
		investor01, baseTime, baseTime,
		investor02, baseTime, baseTime)
	s.Equal(expectedFindOrdersByCampaignIdOutput, string(findOrdersByCampaignIdOutput.Reports[0].Payload))
}

func (s *OrderSuite) TestFindOrdersByInvestorAddress() {
	admin, token, creator, _, verifier, collateral, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// Setup
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)

	createInvestorInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createInvestorInput)

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	s.Tester.DepositERC20(token, investor01, big.NewInt(10000), createOrderInput)

	// Find orders by investor address
	findOrdersByInvestorAddressInput := []byte(fmt.Sprintf(`{"path":"order/investor","data":{"investor_address":"%s"}}`, investor01))
	findOrdersByInvestorAddressOutput := s.Tester.Inspect(findOrdersByInvestorAddressInput)
	s.Len(findOrdersByInvestorAddressOutput.Reports, 1)

	expectedFindOrdersByInvestorAddressOutput := fmt.Sprintf(`[{"id":1,"campaign_id":1,"investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"10000","interest_rate":"9","state":"pending","created_at":%d,"updated_at":0}]`,
		investor01, baseTime, baseTime)
	s.Equal(expectedFindOrdersByInvestorAddressOutput, string(findOrdersByInvestorAddressOutput.Reports[0].Payload))
}

func (s *OrderSuite) TestCancelOrder() {
	admin, token, creator, _, verifier, collateral, _, _ := s.setupCommonAddresses()
	investor01, _, _, _, _ := s.setupInvestorAddresses()
	baseTime, closesAt, maturityAt := s.setupTimeValues()

	// Setup
	createUserInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"creator"}}`, creator))
	s.Tester.Advance(admin, createUserInput)

	createSocialAccountInput := []byte(fmt.Sprintf(`{"path":"social/verifier/create","data":{"address":"%s","username":"test","platform":"twitter"}}`, creator))
	s.Tester.Advance(verifier, createSocialAccountInput)

	createCampaignInput := []byte(fmt.Sprintf(`{"path":"campaign/creator/create","data":{"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","max_interest_rate":"10","debt_issued":"100000","closes_at":%d,"maturity_at":%d}}`,
		token,
		closesAt,
		maturityAt,
	))
	s.Tester.DepositERC20(collateral, creator, big.NewInt(10000), createCampaignInput)

	createInvestorInput := []byte(fmt.Sprintf(`{"path":"user/admin/create","data":{"address":"%s","role":"investor"}}`, investor01))
	s.Tester.Advance(admin, createInvestorInput)

	createOrderInput := []byte(`{"path": "order/create", "data": {"campaign_id":1,"interest_rate":"9"}}`)
	s.Tester.DepositERC20(token, investor01, big.NewInt(10000), createOrderInput)

	// Cancel order
	cancelOrderInput := []byte(`{"path":"order/cancel","data":{"id":1}}`)
	cancelOrderOutput := s.Tester.Advance(investor01, cancelOrderInput)
	s.Len(cancelOrderOutput.Notices, 1)

	expectedCancelOrderOutput := fmt.Sprintf(`order canceled - {"id":1,"campaign_id":1,"token":"%s","investor":{"id":4,"role":"investor","address":"%s","social_accounts":[],"created_at":%d,"updated_at":0},"amount":"10000","interest_rate":"9","state":"canceled","created_at":%d,"updated_at":%d}`,
		token, investor01, baseTime, baseTime, baseTime)
	s.Equal(expectedCancelOrderOutput, string(cancelOrderOutput.Notices[0].Payload))
}
