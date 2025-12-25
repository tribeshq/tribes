// SPDX-License-Identifier: MIT

pragma solidity ^0.8.27;

interface IERC1155Mintable {
    function mint(address account, uint256 id, uint256 amount, bytes memory data) external;
    function mintBatch(address to, uint256[] memory ids, uint256[] memory amounts, bytes memory data) external;
}

contract SafeERC1155Mint {
    error NotAContract(address target);

    function mint(IERC1155Mintable nft, address to, uint256 id, uint256 amount, bytes memory data) public {
        if (address(nft).code.length == 0) {
            revert NotAContract(address(nft));
        }

        nft.mint(to, id, amount, data);
    }

    function mintBatch(
        IERC1155Mintable nft,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) public {
        if (address(nft).code.length == 0) {
            revert NotAContract(address(nft));
        }

        nft.mintBatch(to, ids, amounts, data);
    }
}
