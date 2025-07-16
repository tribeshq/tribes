// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {OutputSafeCall} from "../src/delegatecall/OutputSafeCall.sol";

contract DeployOutputSafeCall is Script {
    OutputSafeCall public outputSafeCall;

    function run() external {
        console.log("Starting SafeCall deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        console.log("Deploying OutputSafeCall contract...");
        outputSafeCall = new OutputSafeCall();
        console.log("OutputSafeCall deployed to:", address(outputSafeCall));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("OutputSafeCall deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"outputSafeCall":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"outputSafeCall":"',
            vm.toString(address(outputSafeCall)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/outputSafeCall.json");
    }
}
