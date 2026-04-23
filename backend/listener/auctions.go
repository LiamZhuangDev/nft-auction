package listener

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"nft-backend/models"
	"nft-backend/repository"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WatchAuctionEvents(client *ethclient.Client, auctionRepo *repository.AuctionRepo, bidRepo *repository.BidRepo) {
	// Load the contract ABI
	file, err := os.ReadFile("abi/auction.json")
	if err != nil {
		log.Fatal("Failed to read ABI file:", err)
	}

	var abiJson struct {
		ABI json.RawMessage `json:"abi"`
	}

	err = json.Unmarshal(file, &abiJson)
	if err != nil {
		log.Fatal("Failed to unmarshal ABI JSON:", err)
	}

	abi, err := abi.JSON(strings.NewReader(string(abiJson.ABI)))
	if err != nil {
		log.Fatal("Failed to parse ABI:", err)
	}

	contractAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// backfill
	pastLogs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for _, vLog := range pastLogs {
		handleAuctionLog(abi, vLog, auctionRepo, bidRepo)
	}

	// subscribe
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal("Failed to subscribe to filter logs:", err)
	}

	log.Println("Listening for Auction events...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case vLog := <-logs:
			handleAuctionLog(abi, vLog, auctionRepo, bidRepo)
		}
	}
}

func handleAuctionLog(abi abi.ABI, vLog types.Log, auctionRepo *repository.AuctionRepo, bidRepo *repository.BidRepo) {
	// Decode the log using the ABI
	event, err := abi.EventByID(vLog.Topics[0])
	if err != nil {
		log.Println("Failed to get event by ID:", err)
		return
	}

	switch event.Name {
	case "AuctionCreated":
		log.Println("receive AuctionCreated event")

		var data struct {
			TokenId    *big.Int
			StartPrice *big.Int
			EndTime    *big.Int
		}

		// Unpack the non-indexed event data into the struct
		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v", err)
			return
		}

		// Extract indexed parameters from the topics
		auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		seller := common.HexToAddress(vLog.Topics[2].Hex())
		nftContract := common.HexToAddress(vLog.Topics[3].Hex())

		// store.Auctions[auctionId.Uint64()] = store.Auction{
		// 	ID:          auctionId.Uint64(),
		// 	Seller:      seller.Hex(),
		// 	NftContract: nftContract.Hex(),
		// 	TokenId:     data.TokenId.String(),
		// 	StartPrice:  data.StartPrice.String(),
		// 	EndTime:     data.EndTime.Uint64(),
		// 	Active:      true,
		// }
		err = auctionRepo.CreateAuction(&models.Auction{
			Seller:      seller.Hex(),
			NftContract: nftContract.Hex(),
			TokenId:     data.TokenId.String(),
			StartPrice:  data.StartPrice.String(),
			EndTime:     data.EndTime.Uint64(),
			Active:      true,
		})
		if err != nil {
			log.Printf("Failed to create auction: %v", err)
			return
		}

		log.Printf("AuctionCreated: ID=%d, Seller=%s, NFT=%s, TokenId=%s, StartPrice=%s, EndTime=%d",
			auctionId.Uint64(), seller.Hex(), nftContract.Hex(), data.TokenId.String(), data.StartPrice.String(), data.EndTime.Uint64())
	case "BidPlaced":
		log.Println("receive BidPlaced event")

		var data struct {
			Amount *big.Int
		}

		// Unpack the non-indexed event data into the struct
		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v", err)
			return
		}
		// Extract indexed parameters from the topics
		auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		bidder := common.HexToAddress(vLog.Topics[2].Hex())

		// store.Bids[auctionId.Uint64()] = store.Bid{
		// 	AuctionId: auctionId.Uint64(),
		// 	Bidder:    bidder.Hex(),
		// 	Amount:    data.Amount.String(),
		// 	Timestamp: uint64(time.Now().Unix()),
		// }
		err = bidRepo.CreateBid(&models.Bid{
			AuctionID: auctionId.Uint64(),
			Bidder:    bidder.Hex(),
			Amount:    data.Amount.String(),
			Timestamp: uint64(time.Now().Unix()),
		})
		if err != nil {
			log.Printf("Failed to create bid, err: %v", err)
			return
		}

		log.Printf("BidPlaced: AuctionId=%d, Bidder=%s, Amount=%s", auctionId.Uint64(), bidder.Hex(), data.Amount.String())
	case "AuctionEnded":
		log.Println("receive AuctionEnded event")

		var data struct {
			Amount *big.Int
		}

		// Unpack the non-indexed event data into the struct
		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v", err)
			return
		}

		// Extract indexed parameters from the topics
		auctionId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		winner := common.HexToAddress(vLog.Topics[2].Hex())

		// if auction, exists := store.Auctions[auctionId.Uint64()]; exists {
		// 	auction.Active = false
		// 	store.Auctions[auctionId.Uint64()] = auction
		// }

		err = auctionRepo.UpdateAuctionStatus(auctionId.Uint64(), false)
		if err != nil {
			log.Printf("Failed to update auction %d status, error: %v", auctionId.Int64(), err)
			return
		}

		log.Printf("AuctionEnded: AuctionId=%d, Winner=%s, Amount=%s", auctionId.Uint64(), winner.Hex(), data.Amount.String())
	default:
		log.Printf("Unknown event: %s", vLog.Topics[0].Hex())
	}
}
