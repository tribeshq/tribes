// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {WebProofXProver} from "../src/vlayer/WebProofXProver.sol";
import {WebProofXVerifier} from "../src/vlayer/WebProofXVerifier.sol";

contract DeployVlayer is Script {
    WebProofXProver prover;
    WebProofXVerifier verifier;

    function run() external {
        vm.startBroadcast();
        prover = new WebProofXProver();

        address proverAddress = vm.parseAddress(vm.prompt("Enter prover address"));
        address inputBoxAddress = vm.parseAddress(vm.prompt("Enter input box address"));
        verifier = new WebProofXVerifier(proverAddress, inputBoxAddress);
        vm.stopBroadcast();

        console.log("Prover deployed to: ", address(prover));
        console.log("Verifier deployed to: ", address(verifier));
    }
}
