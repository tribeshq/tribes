// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {IERC20} from "openzeppelin-contracts/token/ERC20/IERC20.sol";

contract SafeEmergencyWithdraw {
    function safeEmergencyERC20Withdraw(IERC20 token, address to) public {
        token.transfer(to, token.balanceOf(address(this)));
    }

    function safeEmergencyETHWithdraw(address to) public {
        uint256 balance = address(this).balance;
        require(balance > 0, "No ETH to withdraw");
        (bool success,) = to.call{value: balance}("");
        require(success, "ETH transfer failed");
    }
}
