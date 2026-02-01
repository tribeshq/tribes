// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SafeEmergencyWithdraw} from "../src/delegatecall/SafeEmergencyWithdraw.sol";

contract DeployEmergency is Script {
    SafeEmergencyWithdraw public emergencyWithdraw;

    function run() external {
        console.log("Starting emergency contracts deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        console.log("Deploying Emergency Withdraw...");
        emergencyWithdraw = new SafeEmergencyWithdraw();
        console.log("Emergency Withdraw deployed to:", address(emergencyWithdraw));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("Emergency deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"emergency":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"emergencyWithdraw":"',
            vm.toString(address(emergencyWithdraw)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/emergency.json");
    }
}
