package api

import (
	"log"
	"net/http"
	"nft-backend/repository"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer(
	listingRepo *repository.ListingRepo,
	auctionRepo *repository.AuctionRepo,
	bidRepo *repository.BidRepo) {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/listings", func(c *gin.Context) {
		listings, err := listingRepo.GetAllListings()
		if err != nil {
			log.Printf("failed to fetch listings: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch listings",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": listings,
		})
	})

	r.GET("/auctions", func(c *gin.Context) {
		auctions, err := auctionRepo.GetAllAuctions()
		if err != nil {
			log.Printf("failed to fetch auctions: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch auctions",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": auctions,
		})
	})

	r.GET("/bids/:id", func(c *gin.Context) {
		id := c.Param("id")
		auctionId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
			return
		}

		bids, err := bidRepo.GetBids(auctionId)
		if err != nil {
			log.Printf("failed to fetch auction %d bids: %v", auctionId, err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch bids",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": bids,
		})
	})

	log.Println("API Server running on http://localhost:8081")
	r.Run(":8081")
	log.Println("API Server stopped")
}
