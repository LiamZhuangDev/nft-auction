import { ethers } from "https://cdn.jsdelivr.net/npm/ethers@6.10.0/+esm";
import { getContracts } from "./contracts.js";
import { NFT_ADDRESS } from "./config.js";

export async function listNFT() {
    const { marketplace } = getContracts();

    const tokenId = document.getElementById("listTokenId").value;

    const tx = await marketplace.listNFT(NFT_ADDRESS, tokenId);
    await tx.wait();

    console.log("NFT listed for sale");
}

export async function createAuction() {
    const { marketplace } = getContracts();

    const listingId = document.getElementById("listingId").value;
    const startPrice = document.getElementById("startPrice").value;
    const duration = document.getElementById("duration").value;

    const tx = await marketplace.createAuction(
        listingId,
        ethers.parseEther(startPrice),
        duration
    );
    await tx.wait();

    console.log("Auction created");
}