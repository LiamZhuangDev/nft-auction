import {network} from "hardhat";
import {ContractTransactionResponse} from "ethers";

async function main(): Promise<void> {
    const AUCTION_ADDRESS = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    const AUCTION_ID = 0;

    const {ethers} = await network.connect();
    const [seller, bidder1, bidder2] = await ethers.getSigners();
    const auction = await ethers.getContractAt(
        "NFTAuctionHouse",
        AUCTION_ADDRESS,
        seller
    );

    async function withdraw(bidder: any): Promise<void> {
        const pending = await auction.pendingReturns(AUCTION_ID, bidder.address);

        console.log(`\n${bidder.address} has ${ethers.formatEther(pending)} ETH pending for withdrawal`);

        if (pending == 0n) {
            console.log("No funds to withdraw");
            return;
        } else {
            const before = await ethers.provider.getBalance(bidder.address);

            const tx: ContractTransactionResponse = await auction.connect(bidder).withdrawBid(AUCTION_ID);
            await tx.wait();

            const after = await ethers.provider.getBalance(bidder.address);

            console.log(
                `Withdrawn: ${ethers.formatEther(after - before)} ETH`
            );
        }
    }

    await withdraw(bidder1);
    await withdraw(bidder2);
}

main().catch(console.error);