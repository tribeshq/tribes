// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {LibDeployValidator} from "../library/LibDeployValidator.sol";

contract Deployer {
    using LibDeployValidator for address;

    event ContractDeployed(address indexed contractAddress, bytes32 indexed salt);
    event ContractAlreadyExists(address indexed contractAddress, bytes32 indexed salt);

    receive() external payable {}

    function deploy(bytes memory _code, bytes32 _salt) external payable returns (address addr) {
        (bool exists, address contractAddress) = LibDeployValidator.checkIfExists(address(this), _salt, _code);
        if (exists) {
            emit ContractAlreadyExists(contractAddress, _salt);
            return contractAddress;
        }

        assembly {
            addr := create2(callvalue(), add(_code, 0x20), mload(_code), _salt)
        }

        LibDeployValidator.validateCreate2Deployment(address(this), _salt, _code);
        emit ContractDeployed(addr, _salt);
    }
}
