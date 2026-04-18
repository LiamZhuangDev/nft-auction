import { network } from "hardhat";

async function main() {
  const { ethers } = await network.connect();
  const [deployer] = await ethers.getSigners();

  console.log("Deploying with:", deployer.address);

  // 1. AuctionHouse
  const auctionFactory = await ethers.getContractFactory("NFTAuctionHouse");
  const auction = await auctionFactory.deploy(deployer.address);
  await auction.waitForDeployment();

  console.log("AuctionHouse:", await auction.getAddress());

  // 2. Marketplace
  const marketFactory = await ethers.getContractFactory("NFTMarketplace");
  const market = await marketFactory.deploy(await auction.getAddress());
  await market.waitForDeployment();

  console.log("Marketplace:", await market.getAddress());

  // 3. NFT
  const nftFactory = await ethers.getContractFactory("NFT");
  const nft = await nftFactory.deploy();
  await nft.waitForDeployment();

  console.log("NFT:", await nft.getAddress());
}

main().catch(console.error);