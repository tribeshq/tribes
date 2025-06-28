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
    enum PayFeesIn {
        Native,
        LINK
    }

    address immutable i_link;

    event MessageSent(bytes32 messageId);

    constructor(address link) {
        i_link = link;
    }

    receive() external payable {}

    function mint(address router, uint64 destinationChainSelector, address receiver, PayFeesIn payFeesIn) external {
        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(receiver),
            data: abi.encodeWithSignature("mint(address)", msg.sender),
            tokenAmounts: new Client.EVMTokenAmount[](0),
            extraArgs: "",
            feeToken: payFeesIn == PayFeesIn.LINK ? i_link : address(0)
        });

        uint256 fee = IRouterClient(router).getFee(destinationChainSelector, message);

        bytes32 messageId;

        if (payFeesIn == PayFeesIn.LINK) {
            LinkTokenInterface(i_link).approve(router, fee);
            messageId = IRouterClient(router).ccipSend(destinationChainSelector, message);
        } else {
            messageId = IRouterClient(router).ccipSend{value: fee}(destinationChainSelector, message);
        }

        emit MessageSent(messageId);
    }
}
