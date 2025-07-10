// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {NFT} from "../token/ERC721/NFT.sol";
import {Client} from "@chainlink/contracts-ccip/contracts/libraries/Client.sol";
import {CCIPReceiver} from "@chainlink/contracts-ccip/contracts/applications/CCIPReceiver.sol";

interface INFT {
    function mint(address to) external;
}

contract DestinationMinter is CCIPReceiver {
    NFT public immutable nft;

    error InvalidSelector(bytes4 selector);
    error InvalidPayloadLength(uint256 length);
    error MintCallFailed(address to, bytes32 messageId);

    event MintCallSuccessful(address indexed to, bytes32 indexed messageId);

    constructor(NFT _nft, address router) CCIPReceiver(router) {
        nft = _nft;
    }

    function _ccipReceive(Client.Any2EVMMessage memory message) internal override {
        bytes memory payload = message.data;

        if (payload.length != 4 + 32) {
            revert InvalidPayloadLength(payload.length);
        }

        bytes4 selector =
            (bytes4(payload[0]) | (bytes4(payload[1]) >> 8) | (bytes4(payload[2]) >> 16) | (bytes4(payload[3]) >> 24));

        if (selector != INFT.mint.selector) {
            revert InvalidSelector(selector);
        }

        bytes memory data = new bytes(payload.length - 4);
        for (uint256 i = 0; i < data.length; i++) {
            data[i] = payload[i + 4];
        }
        address to = abi.decode(data, (address));

        try nft.mint(to) {
            emit MintCallSuccessful(to, message.messageId);
        } catch {
            revert MintCallFailed(to, message.messageId);
        }
    }
}
