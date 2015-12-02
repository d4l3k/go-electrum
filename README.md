# go-electrum [![GoDoc](https://godoc.org/github.com/d4l3k/go-electrum?status.svg)](https://godoc.org/github.com/d4l3k/go-electrum)
A pure Go [Electrum](https://electrum.org/) bitcoin library. This makes it easy to write bitcoin based services using Go without having to run a full bitcoin node.

![go-electrum](https://rawgit.com/d4l3k/go-electrum/master/media/logo.png)

This is very much WIP and has a number of unimplemented methods. This will eventually be rewritten into a more Go-esque library and handle wallet generation.

Packages provided

* [electrum](https://godoc.org/github.com/d4l3k/go-electrum/electrum) - Library for using JSON-RPC to talk directly to Electrum servers.
* [wallet](https://godoc.org/github.com/d4l3k/go-electrum/wallet) - A bitcoin wallet built on [btcwallet](https://github.com/btcsuite/btcwallet) with Electrum as the backend.
* [irc](https://godoc.org/github.com/d4l3k/go-electrum/irc) - A helper module for finding electrum servers using the [#electrum IRC channel](http://docs.electrum.org/en/latest/protocol.html?highlight=irc#server-peers-subscribe) on Freenode.

## Usage
See [example/](https://github.com/d4l3k/go-electrum/tree/master/example) for more.

### electrum [![GoDoc](https://godoc.org/github.com/d4l3k/go-electrum/electrum?status.svg)](https://godoc.org/github.com/d4l3k/go-electrum/electrum)
```bash
$ go get -u github.com/d4l3k/go-electrum/electrum
```

```go
package main

import (
  "log"

  "github.com/d4l3k/go-electrum/electrum"
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

### wallet [![GoDoc](https://godoc.org/github.com/d4l3k/go-electrum/wallet?status.svg)](https://godoc.org/github.com/d4l3k/go-electrum/wallet)

```bash
$ go get -u github.com/d4l3k/go-electrum/wallet
```

```go
package main

import (
  "log"

  "github.com/btcsuite/btcutil"
  "github.com/d4l3k/go-electrum/wallet"
)

func main() {
  seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
  if err != nil {
    log.Fatal(err)
  }
  w, err := wallet.Create("test.wallet", "pass", seed)
  if err != nil {
    log.Fatal(err)
  }
  addrs, err := w.GenAddresses(1)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Address: %s", addrs[0])
  err = w.SendBitcoin(map[string]btcutil.Amount{
    "18mS21JLSWJcTwKV8ZEv5SvroKAqkbYfPy": btcutil.NewAmount(1.0),
  }, 6)
  if err != nil {
    log.Fatal(err)
  }
}
```

### irc [![GoDoc](https://godoc.org/github.com/d4l3k/go-electrum/irc?status.svg)](https://godoc.org/github.com/d4l3k/go-electrum/irc)
```bash
$ go get -u github.com/d4l3k/go-electrum/irc
```

```go
package main

import (
	"log"

	"github.com/d4l3k/go-electrum/irc"
)

func main() {
	log.Println(irc.FindElectrumServers())
}
```



# License
go-electrum is licensed under the MIT license.

Made by [Tristan Rice](https://fn.lc).
