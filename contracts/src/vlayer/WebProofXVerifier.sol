// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.21;

import {Proof} from "vlayer-1.2.0/src/Proof.sol";
import {WebProofXProver} from "./WebProofXProver.sol";
import {Verifier} from "vlayer-1.2.0/src/Verifier.sol";

contract WebProofXVerifier is Verifier {
    address public prover;

    constructor(address _prover) {
        prover = _prover;
    }

    function verify(Proof calldata proof, string memory username, address account)
        public
        onlyVerified(prover, WebProofXProver.main.selector)
    {
        // TODO:
        // 1. Send social account input to application with username and account address (msg.sender)
    }
}
