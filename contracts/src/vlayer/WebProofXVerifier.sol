// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Proof} from "vlayer-1.2.0/src/Proof.sol";
import {WebProofXProver} from "./WebProofXProver.sol";
import {Verifier} from "vlayer-1.2.0/src/Verifier.sol";
import {IInputBox} from "cartesi-rollups-contracts-2.0.0/src/inputs/IInputBox.sol";

contract WebProofXVerifier is Verifier {
    address public prover;
    address public inputBox;

    constructor(address _prover, address _inputBox) {
        prover = _prover;
        inputBox = _inputBox;
    }

    function verify(Proof calldata, string memory username, address account, address application)
        public
        onlyVerified(prover, WebProofXProver.main.selector)
    {
        bytes memory input = abi.encodePacked(
            '{"path":"social/verifier/create","data":{"address":"',
            toString(account),
            '","username":"',
            username,
            '","platform":"twitter"}}'
        );
        IInputBox(inputBox).addInput(application, input);
    }

    function toString(address account) public pure returns (string memory) {
        return toString(abi.encodePacked(account));
    }

    function toString(bytes memory data) public pure returns (string memory) {
        bytes memory alphabet = "0123456789abcdef";

        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        unchecked {
            for (uint256 i = 0; i < data.length; i++) {
                str[2 + i * 2] = alphabet[uint256(uint8(data[i] >> 4))];
                str[3 + i * 2] = alphabet[uint256(uint8(data[i] & 0x0f))];
            }
        }
        return string(str);
    }
}
