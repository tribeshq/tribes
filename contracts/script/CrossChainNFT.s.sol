// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Helper} from "./Helper.sol";
import {NFT} from "../src/token/ERC721/NFT.sol";
import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SourceMinter} from "../src/chainlink/SourceMinter.sol";
import {DestinationMinter} from "../src/chainlink/DestinationMinter.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";

contract CrossChainNFT is Script, Helper {
    function run(SupportedNetworks destination, SupportedNetworks source) external {
        // Deploy destination minter on ethereumSepolia
        vm.createSelectFork("ethereumSepolia");
        vm.startBroadcast();
        (address destinationRouter,,,) = getConfigFromNetwork(destination);

        NFT nft = new NFT("NFT", "NFT");
        console.log("NFT deployed on ", networks[destination], "with address: ", address(nft));

        DestinationMinter destinationMinter = new DestinationMinter(destinationRouter, address(nft));
        console.log(
            "DestinationMinter deployed on ", networks[destination], "with address: ", address(destinationMinter)
        );

        address minter = nft.owner();
        console.log("Minter role granted to: ", minter);
        vm.stopBroadcast();

        // Deploy source minter on arbitrumSepolia
        vm.createSelectFork("arbitrumSepolia");
        vm.startBroadcast();
        (address sourceRouter,,,) = getConfigFromNetwork(source);

        SourceMinter sourceMinter = new SourceMinter(sourceRouter);
        console.log("SourceMinter deployed on ", networks[source], "with address: ", address(sourceMinter));
        vm.stopBroadcast();
    }
}

contract SetupApplication is Script, Helper {
    function run() external {
        // Transfer ownership of NFT to application
        vm.createSelectFork("ethereumSepolia");
        vm.startBroadcast();
        NFT(vm.envAddress("ETHEREUM_SEPOLIA_NFT")).transferOwnership(vm.envAddress("APPLICATION"));
        address sepoliaMinter = NFT(vm.envAddress("ETHEREUM_SEPOLIA_NFT")).owner();
        console.log("Minter role granted to: ", sepoliaMinter);
        vm.stopBroadcast();

        // Transfer ownership of NFT to application
        vm.createSelectFork("arbitrumSepolia");
        vm.startBroadcast();
        NFT(vm.envAddress("ARBITRUM_SEPOLIA_NFT")).transferOwnership(vm.envAddress("APPLICATION"));
        address arbitrumSepoliaMinter = NFT(vm.envAddress("ARBITRUM_SEPOLIA_NFT")).owner();
        console.log("Minter role granted to: ", arbitrumSepoliaMinter);
        vm.stopBroadcast();
    }
}
