import { network } from "hardhat";
import { ContractTransactionResponse } from "ethers";

async function main(): Promise<void> {
    const MARKET_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
    const LISTING_ID = 0;
    const {ethers} = await network.connect();
    const START_PRICE = ethers.parseEther("0.1");
    const DURATION = 60; // use 60 seconds for testing

    const [seller] = await ethers.getSigners();

    console.log("Seller:", seller.address);

    const market = await ethers.getContractAt(
        "NFTMarketplace",
        MARKET_ADDRESS,
        seller
    );

    const tx: ContractTransactionResponse = await market.createAuction(
        LISTING_ID,
        START_PRICE,
        DURATION
    );

    await tx.wait();

    console.log("Auction created!");
}

main().catch(console.error);