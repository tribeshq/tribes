// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {NFT} from "../src/token/ERC721/NFT.sol";
import {Test} from "forge-std-1.9.7/src/Test.sol";
import {SourceMinter} from "../src/chainlink/SourceMinter.sol";
import {DestinationMinter} from "../src/chainlink/DestinationMinter.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";
import {CCIPLocalSimulator, IRouterClient} from "@chainlink/local/src/ccip/CCIPLocalSimulator.sol";

contract CrossChainNFT is Test {
    NFT public nft;
    SourceMinter public sourceMinter;
    DestinationMinter public destinationMinter;
    CCIPLocalSimulator public ccipLocalSimulator;

    address public guest;
    uint64 public destinationChainSelector;

    IRouterClient public sourceRouter;
    IRouterClient public destinationRouter;

    function setUp() public {
        ccipLocalSimulator = new CCIPLocalSimulator();

        (uint64 chainSelector, IRouterClient _sourceRouter, IRouterClient _destinationRouter,,,,) =
            ccipLocalSimulator.configuration();

        sourceRouter = _sourceRouter;
        destinationRouter = _destinationRouter;

        nft = new NFT(address(this), "NFT", "NFT");
        destinationMinter = new DestinationMinter(address(nft), address(destinationRouter));
        nft.transferOwnership(address(destinationMinter));

        sourceMinter = new SourceMinter(address(sourceRouter));

        guest = makeAddr("guest");
        destinationChainSelector = chainSelector;
    }

    function test_executeReceivedMessageAsFunctionCall() external {
        vm.startPrank(guest);
        sourceMinter.mint(destinationChainSelector, address(destinationMinter));
        vm.stopPrank();

        assertEq(nft.balanceOf(guest), 1);
    }
}
