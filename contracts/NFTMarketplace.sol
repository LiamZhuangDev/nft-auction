// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/interfaces/IERC721.sol";

interface IAuctionHouse {
    function createAuction(
        address seller,
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 duration
    ) external returns (uint256 auctionId);

    function isActive(address nftContract, uint256 tokenId) external view returns (bool);
}

// Design Overview:
// User
//  ↓
// Marketplace (entry point / orchestrator)
//  └─ Delegates auctions → AuctionHouse
//                           ├─ createAuction
//                           ├─ placeBid
//                           └─ endAuction
contract NFTMarketplace {
    struct Listing {
        address seller;
        address nftContract;
        uint256 tokenId;
        bool active;
    }

    Listing[] public listings;
    mapping(address => mapping(uint256 => bool)) public activeListings; // nftContract => tokenId => isActive

    IAuctionHouse public auctionHouse;

    constructor(address _auctionContract) {
        auctionHouse = IAuctionHouse(_auctionContract);
    }

    function listNFT(address nftContract, uint256 tokenId) 
        external returns (uint256) 
    {
        require(tokenId > 0, "Invalid token ID");
        require(nftContract != address(0), "Invalid NFT contract");
        require(!activeListings[nftContract][tokenId], "NFT is already listed");

        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Only the owner can list the NFT");

        listings.push(Listing({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            active: true
        }));
    
        activeListings[nftContract][tokenId] = true;

        return listings.length - 1; // Return the listing ID
    }

    function delistNFT(uint256 listingId) external {
        require(listingId < listings.length, "Invalid listing ID");

        Listing storage l = listings[listingId];
        require(l.active, "Listing is not active");
        require(l.seller == msg.sender, "Only the seller can delist the NFT");

        l.active = false;
        activeListings[l.nftContract][l.tokenId] = false;
    }

    function createAuction(
        uint256 listingId,
        uint256 startPrice,
        uint256 duration
    ) external returns (uint256 auctionId) 
    {
        // Validate inputs
        require(listingId < listings.length, "Invalid listing ID");
        require(startPrice > 0, "Start price must be greater than zero");
        require(duration > 0, "Duration must be greater than zero");

        Listing storage l = listings[listingId];
        address nftContract = l.nftContract;
        uint256 tokenId = l.tokenId;

        // Validate listing
        require(l.active, "Listing is not active");
        require(tokenId > 0, "Invalid token ID");
        require(nftContract != address(0), "Invalid NFT contract");
        require(!auctionHouse.isActive(nftContract, tokenId), "NFT is already in auction");

        // Verify ownership and approval
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Only the owner can create an auction");
        require(nft.getApproved(tokenId) == address(this) || nft.isApprovedForAll(msg.sender, address(this)), "Marketplace must be approved to transfer the NFT");
        
        // Create the auction
        nft.safeTransferFrom(msg.sender, address(auctionHouse), tokenId); // Transfer NFT to auction house
        auctionId = auctionHouse.createAuction(msg.sender, nftContract, tokenId, startPrice, duration);
    }
}