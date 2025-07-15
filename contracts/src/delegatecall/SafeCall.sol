// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {LibSafeCall} from "../library/LibSafeCall.sol";

contract SafeCall {
    using LibSafeCall for address;

    function safeCall(address target, bytes memory data) external returns (bytes memory) {
        return target.safeCall(data);
    }
}