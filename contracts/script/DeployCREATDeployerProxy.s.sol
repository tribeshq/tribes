// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {CREATEDeployer} from "../src/proxy/CREATEDeployer.sol";

contract DeployCREATDeployerProxy is Script {
    CREATEDeployer public deployer;

    function run() external {
        vm.startBroadcast();
        deployer = new CREATEDeployer();
        console.log("CREATDeployer deployed with address: ", address(deployer));
        vm.stopBroadcast();
    }
}
