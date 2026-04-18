import { network } from "hardhat";
import { ContractTransactionResponse } from "ethers";

async function main(): Promise<void> {
  const nftAddress: string = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0";

  const { ethers } = await network.connect();
  const nft = await ethers.getContractAt("NFT", nftAddress);

  const tx: ContractTransactionResponse = await nft.mint(
    "ipfs://fake-uri",
    {
      value: ethers.parseEther("0.01"),
    }
  );

  await tx.wait();

  console.log("Minted!");
}

main().catch((error: unknown) => {
  console.error(error);
  process.exitCode = 1;
});