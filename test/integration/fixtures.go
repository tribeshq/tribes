package integration

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Fixtures contains sample data and common expected strings
type Fixtures struct{}

// NewFixtures creates a new fixtures instance
func NewFixtures() *Fixtures {
	return &Fixtures{}
}

// ExpectedUserOutput returns the expected string for user creation
func (f *Fixtures) ExpectedUserOutput(address common.Address, id int, role string, baseTime int64) string {
	return fmt.Sprintf(`user created - {"id":%d,"role":"%s","address":"%s","social_accounts":[],"created_at":%d}`,
		id, role, address.Hex(), baseTime)
}

// ExpectedSocialAccountOutput returns the expected string for social account creation
func (f *Fixtures) ExpectedSocialAccountOutput(userID int, baseTime int64) string {
	return fmt.Sprintf(`social account created - {"id":1,"user_id":%d,"username":"test","platform":"twitter","created_at":%d}`,
		userID, baseTime)
}

// ExpectedCreateCampaignOutput returns the expected string for campaign creation
func (f *Fixtures) ExpectedCreateCampaignOutput(
	token common.Address,
	creator common.Address,
	baseTime int64,
	collateral common.Address,
	badgeAddress common.Address,
	closesAt int64,
	maturityAt int64,
) string {
	return fmt.Sprintf(`campaign created - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","state":"ongoing","orders":[],"created_at":%d,"closes_at":%d,"maturity_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		baseTime, closesAt, maturityAt)
}

// ExpectedCloseCampaignOutput returns the expected string for campaign closure
func (f *Fixtures) ExpectedCloseCampaignOutput(
	token common.Address,
	creator common.Address,
	baseTime int64,
	collateral common.Address,
	badgeAddress common.Address,
	closesAt int64,
	maturityAt int64,
	investor01 common.Address,
	investor02 common.Address,
	investor03 common.Address,
	investor04 common.Address,
	investor05 common.Address,
) string {
	return fmt.Sprintf(`campaign closed - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"closed","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"partially_accepted","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"accepted","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, closesAt,
		investor02.Hex(), baseTime, closesAt,
		investor03.Hex(), baseTime, closesAt,
		investor04.Hex(), baseTime, closesAt,
		investor05.Hex(), baseTime, closesAt,
		investor01.Hex(), baseTime, closesAt,
		baseTime, closesAt, maturityAt, closesAt)
}

// ExpectedSettleCampaignOutput returns the expected string for campaign settlement
func (f *Fixtures) ExpectedSettleCampaignOutput(
	token common.Address,
	creator common.Address,
	baseTime int64,
	collateral common.Address,
	badgeAddress common.Address,
	closesAt int64,
	maturityAt int64,
	settledAt int64,
	investor01 common.Address,
	investor02 common.Address,
	investor03 common.Address,
	investor04 common.Address,
	investor05 common.Address,
) string {
	return fmt.Sprintf(`campaign settled - {"id":1,"title":"test","description":"testtesttesttesttest","promotion":"testtesttesttesttest","token":"%s","creator":{"id":3,"role":"creator","address":"%s","social_accounts":[{"id":1,"user_id":3,"username":"test","platform":"twitter","created_at":%d}],"created_at":%d,"updated_at":0},"collateral_address":"%s","collateral_amount":"10000","badge_address":"%s","debt_issued":"100000","max_interest_rate":"10","total_obligation":"108195","total_raised":"100000","state":"settled","orders":[`+
		`{"id":1,"campaign_id":1,"investor":"%s","amount":"59500","interest_rate":"9","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":2,"campaign_id":1,"investor":"%s","amount":"28000","interest_rate":"8","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":3,"campaign_id":1,"investor":"%s","amount":"2000","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":4,"campaign_id":1,"investor":"%s","amount":"5000","interest_rate":"6","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":5,"campaign_id":1,"investor":"%s","amount":"5500","interest_rate":"4","state":"settled","created_at":%d,"updated_at":%d},`+
		`{"id":6,"campaign_id":1,"investor":"%s","amount":"500","interest_rate":"9","state":"rejected","created_at":%d,"updated_at":%d}],`+
		`"created_at":%d,"closes_at":%d,"maturity_at":%d,"updated_at":%d}`,
		token.Hex(),
		creator.Hex(),
		baseTime,
		baseTime,
		collateral.Hex(),
		badgeAddress.Hex(),
		investor01.Hex(), baseTime, settledAt,
		investor02.Hex(), baseTime, settledAt,
		investor03.Hex(), baseTime, settledAt,
		investor04.Hex(), baseTime, settledAt,
		investor05.Hex(), baseTime, settledAt,
		investor01.Hex(), baseTime, closesAt,
		baseTime, closesAt, maturityAt, settledAt)
}

// ExpectedWithdrawOutput returns the expected string for withdrawal
func (f *Fixtures) ExpectedWithdrawOutput(token common.Address, amount string, user common.Address) string {
	return fmt.Sprintf(`ERC20 withdrawn - token: %s, amount: %s, user: %s`, token.Hex(), amount, user.Hex())
}

// ExpectedBalanceOutput returns the expected string for balance query
func (f *Fixtures) ExpectedBalanceOutput(balance string) string {
	return fmt.Sprintf(`"%s"`, balance)
}

// TestData contains common test data
type TestData struct {
	CampaignTitle       string
	CampaignDescription string
	CampaignPromotion   string
	MaxInterestRate     string
	DebtIssued          string
	CollateralAmount    *big.Int
	OrderAmounts        []*big.Int
	InterestRates       []string
}

// NewTestData creates default test data
func NewTestData() *TestData {
	return &TestData{
		CampaignTitle:       "test",
		CampaignDescription: "testtesttesttesttest",
		CampaignPromotion:   "testtesttesttesttest",
		MaxInterestRate:     "10",
		DebtIssued:          "100000",
		CollateralAmount:    big.NewInt(10000),
		OrderAmounts: []*big.Int{
			big.NewInt(60000),
			big.NewInt(28000),
			big.NewInt(2000),
			big.NewInt(5000),
			big.NewInt(5500),
		},
		InterestRates: []string{"9", "8", "4", "6", "4"},
	}
}
