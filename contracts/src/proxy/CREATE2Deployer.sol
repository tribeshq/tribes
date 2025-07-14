// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

contract CREATE2Deployer {
    event ContractDeployed(address indexed contractAddress, bytes32 indexed salt);

    receive() external payable {}

    function deploy2(bytes memory _code, bytes32 _salt) external payable returns (address addr) {
        assembly {
            // create2(v, p, n, s)
            // v = amount of ETH to send
            // p = pointer in memory to start of code
            // n = size of code
            // s = salt
            addr := create2(callvalue(), add(_code, 0x20), mload(_code), _salt)
        }
        // return address 0 on error
        require(addr != address(0), "deploy failed");

        emit ContractDeployed(addr, _salt);
    }
}
