import {network} from "hardhat";
import {ContractTransactionResponse} from "ethers";

async function main(): Promise<void> {
    const NFT_ADDRESS = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0";
    const MARKET_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";
    const TOKEN_ID = 1;

    const {ethers} = await network.connect();
    const nft = await ethers.getContractAt("NFT", NFT_ADDRESS);

    // Approve the marketplace to transfer the NFT
    const approveTx: ContractTransactionResponse = await nft.approve(MARKET_ADDRESS, TOKEN_ID);
    await approveTx.wait();
    console.log(`Approved marketplace to transfer NFT with token ID ${TOKEN_ID}`);
}

main().catch(console.error);