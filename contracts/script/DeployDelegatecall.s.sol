// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import {Script} from "forge-std-1.9.7/src/Script.sol";
import {console} from "forge-std-1.9.7/src/console.sol";
import {EmergencyWithdraw} from "../src/delegatecall/EmergencyWithdraw.sol";

contract DeployDelegatecall is Script {
    EmergencyWithdraw emergencyWithdraw;

    function run() external {
        emergencyWithdraw = new EmergencyWithdraw();

        console.log("EmergencyWithdraw deployed to: ", address(emergencyWithdraw));
    }
}