package api

import (
	"encoding/json"
	"log"
	"net/http"
	"nft-backend/store"
)

func StartServer() {
	http.HandleFunc("/listings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		list := []store.Listing{}
		for _, listing := range store.Listings {
			list = append(list, listing)
		}

		json.NewEncoder(w).Encode(list)
	})

	http.HandleFunc("/auctions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println("Received request for /auctions")

		list := []store.Auction{}
		for _, auction := range store.Auctions {
			list = append(list, auction)
		}

		log.Printf("Sending response for /auctions: %d auctions", len(list))

		json.NewEncoder(w).Encode(list)
	})

	log.Println("API Server running on http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
	log.Println("API Server stopped")
}
