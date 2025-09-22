// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {Badge} from "../src/token/ERC721/Badge.sol";
import {MockApplication} from "./mock/MockApplication.sol";
import {BadgeFactory} from "../src/token/ERC721/BadgeFactory.sol";
import {Outputs} from "cartesi-rollups-contracts-2.0.0/src/common/Outputs.sol";
import {SafeERC721Mint, IERC721Mintable} from "../src/delegatecall/SafeERC721Mint.sol";

contract BadgeTest is Test {
    BadgeFactory public badgeFactory;
    SafeERC721Mint public safeERC721Mint;
    MockApplication public mockApplication;

    address public user;

    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event BadgeDeployed(address indexed nft, bytes32 salt);

    function setUp() public {
        user = makeAddr("user");

        badgeFactory = new BadgeFactory();
        safeERC721Mint = new SafeERC721Mint();
        mockApplication = new MockApplication();
    }

    function test_DeterministicDeploymentOfNFTViaFactoryThroughVoucherExecution() public {
        string memory symbol = "MTK";
        string memory name = "MyToken";
        bytes32 salt = keccak256("test-salt");

        bytes memory encodedDeployTx = abi.encodeCall(BadgeFactory.newBadge, (address(mockApplication), salt, name, symbol));
        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(badgeFactory), 0, encodedDeployTx));

        address predictedAddress = address(
            uint160(
                uint256(
                    keccak256(
                        abi.encodePacked(
                            bytes1(0xff),
                            address(badgeFactory),
                            salt,
                            keccak256(
                                abi.encodePacked(
                                    type(Badge).creationCode, abi.encode(address(mockApplication), name, symbol)
                                )
                            )
                        )
                    )
                )
            )
        );

        vm.expectEmit(true, true, false, true);
        emit BadgeDeployed(predictedAddress, salt);
        mockApplication.executeOutput(voucher);

        assertEq(Badge(predictedAddress).name(), name);
        assertEq(Badge(predictedAddress).symbol(), symbol);
    }

    function test_MintNFTThroughDelegatecallVoucher() public {
        string memory symbol = "MTK";
        string memory name = "MyToken";
        bytes32 salt = keccak256("test-salt");

        bytes memory encodedDeployTx = abi.encodeCall(BadgeFactory.newBadge, (address(mockApplication), salt, name, symbol));
        bytes memory voucher = abi.encodeCall(Outputs.Voucher, (address(badgeFactory), 0, encodedDeployTx));

        address predictedAddress = address(
            uint160(
                uint256(
                    keccak256(
                        abi.encodePacked(
                            bytes1(0xff),
                            address(badgeFactory),
                            salt,
                            keccak256(
                                abi.encodePacked(
                                    type(Badge).creationCode, abi.encode(address(mockApplication), name, symbol)
                                )
                            )
                        )
                    )
                )
            )
        );

        vm.expectEmit(true, true, false, true);
        emit BadgeDeployed(predictedAddress, salt);
        mockApplication.executeOutput(voucher);

        bytes memory encodedMintTx =
            abi.encodeCall(SafeERC721Mint.safeMint, (IERC721Mintable(predictedAddress), user, "ipfs://test-uri"));
        bytes memory delegateCallVoucher =
            abi.encodeCall(Outputs.DelegateCallVoucher, (address(safeERC721Mint), encodedMintTx));

        vm.expectEmit(true, true, false, true);
        emit Transfer(address(0), user, 0);
        mockApplication.executeOutput(delegateCallVoucher);
        assertEq(Badge(predictedAddress).ownerOf(0), user);
        assertEq(Badge(predictedAddress).tokenURI(0), "ipfs://test-uri");
        assertEq(Badge(predictedAddress).name(), name);
        assertEq(Badge(predictedAddress).symbol(), symbol);
    }
}