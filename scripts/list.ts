import {network} from "hardhat";

async function main(): Promise<void> {
    const {ethers} = await network.connect();
    const NFT_ADDRESS = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0";
    const MARKET_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
    const TOKEN_ID = 1;

    const market = await ethers.getContractAt("NFTMarketplace", MARKET_ADDRESS);
    const tx = await market.listNFT(NFT_ADDRESS, TOKEN_ID);
    const receipt = await tx.wait();
    console.log(`NFT listed for sale with token ID ${TOKEN_ID}, receipt: ${receipt}`);
}

main().catch(console.error);