pragma solidity ^0.8.27;

import {Test} from "forge-std-1.9.7/src/Test.sol";
import {Badge} from "../src/token/ERC1155/Badge.sol";
import {CREATEDeployer} from "../src/proxy/CREATEDeployer.sol";
import {CREATE2Deployer} from "../src/proxy/CREATE2Deployer.sol";

contract BadgeTest is Test {
    Badge public badge;
    CREATEDeployer public deployer;
    CREATE2Deployer public deployer2;

    address public user;
    address public unauthorized;
    address public applicationAddress;

    event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
    event ContractDeployed(address indexed contractAddress, bytes32 indexed salt);

    function setUp() public {
        user = makeAddr("user");
        unauthorized = makeAddr("unauthorized");
        applicationAddress = makeAddr("applicationAddress");

        deployer = new CREATEDeployer();
        deployer2 = new CREATE2Deployer();
    }

    function test_DeployBadgeViaProxy() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        address deployedAddress = deployer.deploy(badgeBytecode);
        badge = Badge(deployedAddress);
        assertEq(badge.owner(), applicationAddress);
    }

    function test_OnlyApplicationCanMint() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        address deployedAddress = deployer.deploy(badgeBytecode);
        badge = Badge(deployedAddress);

        uint256 id = 1;
        uint256 amount = 1;
        bytes memory data = "";

        vm.prank(applicationAddress);
        vm.expectEmit(true, true, true, true);
        emit TransferSingle(applicationAddress, address(0), user, id, amount);
        badge.mint(user, id, amount, data);

        assertEq(badge.balanceOf(user, id), amount);
    }

    function test_UnauthorizedCannotMint() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        address deployedAddress = deployer.deploy(badgeBytecode);
        badge = Badge(deployedAddress);

        uint256 id = 1;
        uint256 amount = 1;
        bytes memory data = "";

        vm.prank(unauthorized);
        vm.expectRevert();
        badge.mint(user, id, amount, data);
    }

    function test_RandomAddressCannotMint() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        address deployedAddress = deployer.deploy(badgeBytecode);
        badge = Badge(deployedAddress);

        uint256 id = 1;
        uint256 amount = 1;
        bytes memory data = "";

        vm.prank(address(0x123));
        vm.expectRevert();
        badge.mint(user, id, amount, data);
    }

    // CREATE2 Tests
    function test_DeployBadgeViaCreate2() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        bytes32 salt = keccak256("test-salt");
        address predictedAddress = _computeCreate2Address(address(deployer2), salt, badgeBytecode);

        vm.prank(applicationAddress);
        vm.expectEmit(true, true, false, true);
        emit ContractDeployed(predictedAddress, salt);

        address deployedAddress = deployer2.deploy2(badgeBytecode, salt);
        assertEq(deployedAddress, predictedAddress);

        badge = Badge(deployedAddress);
        assertEq(badge.owner(), applicationAddress);
    }

    function test_OnlyApplicationCanMintCreate2() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        bytes32 salt = keccak256("mint-test-salt");
        address deployedAddress = deployer2.deploy2(badgeBytecode, salt);
        badge = Badge(deployedAddress);

        uint256 id = 1;
        uint256 amount = 1;
        bytes memory data = "";

        vm.prank(applicationAddress);
        vm.expectEmit(true, true, true, true);
        emit TransferSingle(applicationAddress, address(0), user, id, amount);
        badge.mint(user, id, amount, data);

        assertEq(badge.balanceOf(user, id), amount);
    }

    function test_UnauthorizedCannotMintCreate2() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        bytes32 salt = keccak256("unauthorized-test-salt");
        address deployedAddress = deployer2.deploy2(badgeBytecode, salt);
        badge = Badge(deployedAddress);

        uint256 id = 1;
        uint256 amount = 1;
        bytes memory data = "";

        vm.prank(unauthorized);
        vm.expectRevert();
        badge.mint(user, id, amount, data);
    }

    function test_DeployerCannotMintCreate2() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        bytes32 salt = keccak256("deployer-test-salt");
        address deployedAddress = deployer2.deploy2(badgeBytecode, salt);
        badge = Badge(deployedAddress);

        uint256 id = 1;
        uint256 amount = 1;
        bytes memory data = "";

        vm.prank(address(0x123));
        vm.expectRevert();
        badge.mint(user, id, amount, data);
    }

    function test_SameSaltSameAddress() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        bytes32 salt = keccak256("same-salt");
        address predictedAddress = _computeCreate2Address(address(deployer2), salt, badgeBytecode);

        address deployedAddress1 = deployer2.deploy2(badgeBytecode, salt);

        vm.expectRevert();
        deployer2.deploy2(badgeBytecode, salt);

        assertEq(deployedAddress1, predictedAddress);
    }

    function test_DifferentSaltDifferentAddress() public {
        bytes memory badgeBytecode = abi.encodePacked(vm.getCode("Badge.sol:Badge"), abi.encode(applicationAddress));

        bytes32 salt1 = keccak256("salt-1");
        bytes32 salt2 = keccak256("salt-2");

        address predictedAddress1 = _computeCreate2Address(address(deployer2), salt1, badgeBytecode);
        address predictedAddress2 = _computeCreate2Address(address(deployer2), salt2, badgeBytecode);

        assertTrue(predictedAddress1 != predictedAddress2);

        address deployedAddress1 = deployer2.deploy2(badgeBytecode, salt1);
        address deployedAddress2 = deployer2.deploy2(badgeBytecode, salt2);

        assertEq(deployedAddress1, predictedAddress1);
        assertEq(deployedAddress2, predictedAddress2);
        assertTrue(deployedAddress1 != deployedAddress2);
    }

    function _computeCreate2Address(address factory, bytes32 salt, bytes memory bytecode)
        internal
        pure
        returns (address)
    {
        bytes32 hash = keccak256(abi.encodePacked(bytes1(0xff), factory, salt, keccak256(bytecode)));
        return address(uint160(uint256(hash)));
    }
}
