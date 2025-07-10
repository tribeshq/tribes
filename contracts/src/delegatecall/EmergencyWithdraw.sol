// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.27;

import {IERC20} from "openzeppelin-contracts/token/ERC20/IERC20.sol";
import {AccessControl} from "openzeppelin-contracts/access/AccessControl.sol";

contract EmergencyWithdraw is AccessControl {
    error ZeroBalance();
    error ETHTransferFailed();

    bytes32 public constant APP_ROLE = keccak256("APP_ROLE");

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }

    function emergencyERC20Withdraw(address to, IERC20 token) public onlyRole(APP_ROLE) {
        token.transfer(to, token.balanceOf(address(this)));
    }

    function emergencyETHWithdraw(address to) public onlyRole(APP_ROLE) {
        uint256 balance = address(this).balance;
        if (balance == 0) {
            revert ZeroBalance();
        }
        (bool success,) = to.call{value: balance}("");
        if (!success) {
            revert ETHTransferFailed();
        }
    }
}
