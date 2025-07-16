// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {LibOutputSafeCall} from "../library/LibOutputSafeCall.sol";

contract OutputSafeCall {
    using LibOutputSafeCall for address;

    function safeCall(address target, bytes memory data) external returns (bytes memory) {
        return target.safeCall(data);
    }
}