// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {NFT} from "../token/ERC721/NFT.sol";
import {Client} from "@chainlink/contracts-ccip/contracts/libraries/Client.sol";
import {Ownable} from "openzeppelin-contracts/access/Ownable.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
contract SourceMinter is Ownable {
    NFT public immutable nft;
    address public immutable sourceRouter;

    event MessageSent(bytes32 messageId);

    error InsufficientFee(uint256 currentValue, uint256 requiredFee);

    constructor(NFT _nft, address _sourceRouter) Ownable(msg.sender) {
        nft = _nft;
        sourceRouter = _sourceRouter;
    }

    receive() external payable {}

    function mint(uint64 destinationChainSelector, address to, address minter) external payable onlyOwner {
        if (destinationChainSelector == 16015286601757825753) {
            nft.mint(to);
            return;
        }

        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(minter),
            data: abi.encodeWithSignature("mint(address)", to),
            tokenAmounts: new Client.EVMTokenAmount[](0),
            extraArgs: "",
            feeToken: address(0)
        });
        uint256 fee = IRouterClient(sourceRouter).getFee(destinationChainSelector, message);
        
        if (msg.value < fee) {
            revert InsufficientFee(msg.value, fee);
        }

        bytes32 messageId = IRouterClient(sourceRouter).ccipSend{value: fee}(destinationChainSelector, message);

        emit MessageSent(messageId);
    }
}
