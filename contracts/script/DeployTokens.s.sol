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
        string memory collateralName = vm.prompt("Collateral name");
        string memory collateralSymbol = vm.prompt("Collateral symbol");
        address collateralInitialOwner = vm.parseAddress(vm.prompt("Collateral initial owner"));
        collateral = new Token(collateralName, collateralSymbol, collateralInitialOwner);
        console.log("Collateral deployed to:", address(collateral));

        console.log("Deploying Stablecoin Token...");
        string memory stablecoinName = vm.prompt("Stablecoin name");
        string memory stablecoinSymbol = vm.prompt("Stablecoin symbol");
        address stablecoinInitialOwner = vm.parseAddress(vm.prompt("Stablecoin initial owner"));
        stablecoin = new Token(stablecoinName, stablecoinSymbol, stablecoinInitialOwner);
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
