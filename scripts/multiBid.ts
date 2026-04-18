import { network } from "hardhat";
import { ContractTransactionResponse } from "ethers";

async function main(): Promise<void> {
    const AUCTION_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    const AUCTION_ID = 0;

    const { ethers } = await network.connect();
    const [seller, bidder1, bidder2, bidder3] = await ethers.getSigners();

    const auction = await ethers.getContractAt(
        "NFTAuctionHouse",
        AUCTION_ADDRESS,
        seller
    );

    console.log("Bidders:");
    console.log("Bidder1:", bidder1.address);
    console.log("Bidder2:", bidder2.address);
    console.log("Bidder3:", bidder3.address);

    // Helper function
    async function placeBid(
        bidder: any,
        amount: string
    ): Promise<void> {
        const before = await ethers.provider.getBalance(bidder.address);

        const tx: ContractTransactionResponse =
        await auction.connect(bidder).placeBid(AUCTION_ID, {
            value: ethers.parseEther(amount),
        });

        await tx.wait();

        const after = await ethers.provider.getBalance(bidder.address);

        console.log(
            `\n${bidder.address} bid ${amount} ETH`
            );
        console.log(
            "Spent:",
            ethers.formatEther(before - after)
        );
    }

    // Sequence of bids
    await placeBid(bidder1, "0.3"); // first bid, must be above starting price
    await placeBid(bidder2, "0.4"); // outbids bidder1
    await placeBid(bidder3, "0.5"); // outbids bidder2

    // Check pending returns
    console.log("\n=== Pending Returns ===");

    console.log(
        "Bidder1:",
        ethers.formatEther(
            await auction.pendingReturns(AUCTION_ID, bidder1.address)
        )
    );

    console.log(
        "Bidder2:",
        ethers.formatEther(
            await auction.pendingReturns(AUCTION_ID, bidder2.address)
        )
    );

    console.log(
        "Bidder3:",
        ethers.formatEther(
            await auction.pendingReturns(AUCTION_ID, bidder3.address)
        )
    );
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});