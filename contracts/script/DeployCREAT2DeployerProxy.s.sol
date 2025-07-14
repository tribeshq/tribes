// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {CREATE2Deployer} from "../src/proxy/CREATE2Deployer.sol";

contract DeployCREAT2DeployerProxy is Script {
    CREATE2Deployer public deployer;

    function run() external {
        vm.startBroadcast();
        deployer = new CREATE2Deployer();
        console.log("CREAT2Deployer deployed with address: ", address(deployer));
        vm.stopBroadcast();
    }
}
