// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {Badge} from "../src/token/ERC1155/Badge.sol";
import {MockApplication} from "./mock/MockApplication.sol";
import {BadgeFactory} from "../src/token/ERC1155/BadgeFactory.sol";
import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";
import {SafeERC1155Mint, IERC1155Mintable} from "../src/delegatecall/SafeERC1155Mint.sol";

contract BadgeTest is Test {
    BadgeFactory public badgeFactory;
    SafeERC1155Mint public safeERC1155Mint;
    MockApplication public mockApplication;

    address public user;

    event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
    event BadgeDeployed(address indexed nft, bytes32 salt);

    function setUp() public {
        user = makeAddr("user");

        badgeFactory = new BadgeFactory();
        safeERC1155Mint = new SafeERC1155Mint();
        mockApplication = new MockApplication();
    }

    function test_DeterministicDeploymentOfNFTViaFactoryThroughVoucherExecution() public {
        bytes32 salt = keccak256("test-salt");

        bytes memory encodedDeployTx = abi.encodeCall(BadgeFactory.newBadge, (address(mockApplication), salt));
        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(badgeFactory), 0, encodedDeployTx));

        address predictedAddress = badgeFactory.computeAddress(address(mockApplication), salt);

        vm.expectEmit(true, true, false, true);
        emit BadgeDeployed(predictedAddress, salt);
        mockApplication.executeOutput(voucher);

        assertEq(Badge(predictedAddress).owner(), address(mockApplication));
    }

    function test_MintNFTThroughDelegatecallVoucher() public {
        bytes32 salt = keccak256("test-salt");

        bytes memory encodedDeployTx = abi.encodeCall(BadgeFactory.newBadge, (address(mockApplication), salt));
        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(badgeFactory), 0, encodedDeployTx));

        address predictedAddress = badgeFactory.computeAddress(address(mockApplication), salt);

        vm.expectEmit(true, true, false, true);
        emit BadgeDeployed(predictedAddress, salt);
        mockApplication.executeOutput(voucher);

        uint256 tokenId = 1;
        uint256 amount = 1;
        bytes memory data = "";

        bytes memory encodedMintTx =
            abi.encodeCall(SafeERC1155Mint.mint, (IERC1155Mintable(predictedAddress), user, tokenId, amount, data));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC1155Mint), encodedMintTx));

        vm.expectEmit(true, true, true, true);
        emit TransferSingle(address(mockApplication), address(0), user, tokenId, amount);
        mockApplication.executeOutput(delegateCallVoucher);
        assertEq(Badge(predictedAddress).balanceOf(user, tokenId), amount);
        assertEq(Badge(predictedAddress).owner(), address(mockApplication));
    }
}
