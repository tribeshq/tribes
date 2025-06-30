// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {LinkTokenInterface} from "@chainlink/contracts/src/v0.8/shared/interfaces/LinkTokenInterface.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";
import {Client} from "@chainlink/contracts-ccip/contracts/libraries/Client.sol";
import {Withdraw} from "./utils/Withdraw.sol";

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
contract SourceMinter is Withdraw {
    address immutable sourceRouter;

    event MessageSent(bytes32 messageId);

    constructor(address _sourceRouter) {
        sourceRouter = _sourceRouter;
    }

    receive() external payable {}

    function mint(uint64 destinationChainSelector, address receiver) external {
        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(receiver),
            data: abi.encodeWithSignature("mint(address)", msg.sender),
            tokenAmounts: new Client.EVMTokenAmount[](0),
            extraArgs: "",
            feeToken: address(0)
        });

        uint256 fee = IRouterClient(sourceRouter).getFee(destinationChainSelector, message);

        bytes32 messageId = IRouterClient(sourceRouter).ccipSend{value: fee}(destinationChainSelector, message);

        emit MessageSent(messageId);
    }
}
