// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.27;

import {IERC20} from "openzeppelin-contracts/token/ERC20/IERC20.sol";

contract EmergencyWithdraw {
    error ZeroBalance();
    error ETHTransferFailed();
    error NotAdmin(address admin);

    modifier onlyAdmin(address admin) {
        if (msg.sender != admin) {
            revert NotAdmin(admin);
        }
        _;
    }

    function emergencyERC20Withdraw(address to, address admin, IERC20 token) public onlyAdmin(admin) {
        token.transfer(to, token.balanceOf(address(this)));
    }

    function emergencyETHWithdraw(address to, address admin) public onlyAdmin(admin) {
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
