// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SafeCall} from "../src/delegatecall/SafeCall.sol";

contract DeploySafeCall is Script {
    SafeCall public safeCall;

    function run() external {
        console.log("Starting SafeCall deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        console.log("Deploying SafeCall contract...");
        safeCall = new SafeCall();
        console.log("SafeCall deployed to:", address(safeCall));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("SafeCall deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"safeCall":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"safeCall":"',
            vm.toString(address(safeCall)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/safeCall.json");
    }
}
