// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import {Script} from "forge-std-1.9.7/src/Script.sol";

struct PropertyArgs {
    address router;
}

contract SetupProperty is Script {
    PropertyArgs public propertyArgs;

    mapping(uint256 => PropertyArgs) public chainIdToNetworkConfig;

    constructor() {
        chainIdToNetworkConfig[11155111] = getEthereumSepoliaPropertyArgs();
        chainIdToNetworkConfig[1440002] = getArbitrumSepoliaPropertyArgs();
        propertyArgs = chainIdToNetworkConfig[block.chainid];
    }

    function getEthereumSepoliaPropertyArgs()
        internal
        pure
        returns (PropertyArgs memory EthereumSepoliaPropertyArgs)
    {
        EthereumSepoliaPropertyArgs = PropertyArgs({
            router: 0x0000000000000000000000000000000000000000
        });
    }

    function getArbitrumSepoliaPropertyArgs()
        internal
        pure
        returns (PropertyArgs memory ArbitrumSepoliaPropertyArgs)
    {
        ArbitrumSepoliaPropertyArgs = PropertyArgs({
            router: 0x0000000000000000000000000000000000000000
        });
    }
}