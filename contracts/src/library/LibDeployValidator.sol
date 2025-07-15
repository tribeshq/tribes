// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

library LibDeployValidator {
    error ContractNotDeployed(address target);
    error ContractAlreadyExists(address target);

    function _isContract(address target) internal view returns (bool) {
        return target.code.length > 0;
    }

    function _computeCreate2Address(address factory, bytes32 salt, bytes memory bytecode)
        internal
        pure
        returns (address)
    {
        return address(uint160(uint256(keccak256(abi.encodePacked(bytes1(0xff), factory, salt, keccak256(bytecode))))));
    }

    function checkIfExists(address factory, bytes32 salt, bytes memory bytecode)
        internal
        view
        returns (bool, address)
    {
        address contractAddress = _computeCreate2Address(factory, salt, bytecode);
        return (_isContract(contractAddress), contractAddress);
    }

    function validateCreate2Deployment(address factory, bytes32 salt, bytes memory bytecode)
        internal
        view
        returns (address)
    {
        address predicted = _computeCreate2Address(factory, salt, bytecode);
        if (!_isContract(predicted)) revert ContractNotDeployed(predicted);
        return predicted;
    }
}
