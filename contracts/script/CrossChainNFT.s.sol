// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Helper} from "./Helper.sol";
import {NFT} from "../src/token/ERC721/NFT.sol";
import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SourceMinter} from "../src/chainlink/SourceMinter.sol";
import {DestinationMinter} from "../src/chainlink/DestinationMinter.sol";

contract CrossChainNFT is Script, Helper {
    function run(SupportedNetworks destination, SupportedNetworks source) external {
        // Deploy source minter on arbitrumSepolia
        vm.createSelectFork("arbitrumSepolia");
        vm.startBroadcast();
        (address sourceRouter,,,) = getConfigFromNetwork(source);

        NFT sourceNFT = new NFT("NFT", "NFT");
        console.log("Source NFT deployed on ", networks[source], "with address: ", address(sourceNFT));

        SourceMinter sourceMinter = new SourceMinter(sourceNFT, sourceRouter);
        console.log("SourceMinter deployed on ", networks[source], "with address: ", address(sourceMinter));

        sourceNFT.transferOwnership(address(sourceMinter));
        console.log("Source NFT ownership transferred to SourceMinter");

        address sourceNFTOwner = sourceNFT.owner();
        console.log("Owner of Source NFT: ", sourceNFTOwner);
        vm.stopBroadcast();

        // Deploy destination minter on ethereumSepolia
        vm.createSelectFork("ethereumSepolia");
        vm.startBroadcast();
        (address destinationRouter,,,) = getConfigFromNetwork(destination);

        NFT destinationNft = new NFT("NFT", "NFT");
        console.log("NFT deployed on ", networks[destination], "with address: ", address(destinationNft));

        DestinationMinter destinationMinter = new DestinationMinter(destinationNft, destinationRouter);
        console.log(
            "DestinationMinter deployed on ", networks[destination], "with address: ", address(destinationMinter)
        );

        destinationNft.transferOwnership(address(destinationMinter));
        console.log("Destination NFT ownership transferred to DestinationMinter");

        address destinationNFTOwner = destinationNft.owner();
        console.log("Owner of Destination NFT: ", destinationNFTOwner);
        vm.stopBroadcast();
    }
}

contract SetupApplication is Script, Helper {
    function run() external {
        // Transfer ownership of NFT to application
        vm.createSelectFork("arbitrumSepolia");
        vm.startBroadcast();
        SourceMinter(payable(vm.envAddress("ARBITRUM_SEPOLIA_SOURCE_MINTER"))).transferOwnership(vm.envAddress("APPLICATION"));
        address arbitrumSepoliaMinter = SourceMinter(payable(vm.envAddress("ARBITRUM_SEPOLIA_SOURCE_MINTER"))).owner();
        console.log("Minter role granted to: ", arbitrumSepoliaMinter);
        vm.stopBroadcast();
    }
}
