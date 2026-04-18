import { network } from "hardhat";
import { ContractTransactionResponse } from "ethers";

async function main(): Promise<void> {
    const AUCTION_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    const AUCTION_ID = 0;

    const {ethers} = await network.connect();
    const [, bidder] = await ethers.getSigners(); // second account

    console.log("Bidder:", bidder.address);

    const auction = await ethers.getContractAt(
        "NFTAuctionHouse",
        AUCTION_ADDRESS,
        bidder
    );

    const before = await ethers.provider.getBalance(bidder.address);

    const tx: ContractTransactionResponse = await auction.placeBid(
        AUCTION_ID,
        {
        value: ethers.parseEther("0.2"),
        }
    );

    await tx.wait();

    console.log("Bid placed!");

    const after = await ethers.provider.getBalance(bidder.address);

    console.log(
        "Bidder spent:",
        ethers.formatEther(before - after)
    );
}

main().catch(console.error);