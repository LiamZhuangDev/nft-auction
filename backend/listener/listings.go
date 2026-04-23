package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"nft-backend/models"
	"nft-backend/repository"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WatchListingEvents(client *ethclient.Client, listingRepo *repository.ListingRepo) {
	// Load the contract ABI
	file, err := os.ReadFile("abi/marketplace.json")
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

	contractAddress := common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")

	// Create a filter query specifically for the Listing events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal("Failed to subscribe to filter logs:", err)
	}

	fmt.Println("Listening for Listing events...")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal("Subscription error:", err)
		case vLog := <-logs:
			handleListingLog(abi, vLog, listingRepo)
		}
	}
}

func handleListingLog(abi abi.ABI, vLog types.Log, listingRepo *repository.ListingRepo) {
	event, err := abi.EventByID(vLog.Topics[0])
	if err != nil {
		log.Printf("Unknown event: %s", vLog.Topics[0].Hex())
		return
	}

	switch event.Name {
	case "Listed":
		log.Println("Handling Listed event...")

		var data struct {
			TokenId   *big.Int
			ListingId *big.Int
		}

		// Unpack the non-indexed event data into the struct
		err := abi.UnpackIntoInterface(&data, event.Name, vLog.Data)
		if err != nil {
			log.Printf("Failed to unpack log: %v", err)
			return
		}

		// Extract indexed parameters from the topics
		seller := common.HexToAddress(vLog.Topics[1].Hex())
		nftContract := common.HexToAddress(vLog.Topics[2].Hex())

		// store.Listings[data.ListingId.Uint64()] = store.Listing{
		// 	ID:          data.ListingId.Uint64(),
		// 	Seller:      seller.Hex(),
		// 	NftContract: nftContract.Hex(),
		// 	TokenId:     data.TokenId.String(),
		// 	Active:      true,
		// }
		err = listingRepo.CreateListing(&models.Listing{
			Seller:      seller.Hex(),
			NftContract: nftContract.Hex(),
			TokenId:     data.TokenId.String(),
			Active:      true,
		})
		if err != nil {
			log.Printf("Failed to list NFT: %v", err)
		}

		log.Printf("New listing: Seller=%s, NFT Contract=%s, Token ID=%s, Listing ID=%s",
			seller.Hex(), nftContract.Hex(), data.TokenId.String(), data.ListingId.String())
	default:
		log.Printf("Unhandled event: %s", event.Name)
	}
}
