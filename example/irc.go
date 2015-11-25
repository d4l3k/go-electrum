package main

import (
	"log"

	"github.com/d4l3k/go-electrum/irc"
)

func main() {
	log.Println(irc.FindElectrumServers())
}
