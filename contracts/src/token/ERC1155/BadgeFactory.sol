// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Badge} from "./Badge.sol";

contract BadgeFactory {
    event BadgeDeployed(address indexed nft, bytes32 salt);
    event BadgeAlreadyDeployed(address indexed nft, bytes32 salt);

    function newBadge(address initialOwner, bytes32 salt) external returns (Badge) {
        address predicted = computeAddress(initialOwner, salt);

        if (predicted.code.length > 0) {
            emit BadgeAlreadyDeployed(predicted, salt);
            return Badge(predicted);
        }

        Badge badge = new Badge{salt: salt}(initialOwner);

        emit BadgeDeployed(address(badge), salt);
        return badge;
    }

    function computeAddress(address initialOwner, bytes32 salt) public view returns (address) {
        return address(
            uint160(
                uint256(
                    keccak256(
                        abi.encodePacked(
                            bytes1(0xff),
                            address(this),
                            salt,
                            keccak256(abi.encodePacked(type(Badge).creationCode, abi.encode(initialOwner)))
                        )
                    )
                )
            )
        );
    }
}
