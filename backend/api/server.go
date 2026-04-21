package api

import (
	"log"
	"net/http"
	"nft-backend/store"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/listings", func(c *gin.Context) {
		log.Printf("len of listed NFT: %d", len(store.Listings))

		list := []store.Listing{}
		for _, listing := range store.Listings {
			list = append(list, listing)
		}
		c.JSON(http.StatusOK, list)
	})

	r.GET("/auctions", func(c *gin.Context) {
		// print the count of auctions
		log.Printf("len of auctions: %d", len(store.Auctions))

		list := []store.Auction{}
		for _, auction := range store.Auctions {
			list = append(list, auction)
		}
		c.JSON(http.StatusOK, list)
	})

	r.GET("/bids/:id", func(c *gin.Context) {
		id := c.Param("id")
		auctionId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction ID"})
			return
		}

		result := []store.Bid{}
		for _, bid := range store.Bids {
			if bid.AuctionId == auctionId {
				result = append(result, bid)
			}
		}
		c.JSON(http.StatusOK, result)
	})

	log.Println("API Server running on http://localhost:8081")
	r.Run(":8081")
	log.Println("API Server stopped")
}
