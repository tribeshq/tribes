// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {WebProofXProver} from "../src/vlayer/WebProofXProver.sol";
import {WebProofXVerifier} from "../src/vlayer/WebProofXVerifier.sol";

contract DeployVLayer is Script {
    WebProofXProver public prover;
    WebProofXVerifier public verifier;

    function run() external {
        console.log("Starting VLayer deployment on chain ID:", block.chainid);

        vm.startBroadcast();
        console.log("Deploying WebProofX Prover...");
        prover = new WebProofXProver();
        console.log("Prover deployed to:", address(prover));

        console.log("Deploying WebProofX Verifier...");
        verifier = new WebProofXVerifier(address(prover), 0xc70074BDD26d8cF983Ca6A5b89b8db52D5850051);
        console.log("Verifier deployed to:", address(verifier));
        vm.stopBroadcast();

        _saveDeploymentInfo();

        console.log("VLayer deployment completed!");
    }

    function _saveDeploymentInfo() internal {
        string memory deploymentInfo = string.concat(
            '{"vlayer":{',
            '"chainId":',
            vm.toString(block.chainid),
            ",",
            '"timestamp":',
            vm.toString(block.timestamp),
            ",",
            '"contracts":{',
            '"prover":"',
            vm.toString(address(prover)),
            '",',
            '"verifier":"',
            vm.toString(address(verifier)),
            '",',
            '"inputBoxAddress":"0xc70074BDD26d8cF983Ca6A5b89b8db52D5850051"',
            "}",
            "}}"
        );

        vm.writeJson(deploymentInfo, "./deployments/vlayer.json");
    }
}
