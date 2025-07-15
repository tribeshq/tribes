// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {Deployer} from "../src/proxy/Deployer.sol";

contract DeployDeployer is Script {
    Deployer public deployer;

    function run() external {
        console.log("Starting Deployer deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        deployer = new Deployer();
        console.log("Deployer deployed to:", address(deployer));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("Deployer deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"deployer":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"deployer":"',
            vm.toString(address(deployer)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/deployer.json");
    }
}
