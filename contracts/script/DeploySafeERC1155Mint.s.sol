// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SafeERC1155Mint} from "../src/delegatecall/SafeERC1155Mint.sol";

contract DeploySafeERC1155Mint is Script {
    SafeERC1155Mint public safeERC1155Mint;

    function run() external {
        console.log("Starting SafeERC1155Mint deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        console.log("Deploying SafeERC1155Mint...");
        safeERC1155Mint = new SafeERC1155Mint();
        console.log("SafeERC1155Mint deployed to:", address(safeERC1155Mint));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("SafeERC1155Mint deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"safeERC1155Mint":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"safeERC1155Mint":"',
            vm.toString(address(safeERC1155Mint)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/safeERC1155Mint.json");
    }
}
