package api

import (
	"encoding/json"
	"log"
	"net/http"
	"nft-backend/store"
)

func StartServer() {
	http.HandleFunc("/listings", func(w http.ResponseWriter, r *http.Request) {
		list := []store.Listing{}
		for _, listing := range store.Listings {
			list = append(list, listing)
		}

		json.NewEncoder(w).Encode(list)
	})

	log.Println("🚀 API Server running on http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
	log.Println("🛑 API Server stopped")
}
