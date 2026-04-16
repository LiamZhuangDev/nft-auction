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

    mapping(uint256 => Listing) public listings; // listingId => Listing
    mapping(address => mapping(uint256 => bool)) public activeListings; // nftContract => tokenId => isActive
    uint256 public listingCount;

    IAuctionHouse public auctionHouse;
    address public feeRecipient;
    uint256 public constant ListingfeePercent = 250; // 2.5% fee = 250 / 10,000

    constructor(address _auctionContract, address _feeRecipient) {
        auctionHouse = IAuctionHouse(_auctionContract);
        feeRecipient = _feeRecipient;
    }

    function listNFT(address nftContract, uint256 tokenId) 
        external returns (uint256) 
    {
        require(tokenId > 0, "Invalid token ID");
        require(nftContract != address(0), "Invalid NFT contract");
        require(!activeListings[nftContract][tokenId], "NFT is already listed");

        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Only the owner can list the NFT");
        nft.approve(address(this), tokenId); // Approve marketplace to transfer the NFT

        listingCount++;
        listings[listingCount] = Listing({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            active: true
        });
    
        activeListings[nftContract][tokenId] = true;

        return listingCount; // Return the listing ID
    }

    function delistNFT(uint256 listingId) external {
        Listing storage l = listings[listingId];
        require(l.active, "Listing is not active");
        require(l.seller == msg.sender, "Only the seller can delist the NFT");

        l.active = false;
        activeListings[l.nftContract][l.tokenId] = false;
    }

    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 duration
    ) external returns (uint256 auctionId) 
    {
        // Validate inputs
        require(startPrice > 0, "Start price must be greater than zero");
        require(tokenId > 0, "Invalid token ID");
        require(duration > 0, "Duration must be greater than zero");
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