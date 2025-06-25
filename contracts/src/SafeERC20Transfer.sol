// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

pragma solidity ^0.8.20;

import {IERC20} from "openzeppelin-contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "openzeppelin-contracts/token/ERC20/utils/SafeERC20.sol";

contract SafeERC20Transfer {
    using SafeERC20 for IERC20;

    error NotTarget(address target);

    function safeTransferTargeted(IERC20 token, address target, address to, uint256 value) external {
        if (msg.sender != target) {
            revert NotTarget(target);
        }
        token.safeTransfer(to, value);
    }
}