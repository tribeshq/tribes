// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {NFT} from "../src/token/ERC721/NFT.sol";
import {Test} from "forge-std-1.9.7/src/Test.sol";
import {SourceMinter} from "../src/chainlink/SourceMinter.sol";
import {DestinationMinter} from "../src/chainlink/DestinationMinter.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";
import {CCIPLocalSimulator, IRouterClient} from "@chainlink/local/src/ccip/CCIPLocalSimulator.sol";

contract CrossChainNFT is Test {
    NFT public sourceNFT;
    NFT public destinationNft;
    SourceMinter public sourceMinter;
    DestinationMinter public destinationMinter;
    CCIPLocalSimulator public ccipLocalSimulator;

    address public guest;
    address public application;
    uint64 public destinationChainSelector;

    IRouterClient public sourceRouter;
    IRouterClient public destinationRouter;

    function setUp() public {
        ccipLocalSimulator = new CCIPLocalSimulator();

        (uint64 chainSelector, IRouterClient _sourceRouter, IRouterClient _destinationRouter,,,,) =
            ccipLocalSimulator.configuration();

        sourceRouter = _sourceRouter;
        destinationRouter = _destinationRouter;

        sourceNFT = new NFT("NFT", "NFT");
        sourceMinter = new SourceMinter(sourceNFT, address(sourceRouter));
        sourceNFT.transferOwnership(address(sourceMinter));

        application = makeAddr("application");
        sourceMinter.transferOwnership(application);

        destinationNft = new NFT("NFT", "NFT");
        destinationMinter = new DestinationMinter(destinationNft, address(destinationRouter));
        destinationNft.transferOwnership(address(destinationMinter));

        guest = makeAddr("guest");
        destinationChainSelector = chainSelector;
    }

    function test_executeReceivedMessageAsFunctionCall() external {
        vm.startPrank(application);
        sourceMinter.mint(destinationChainSelector, guest, address(destinationMinter));
        vm.stopPrank();
        assertEq(destinationNft.balanceOf(guest), 1);
    }
}
