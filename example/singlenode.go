package main

import (
	"log"

	"github.com/d4l3k/go-electrum"
)

func main() {
	node, err := electrum.NewNode("electrum.dragonzone.net:50001")
	if err != nil {
		log.Fatal(err)
	}
	version, err := node.ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Version: %s", version)
	banner, err := node.ServerBanner()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Banner: %s", banner)
	address, err := node.ServerDonationAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Banner: %s", address)
}
