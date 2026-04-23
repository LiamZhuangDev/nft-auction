package listener

import (
	"log"
	"nft-backend/repository"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

func StartListener(
	listingRepo *repository.ListingRepo,
	auctionRepo *repository.AuctionRepo,
	bidRepo *repository.BidRepo) {
	log.Println("Listener STARTED")

	client, err := ethclient.Dial("ws://127.0.0.1:8545/")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	log.Println("Subscribing to logs...")

	// Use a WaitGroup to keep the main goroutine alive while the listeners run in parallel
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		WatchListingEvents(client, listingRepo)
	}()

	go func() {
		defer wg.Done()
		WatchAuctionEvents(client, auctionRepo, bidRepo)
	}()

	wg.Wait()

	log.Println("Listener STOPPED")
}
