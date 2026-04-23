package main

import (
	"log"
	"nft-backend/api"
	"nft-backend/config"
	"nft-backend/listener"
	"nft-backend/repository"
)

func main() {
	log.Println("Starting app...")

	// Initialize database
	db := config.ConnectDB()

	// Create data repositories
	listingRepo := repository.NewListingRepo(db)
	auctionRepo := repository.NewAuctionRepo(db)
	bidRepo := repository.NewBidRepo(db)

	log.Println("Database connected")

	// Start blockchain listener (background)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Listener crashed:", r)
			}
		}()
		listener.StartListener(listingRepo, auctionRepo, bidRepo)
	}()

	// Start API server (blocking)
	api.StartServer(listingRepo, auctionRepo, bidRepo)
}
