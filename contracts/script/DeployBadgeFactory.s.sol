// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {BadgeFactory} from "../src/token/ERC1155/BadgeFactory.sol";

contract DeployBadgeFactory is Script {
    BadgeFactory public badgeFactory;

    function run() external {
        console.log("Starting Deployer deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        badgeFactory = new BadgeFactory();
        console.log("BadgeFactory deployed to:", address(badgeFactory));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("BadgeFactory deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"badgeFactory":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"badgeFactory":"',
            vm.toString(address(badgeFactory)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/badgeFactory.json");
    }
}
