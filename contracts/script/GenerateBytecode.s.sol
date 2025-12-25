// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import "forge-std-1.9.7/src/Script.sol";
import {Badge} from "../src/token/ERC1155/Badge.sol";

contract GenerateBytecode is Script {
    function run() external {
        bytes memory bytecode = type(Badge).creationCode;
        string memory bytecodeHex = vm.toString(bytecode);
        string memory json = string.concat('{"bytecode":"', bytecodeHex, '"}');

        vm.writeJson(json, "../assets/artifacts/Badge.json");
    }
}
