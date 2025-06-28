// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {CCIPReceiver} from "@chainlink/contracts-ccip/contracts/applications/CCIPReceiver.sol";
import {Client} from "@chainlink/contracts-ccip/contracts/libraries/Client.sol";
import {IERC721} from "openzeppelin-contracts/token/ERC721/IERC721.sol";

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
contract DestinationMinter is CCIPReceiver {
    IERC721 nft;

    event MintCallSuccessfull();

    constructor(address nftAddress, address router) CCIPReceiver(router) {
        nft = IERC721(nftAddress);
    }

    function _ccipReceive(Client.Any2EVMMessage memory message) internal override {
        (bool success,) = address(nft).call(message.data);
        require(success);
        emit MintCallSuccessfull();
    }
}
