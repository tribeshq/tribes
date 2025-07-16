// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {LibError} from "cartesi-rollups-contracts-2.0.0/src/library/LibError.sol";

library LibOutputSafeCall {
    using LibError for bytes;

    error NotAContract(address target);

    function safeCall(address target, bytes memory data) internal returns (bytes memory) {
        if (target.code.length == 0) revert NotAContract(target);

        (bool success, bytes memory returndata) = target.call(data);
        if (!success) returndata.raise();

        return returndata;
    }
}
