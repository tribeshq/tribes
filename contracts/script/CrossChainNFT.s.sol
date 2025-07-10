// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {stdJson} from "forge-std-1.9.7/src/StdJson.sol";
import {NFT} from "../src/token/ERC721/NFT.sol";
import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {SourceMinter} from "../src/chainlink/SourceMinter.sol";
import {DestinationMinter} from "../src/chainlink/DestinationMinter.sol";

contract CrossChainNFTSourceMinter is Script {
    function run() external {
        vm.startBroadcast();
        NFT sourceNFT = new NFT("NFT", "NFT");
        console.log("Source NFT deployed on with address: ", address(sourceNFT));

        // Chain selector for ethereum sepolia: 16015286601757825753
        address sourceRouter = 0x2a9C5afB0d0e4BAb2BCdaE109EC4b0c4Be15a165;

        SourceMinter sourceMinter = new SourceMinter(sourceNFT, sourceRouter);
        console.log("SourceMinter deployed on with address: ", address(sourceMinter));

        sourceNFT.transferOwnership(address(sourceMinter));
        console.log("Source NFT ownership transferred to SourceMinter");

        address sourceNFTOwner = sourceNFT.owner();
        console.log("Owner of Source NFT: ", sourceNFTOwner);
        vm.stopBroadcast();
    }
}

contract CrossChainNFTDestinationMinter is Script {
    DestinationMinter destinationMinter;

    function run() external {
        vm.startBroadcast();
        NFT destinationNft = new NFT("NFT", "NFT");
        console.log("NFT deployed on with address: ", address(destinationNft));

        // Chain selector for arbitrum sepolia: 3478487238524512106
        address destinationRouter = 0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59;

        destinationMinter = new DestinationMinter(destinationNft, destinationRouter);
        console.log("DestinationMinter deployed on with address: ", address(destinationMinter));

        destinationNft.transferOwnership(address(destinationMinter));
        console.log("Destination NFT ownership transferred to DestinationMinter");

        address destinationNFTOwner = destinationNft.owner();
        console.log("Owner of Destination NFT: ", destinationNFTOwner);
        vm.stopBroadcast();
    }
}

contract SetupApplication is Script {
    SourceMinter sourceMinter;

    function run() external {
        string memory root = vm.projectRoot();
        string memory path = string.concat(root, "/broadcast/CrossChainNFT.s.sol/11155111/run-latest.json");
        string memory json = vm.readFile(path);

        address sourceMinterAddress = bytesToAddress(stdJson.parseRaw(json, ".transactions[0].contractAddress"));

        vm.startBroadcast();
        address applicationAddress = vm.parseAddress(vm.prompt("Enter application address"));

        sourceMinter = SourceMinter(payable(sourceMinterAddress));
        sourceMinter.transferOwnership(applicationAddress);
        address owner = sourceMinter.owner();
        console.log("Minter role granted to: ", owner);
        vm.stopBroadcast();
    }

    function bytesToAddress(bytes memory bys) private pure returns (address addr) {
        assembly {
            addr := mload(add(bys, 32))
        }
    }
}
