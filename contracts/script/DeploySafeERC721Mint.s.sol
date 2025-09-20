// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SafeERC721Mint} from "../src/delegatecall/SafeERC721Mint.sol";

contract DeploySafeERC721Mint is Script {
    SafeERC721Mint public safeERC721Mint;

    function run() external {
        console.log("Starting safe ERC721 mint contracts deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        console.log("Deploying Safe ERC721 Mint...");
        safeERC721Mint = new SafeERC721Mint{salt: keccak256("1596")}();
        console.log("Safe ERC721 Mint deployed to:", address(safeERC721Mint));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("SafeERC721Mint deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"safeERC721Mint":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"safeERC721Mint":"',
            vm.toString(address(safeERC721Mint)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/safeERC721Mint.json");
    }
}
