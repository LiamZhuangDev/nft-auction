import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";
import { getContracts } from "./contracts.js";

export async function placeBid() {
  const { auction } = getContracts();

  const auctionId = document.getElementById("auctionId").value;
  const bidAmount = document.getElementById("bidAmount").value;

  const tx = await auction.placeBid(auctionId, {
    value: ethers.parseEther(bidAmount)
  });

  await tx.wait();
  console.log("Bid Placed");
}

export async function finalizeAuction() {
  const { auction } = getContracts();

  const auctionId = document.getElementById("finalizeAuctionId").value;

  const tx = await auction.finalizeAuction(auctionId);
  await tx.wait();

  console.log("Auction Finalized");
}