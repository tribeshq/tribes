// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {Token} from "../src/token/ERC20/Token.sol";
import {MockApplication} from "./mock/MockApplication.sol";
import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";
import {InputBox} from "cartesi-rollups-contracts-2.0.0/src/inputs/InputBox.sol";

contract TestBaseLayerGas is Test {
    Token public paymentToken;
    MockApplication public mockApplication;
    InputBox public inputBox;

    address public creator;

    event InputAdded(address indexed appContract, uint256 indexed index, bytes input);
    event Transfer(address indexed from, address indexed to, uint256 value);

    function setUp() public {
        creator = makeAddr("creator");

        mockApplication = new MockApplication();
        inputBox = new InputBox();
        paymentToken = new Token("Payment Token", "PAY", address(mockApplication));

        vm.prank(address(mockApplication));
        paymentToken.mint(address(mockApplication), 1000000 * 10 ** 18);
    }

    function testGasAddInputCloseIssuance() public {
        string memory payload =
            string(abi.encodePacked('{"path":"issuance/close","data":{"creator":"', _addressToString(creator), '"}}'));

        uint256 gasBefore = gasleft();
        inputBox.addInput(address(mockApplication), bytes(payload));
        uint256 gasUsed = gasBefore - gasleft();

        emit log_named_uint("Gas used (input close issuance):", gasUsed);

        assertEq(inputBox.getNumberOfInputs(address(mockApplication)), 1);
    }

    function testGasExecuteVoucherWithdrawRaisedFunds() public {
        uint256 raisedAmount = 100000 * 10 ** 18;

        uint256 initialCreatorBalance = paymentToken.balanceOf(creator);
        uint256 initialAppBalance = paymentToken.balanceOf(address(mockApplication));

        bytes memory transferCall = abi.encodeCall(paymentToken.transfer, (creator, raisedAmount));

        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(paymentToken), 0, transferCall));

        uint256 gasBefore = gasleft();
        vm.expectEmit(true, true, false, true);
        emit Transfer(address(mockApplication), creator, raisedAmount);
        mockApplication.executeOutput(voucher);
        uint256 gasUsed = gasBefore - gasleft();

        emit log_named_uint("Gas used (voucher withdraw):", gasUsed);

        assertEq(paymentToken.balanceOf(creator), initialCreatorBalance + raisedAmount);
        assertEq(paymentToken.balanceOf(address(mockApplication)), initialAppBalance - raisedAmount);
    }

    function testGasCloseIssuanceAndWithdraw() public {
        uint256 raisedAmount = 100000 * 10 ** 18;

        string memory closePayload =
            string(abi.encodePacked('{"path":"issuance/close","data":{"creator":"', _addressToString(creator), '"}}'));

        uint256 gasInput = gasleft();
        inputBox.addInput(address(mockApplication), bytes(closePayload));
        uint256 gasUsedInput = gasInput - gasleft();

        string memory withdrawPayload = string(
            abi.encodePacked(
                '{"path":"user/withdraw","data":{"token":"',
                _addressToString(address(paymentToken)),
                '","amount":"100000"}}'
            )
        );

        uint256 gasWithdrawInput = gasleft();
        inputBox.addInput(address(mockApplication), bytes(withdrawPayload));
        uint256 gasUsedWithdrawInput = gasWithdrawInput - gasleft();

        bytes memory transferCall = abi.encodeCall(paymentToken.transfer, (creator, raisedAmount));

        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(paymentToken), 0, transferCall));

        uint256 gasVoucher = gasleft();
        mockApplication.executeOutput(voucher);
        uint256 gasUsedVoucher = gasVoucher - gasleft();

        emit log_string("=== COMPLETE FINALIZATION FLOW ===");
        emit log_named_uint("Gas used (input close):", gasUsedInput);
        emit log_named_uint("Gas used (input withdraw):", gasUsedWithdrawInput);
        emit log_named_uint("Gas used (voucher transfer):", gasUsedVoucher);
        emit log_named_uint("Total gas (L1 on-chain):", gasUsedInput + gasUsedWithdrawInput + gasUsedVoucher);

        assertEq(inputBox.getNumberOfInputs(address(mockApplication)), 2);
        assertEq(paymentToken.balanceOf(creator), raisedAmount);
    }

    function _addressToString(address _addr) internal pure returns (string memory) {
        bytes memory alphabet = "0123456789abcdef";
        bytes memory data = abi.encodePacked(_addr);
        bytes memory str = new bytes(42);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < 20; i++) {
            str[2 + i * 2] = alphabet[uint8(data[i] >> 4)];
            str[3 + i * 2] = alphabet[uint8(data[i] & 0x0f)];
        }
        return string(str);
    }
}
