package main

import (
	"log"
	"nft-backend/api"
	"nft-backend/listener"
)

func main() {
	log.Println("Starting app...")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Listener crashed:", r)
			}
		}()
		listener.StartListener()
	}()

	api.StartServer()
}
