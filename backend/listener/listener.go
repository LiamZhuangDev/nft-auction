package listener

import (
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

func StartListener() {
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
		WatchListingEvents(client)
	}()

	go func() {
		defer wg.Done()
		WatchAuctionEvents(client)
	}()

	wg.Wait()

	log.Println("Listener STOPPED")
}
