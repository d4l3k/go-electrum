/*
Package wallet provides a simple interface to btcwallet and electrum.

This is designed to make accepting and sending bitcoin really easy using Go.
*/
package wallet

import (
	"log"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/netparams"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/wallet"
	"github.com/btcsuite/btcwallet/walletdb"
	"github.com/btcsuite/btcwallet/wtxmgr"
	"github.com/d4l3k/go-electrum/electrum"

	_ "github.com/btcsuite/btcwallet/walletdb/bdb"
)

var (
	waddrmgrNamespaceKey = []byte("waddrmgrNamespace")
	wtxmgrNamespaceKey   = []byte("wtxmgr")

	bitcoinNetwork = &netparams.MainNetParams
)

type Wallet struct {
	wallet *wallet.Wallet
	node   *electrum.Node
}

// Addresses returns all addresses generated in the current bitcoin wallet.
func (w *Wallet) Addresses() ([]btcutil.Address, error) {
	acc, err := w.wallet.Manager.LastAccount()
	if err != nil {
		return nil, err
	}
	var addrs []btcutil.Address
	err = w.wallet.Manager.ForEachAccountAddress(acc, func(maddr waddrmgr.ManagedAddress) error {
		addrs = append(addrs, maddr.Address())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return addrs, nil
}

// GenAddresses generates a number of addresses for the wallet.
func (w *Wallet) GenAddresses(n int) ([]btcutil.Address, error) {
	acc, err := w.wallet.Manager.LastAccount()
	if err != nil {
		return nil, err
	}
	managedAddrs, err := w.wallet.Manager.NextExternalAddresses(acc, uint32(n))
	if err != nil {
		return nil, err
	}
	return stripManagedAddrs(managedAddrs), nil
}

// SendBitcoin sends some amount of bitcoin specifying minimum confirmations.
func (w *Wallet) SendBitcoin(amounts map[string]btcutil.Amount, minconf int) error {
	account, err := w.wallet.Manager.LastAccount()
	if err != nil {
		return err
	}

	log.Printf("creating tx")

	// TODO: make this work

	// taken from https://github.com/btcsuite/btcwallet/blob/master/wallet/wallet.go SendPairs

	// Create transaction, replying with an error if the creation
	// was not successful.
	createdTx, err := w.wallet.CreateSimpleTx(account, amounts, int32(minconf))
	if err != nil {
		return err
	}

	log.Printf("created tx %#v", createdTx)

	// Create transaction record and insert into the db.
	rec, err := wtxmgr.NewTxRecordFromMsgTx(createdTx.MsgTx, time.Now())
	if err != nil {
		log.Printf("Cannot create record for created transaction: %v", err)
		return err
	}
	log.Printf("new txrecord %#v", rec)
	err = w.wallet.TxStore.InsertTx(rec, nil)
	if err != nil {
		log.Printf("Error adding sent tx history: %v", err)
		return err
	}
	log.Printf("inserted tx")

	if createdTx.ChangeIndex >= 0 {
		err = w.wallet.TxStore.AddCredit(rec, nil, uint32(createdTx.ChangeIndex), true)
		if err != nil {
			log.Printf("Error adding change address for sent "+"tx: %v", err)
			return err
		}
	}

	log.Printf("broadcasting")

	resp, err := w.node.BlockchainTransactionBroadcast(rec.SerializedTx)
	if err != nil {
		return err
	}

	log.Printf("RESP broadcast %#v", resp)

	return nil
}

func stripManagedAddrs(mAddrs []waddrmgr.ManagedAddress) []btcutil.Address {
	addrs := make([]btcutil.Address, len(mAddrs))
	for i, addr := range mAddrs {
		addrs[i] = addr.Address()
	}
	return addrs
}

// Create creates a wallet with the specified path, private key password and seed.
// Seed can be created using: hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
func Create(path, privPass string, seed []byte) (*Wallet, error) {
	db, err := walletdb.Create("bdb", path)
	if err != nil {
		return nil, err
	}
	namespace, err := db.Namespace(waddrmgrNamespaceKey)
	if err != nil {
		return nil, err
	}
	manager, err := waddrmgr.Create(namespace, seed, nil,
		[]byte(privPass), bitcoinNetwork.Params, nil)
	if err != nil {
		return nil, err
	}
	manager.Close()

	return openWallet(db, privPass, seed)
}

func returnBytes(bytes []byte) func() ([]byte, error) {
	return func() ([]byte, error) {
		return bytes, nil
	}
}

// Load loads a wallet with the specified path, private key password and seed.
func Load(path, privPass string, seed []byte) (*Wallet, error) {
	db, err := walletdb.Open("bdb", path)
	if err != nil {
		return nil, err
	}
	return openWallet(db, privPass, seed)
}

func openWallet(db walletdb.DB, privPass string, seed []byte) (*Wallet, error) {
	addrMgrNS, err := db.Namespace(waddrmgrNamespaceKey)
	if err != nil {
		return nil, err
	}
	txMgrNS, err := db.Namespace(wtxmgrNamespaceKey)
	if err != nil {
		return nil, err
	}

	cbs := &waddrmgr.OpenCallbacks{
		ObtainSeed:        returnBytes(seed),
		ObtainPrivatePass: returnBytes([]byte(privPass)),
	}
	backWallet, err := wallet.Open(nil, bitcoinNetwork.Params, db, addrMgrNS, txMgrNS, cbs)
	if err != nil {
		return nil, err
	}

	// TODO: use more than 1 node
	node := electrum.NewNode()
	if err := node.ConnectTCP("btc.mustyoshi.com:50001"); err != nil {
		return nil, err
	}

	w := &Wallet{
		wallet: backWallet,
		node:   node,
	}

	addrs, err := w.Addresses()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if err := w.watchAddress(addr.String()); err != nil {
			return nil, err
		}
	}

	return w, nil
}

func (w *Wallet) watchAddress(addr string) error {
	c, err := w.node.BlockchainAddressSubscribe(addr)
	if err != nil {
		return err
	}
	// TODO(d4l3k) handle history
	// w.node.BlockchainAddressGetHistory
	go w.handleTransactions(c)
	return nil
}

func (w *Wallet) handleTransactions(c <-chan string) {
	var err error
	for txid := range c {
		var tx string
		if tx, err = w.node.BlockchainTransactionGet(txid); err != nil {
			break
		}
		if err = w.insertTx(tx); err != nil {
			break
		}
	}
	if err != nil {
		log.Println(err)
		// TODO(d4l3k): Better error handling.
	}
}

func (w *Wallet) insertTx(tx string) error {
	rec, err := wtxmgr.NewTxRecord([]byte(tx), time.Now())
	if err != nil {
		return err
	}
	return w.wallet.TxStore.InsertTx(rec, nil)
}
