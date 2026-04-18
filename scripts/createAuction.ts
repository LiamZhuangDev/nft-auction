import { network } from "hardhat";
import { ContractTransactionResponse } from "ethers";

async function main(): Promise<void> {
    const TOKEN_ID = 1;
    const LISTING_ID = 0;
    const {ethers} = await network.connect();
    const NFT_ADDRESS="0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"
    const nft = await ethers.getContractAt("NFT", NFT_ADDRESS);
    console.log("Owner:", await nft.ownerOf(TOKEN_ID));
    console.log("Approved:", await nft.getApproved(TOKEN_ID));

    const MARKET_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
    const market = await ethers.getContractAt("NFTMarketplace", MARKET_ADDRESS);
    console.log("Listings:", await market.listings(LISTING_ID));

    const AUCTION_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    const auction = await ethers.getContractAt("NFTAuctionHouse", AUCTION_ADDRESS);
    console.log("isAuctionActive:", await auction.isActive(NFT_ADDRESS, TOKEN_ID));

    await market.createAuction.staticCall(
    LISTING_ID,
    ethers.parseEther("0.1"),
    60
    );

//   const MARKET_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
//   const LISTING_ID = 0;
//   const {ethers} = await network.connect();
//   const START_PRICE = ethers.parseEther("0.1");
//   const DURATION = 60; // use 60 seconds for testing

//   const [seller] = await ethers.getSigners();

//   console.log("Seller:", seller.address);

//   const market = await ethers.getContractAt(
//     "NFTMarketplace",
//     MARKET_ADDRESS,
//     seller
//   );

//   const tx: ContractTransactionResponse = await market.createAuction(
//     LISTING_ID,
//     START_PRICE,
//     DURATION
//   );

//   await tx.wait();

//   console.log("Auction created!");
}

main().catch(console.error);