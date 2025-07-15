// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {Token} from "../src/token/ERC20/Token.sol";

contract DeployTokens is Script {
    Token public collateral;
    Token public stablecoin;

    function run() external {
        console.log("Starting tokens deployment on chain ID:", block.chainid);

        vm.startBroadcast();

        console.log("Deploying Collateral Token...");
        collateral = new Token("Collateral", "COLL");
        console.log("Collateral deployed to:", address(collateral));

        console.log("Deploying Stablecoin Token...");
        stablecoin = new Token("Stablecoin", "STBL");
        console.log("Stablecoin deployed to:", address(stablecoin));

        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("Tokens deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"tokens":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"collateral":"',
            vm.toString(address(collateral)),
            '",',
            '"stablecoin":"',
            vm.toString(address(stablecoin)),
            '"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/tokens.json");
    }
}
