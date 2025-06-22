// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {ERC20Collateral} from "../src/ERC20Collateral.sol";
import {ERC20Stablecoin} from "../src/ERC20Stablecoin.sol";
import {SafeERC20Transfer} from "../src/SafeERC20Transfer.sol";
import {WebProofXProver} from "../src/vlayer/WebProofXProver.sol";
import {WebProofXVerifier} from "../src/vlayer/WebProofXVerifier.sol";

contract Deploy is Script {
    ERC20Collateral public erc20Collateral;
    ERC20Stablecoin public erc20Stablecoin;
    WebProofXProver public webProofXProver;
    SafeERC20Transfer public safeERC20Transfer;
    WebProofXVerifier public webProofXVerifier;

    function setUp() public {}

    function run() public returns (WebProofXVerifier, WebProofXProver, ERC20Collateral, SafeERC20Transfer) {
        vm.startBroadcast();
        webProofXProver = new WebProofXProver();
        webProofXVerifier = new WebProofXVerifier(address(webProofXProver));

        erc20Collateral = new ERC20Collateral();
        erc20Stablecoin = new ERC20Stablecoin();
        safeERC20Transfer = new SafeERC20Transfer();
        vm.stopBroadcast();
        return (webProofXVerifier, webProofXProver, erc20Collateral, safeERC20Transfer);
    }
}