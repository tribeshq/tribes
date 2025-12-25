package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestEmergencySuite(t *testing.T) {
	suite.Run(t, new(EmergencySuite))
}

type EmergencySuite struct {
	DCMRollupSuite
}

func (s *EmergencySuite) TestEmergencyERC20Withdraw() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")
	to := common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")

	// Emergency ERC20 withdraw
	emergencyERC20WithdrawInput := []byte(fmt.Sprintf(`{"path":"user/admin/emergency-erc20-withdraw","data":{"to":"%s","token":"%s"}}`, to.Hex(), token.Hex()))
	emergencyERC20WithdrawOutput := s.Tester.Advance(admin, emergencyERC20WithdrawInput)
	s.Len(emergencyERC20WithdrawOutput.DelegateCallVouchers, 1)

	// Verify the delegate call voucher payload
	abiJSON := `[{
		"type":"function",
		"name":"emergencyERC20Withdraw",
		"inputs":[
			{"type":"address"},
			{"type":"address"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["emergencyERC20Withdraw"].Inputs.Unpack(emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
}

func (s *EmergencySuite) TestEmergencyEtherWithdraw() {
	admin := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	to := common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")

	// Emergency ETH withdraw
	emergencyEtherWithdrawInput := []byte(fmt.Sprintf(`{"path":"user/admin/emergency-ether-withdraw","data":{"to":"%s"}}`, to.Hex()))
	emergencyEtherWithdrawOutput := s.Tester.Advance(admin, emergencyEtherWithdrawInput)
	s.Len(emergencyEtherWithdrawOutput.DelegateCallVouchers, 1)

	// Verify the delegate call voucher payload
	abiJSON := `[{
		"type":"function",
		"name":"emergencyETHWithdraw",
		"inputs":[
			{"type":"address"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["emergencyETHWithdraw"].Inputs.Unpack(emergencyEtherWithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(to, unpacked[0].(common.Address))
}
