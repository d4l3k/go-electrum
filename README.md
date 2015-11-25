# go-electrum [![GoDoc](https://godoc.org/github.com/d4l3k/go-electrum?status.svg)](https://godoc.org/github.com/d4l3k/go-electrum)
A pure Go (golang) electrum bitcoin library. This is very much WIP and has a number of unimplemented methods. This will eventually be rewritten into a more Go-esque library and handle wallet generation.

```bash
go get -u "github.com/d4l3k/go-electrum"
```

Example
```go
package main

import (
  "log"
  "github.com/d4l3k/go-electrum"
)

func main() {
	node := electrum.NewNode()
	if err := node.ConnectTCP("electrum.dragonzone.net:50001"); err != nil {
		log.Fatal(err)
	}
	balance, err := node.BlockchainAddressGetBalance("1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Address balance: %+v", balance)
}
```
See [examples/](https://github.com/d4l3k/go-electrum/tree/master/example) for more.

# License
go-electrum is licensed under the MIT license.

Made by [Tristan Rice](https://fn.lc).
