import { network } from "hardhat";
import { ContractTransactionResponse } from "ethers";

async function main(): Promise<void> {
    const AUCTION_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    const NFT_ADDRESS = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0";
    const AUCTION_ID = 0;
    const TOKEN_ID = 1;

    const { ethers } = await network.connect();
    const [seller, bidder1, bidder2, bidder3] = await ethers.getSigners();
    const highestBidder = bidder3; // Assuming bidder3 is the highest bidder from previous bids

    console.log("Seller:", seller.address);
    console.log("Bidder:", highestBidder.address);

    const auction = await ethers.getContractAt(
        "NFTAuctionHouse",
        AUCTION_ADDRESS,
        seller
    );

    const nft = await ethers.getContractAt("NFT", NFT_ADDRESS);

    // Check auction state
    const a = await auction.auctions(AUCTION_ID);
    console.log("\n=== BEFORE FINALIZE ===");
    console.log("Highest Bidder:", a.highestBidder);
    console.log("Highest Bid:", ethers.formatEther(a.highestBid));
    console.log("Active:", a.active);

    // Balances before
    const sellerBefore = await ethers.provider.getBalance(seller.address);
    const bidderBefore = await ethers.provider.getBalance(highestBidder.address);

    console.log("\nBalances BEFORE:");
    console.log("Seller:", ethers.formatEther(sellerBefore));
    console.log("Bidder:", ethers.formatEther(bidderBefore));

    // Finalize auction
    const tx: ContractTransactionResponse =
        await auction.finalizeAuction(AUCTION_ID);
    const receipt = await tx.wait();

    console.log("\nTx hash:", receipt?.hash);

    // Balances after
    const sellerAfter = await ethers.provider.getBalance(seller.address);
    const bidderAfter = await ethers.provider.getBalance(highestBidder.address);

    console.log("\n=== AFTER FINALIZE ===");

    console.log("Seller balance:", ethers.formatEther(sellerAfter));
    console.log("Bidder balance:", ethers.formatEther(bidderAfter));

    console.log(
        "Seller earned:",
        ethers.formatEther(sellerAfter - sellerBefore)
    );

    console.log(
        "Bidder spent:",
        ethers.formatEther(bidderBefore - bidderAfter)
    );

    // Check NFT ownership
    const newOwner = await nft.ownerOf(TOKEN_ID);

    console.log("\nNFT Owner:", newOwner);

    console.log(
        "\n NFT transferred to bidder:",
        newOwner.toLowerCase() === highestBidder.address.toLowerCase()
    );

    console.log("\n Auction finalized and verified!");
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});