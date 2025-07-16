// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.27;

import {LibError} from "cartesi-rollups-contracts-2.0.0/src/library/LibError.sol";
import {IERC20} from "openzeppelin-contracts/token/ERC20/IERC20.sol";

contract EmergencyWithdraw {
    using LibError for bytes;

    error ZeroBalance();

    function emergencyERC20Withdraw(address to, IERC20 token) public {
        token.transfer(to, token.balanceOf(address(this)));
    }

    function emergencyETHWithdraw(address to) public {
        uint256 balance = address(this).balance;
        if (balance == 0) revert ZeroBalance();

        (bool success, bytes memory returndata) = to.call{value: balance}("");
        if (!success) returndata.raise();
    }
}
