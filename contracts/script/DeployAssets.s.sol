// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {Token} from "../src/token/ERC20/Token.sol";
import {console} from "forge-std-1.9.7/src/console.sol";

contract DeployAssets is Script {
    Token collateral;
    Token stablecoin;

    function run() external {
        vm.startBroadcast();
        collateral = new Token("Collateral", "COLL");
        stablecoin = new Token("Stablecoin", "STBL");
        vm.stopBroadcast();

        console.log("Collateral deployed to: ", address(collateral));
        console.log("Stablecoin deployed to: ", address(stablecoin));
    }
}
