// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol"; // npm install @openzeppelin/contracts
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
    
contract NFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private _tokenCount;

    uint256 public constant MAX_SUPPLY = 10000;
    uint public mintPrice = 0.01 ether;

    event Minted(address indexed minter, uint256 indexed tokenId, string uri);

    constructor() ERC721("GopherNFT", "GONFT") Ownable(msg.sender) {}

    function setMintPrice(uint256 price) public onlyOwner {
        require(price > 0, "Price must be greater than zero");
        mintPrice = price;
    }

    // Create a new NFT and assign ownership to msg.sender
    function mint(string memory uri) public payable returns (uint256) {
        require(_tokenCount < MAX_SUPPLY, "Max supply reached");
        require(msg.value >= mintPrice, "Insufficient payment");

        _tokenCount++;
        _safeMint(msg.sender, _tokenCount);
        _setTokenURI(_tokenCount, uri);

        emit Minted(msg.sender, _tokenCount, uri);
        
        return _tokenCount;
    }

    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        (bool success, ) = owner().call{value: balance}("");
        require(success, "Withdrawal failed");
    }

    // The following functions are overrides required by Solidity.
    // Both ERC721 and ERC721URIStorage implement tokenURI. We need to specify which one to use.
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }

    // Both ERC721 and ERC721URIStorage implement supportsInterface. We need to specify which one to use.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}