import { connectWallet } from "./wallet.js";
import { mintNFT, approveNFT } from "./nft.js";
import { listNFT, createAuction } from "./marketplace.js";
import { placeBid, finalizeAuction } from "./auction.js";

// Expose functions to HTML
window.connectWallet = connectWallet;
window.mintNFT = mintNFT;
window.approveNFT = approveNFT;
window.listNFT = listNFT;
window.createAuction = createAuction;
window.placeBid = placeBid;
window.finalizeAuction = finalizeAuction;