package main

import (
	"log"

	"github.com/d4l3k/go-electrum/electrum"
)

func main() {
	node := electrum.NewNode()
	if err := node.ConnectTCP("btc.mustyoshi.com:50001"); err != nil {
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
	log.Printf("Address: %s", address)

	peers, err := node.ServerPeersSubscribe()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Peers: %+v", peers)

	numblocks, err := node.BlockchainNumBlocksSubscribe()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Numblocks: %+v", numblocks)

	headerChan, err := node.BlockchainHeadersSubscribe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for header := range headerChan {
			log.Printf("Headers: %+v", header)
		}
	}()

	hashChan, err := node.BlockchainAddressSubscribe("1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for hash := range hashChan {
			log.Printf("Address history hash: %+v", hash)
		}
	}()

	history, err := node.BlockchainAddressGetHistory("1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Address history: %+v", history)

	transaction, err := node.BlockchainTransactionGet("0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Transaction: %s", transaction)

	transactions, err := node.BlockchainAddressListUnspent("1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unspent transactions: %+v", transactions)

	// TODO(d4l3k) seems to not work, need to subscribe first maybe?
	balance, err := node.BlockchainAddressGetBalance("1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Address balance: %+v", balance)
}
