package wallet

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
)

func TestWalletCreationAndLoad(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "wallet.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	privPass := "testtest"
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := CreateWallet(file.Name(), privPass, seed)
	if err != nil {
		t.Fatal(err)
	}

	if addrs, err := wallet.Addresses(); err != nil {
		t.Fatal(err)
	} else if len(addrs) != 0 {
		t.Fatalf("wallet doesn't start with 0 addresses, len = %d", len(addrs))
	}

	if addrs, err := wallet.GenAddresses(10); err != nil {
		t.Fatal(err)
	} else if len(addrs) != 10 {
		t.Fatalf("generated wrong number of addresses, len = %d", len(addrs))
	}

	if addrs, err := wallet.Addresses(); err != nil {
		t.Fatal(err)
	} else if len(addrs) != 10 {
		t.Fatalf("wallet doesn't have new addresses, len = %d", len(addrs))
	} else {
		for _, addr := range addrs {
			fmt.Printf("addr %s\n", addr.String())
		}
	}
}
