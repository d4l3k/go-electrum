package electrum

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
)

const (
	ClientVersion   = "0.0.1"
	ProtocolVersion = "1.0"
)

var ErrNotImplemented = errors.New("error not implemented")

type Transport interface {
	SendMessage([]byte) error
	Responses() <-chan []byte
	Errors() <-chan error
}

func FindElectrumServersIRC() ([]*Node, error) {
	return nil, ErrNotImplemented
}

type Node struct {
	transport    Transport
	handlers     map[int]chan []byte
	handlersLock sync.RWMutex

	nextId int
}

func NewNode(addr string) (*Node, error) {
	transport, err := NewTCPTransport(addr)
	if err != nil {
		return nil, err
	}
	n := &Node{
		transport: transport,
		handlers:  make(map[int]chan []byte),
	}
	go n.listen()
	return n, nil
}

type resp struct {
	Id int `json:"id"`
}

func (n *Node) err(err error) {
	// TODO (d4l3k) Better error handling.
	log.Fatal(err)
}

func (n *Node) listen() {
	for {
		select {
		case err := <-n.transport.Errors():
			n.err(err)
			return
		case bytes := <-n.transport.Responses():
			msg := &resp{}
			if err := json.Unmarshal(bytes, msg); err != nil {
				n.err(err)
				return
			}

			n.handlersLock.RLock()
			c, ok := n.handlers[msg.Id]
			n.handlersLock.RUnlock()

			if ok {
				c <- bytes
			}
		}
	}
}

type message struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

func (n *Node) request(method string, params []string, v interface{}) error {
	msg := message{
		Id:     n.nextId,
		Method: method,
		Params: params,
	}
	n.nextId++
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	bytes = append(bytes, delim)
	if err := n.transport.SendMessage(bytes); err != nil {
		return err
	}

	c := make(chan []byte, 1)

	n.handlersLock.Lock()
	n.handlers[msg.Id] = c
	n.handlersLock.Unlock()

	resp := <-c

	n.handlersLock.Lock()
	defer n.handlersLock.Unlock()
	delete(n.handlers, msg.Id)

	if err := json.Unmarshal(resp, v); err != nil {
		return nil
	}
	return nil
}

type basicResp struct {
	Result string `json:"result"`
}

// ServerVersion returns the server's version.
// http://docs.electrum.org/en/latest/protocol.html#server-version
func (n *Node) ServerVersion() (string, error) {
	resp := &basicResp{}
	err := n.request("server.version", []string{ClientVersion, ProtocolVersion}, resp)
	return resp.Result, err
}

// ServerBanner returns the server's banner.
// http://docs.electrum.org/en/latest/protocol.html#server-banner
func (n *Node) ServerBanner() (string, error) {
	resp := &basicResp{}
	err := n.request("server.banner", []string{ClientVersion, ProtocolVersion}, resp)
	return resp.Result, err
}

// ServerDonationAddress returns the donation address of the server.
// http://docs.electrum.org/en/latest/protocol.html#server-donation-address
func (n *Node) ServerDonationAddress() (string, error) {
	resp := &basicResp{}
	err := n.request("server.donation_address", []string{ClientVersion, ProtocolVersion}, resp)
	return resp.Result, err
}

// ServerPeersSubscribe requests peers from a server.
// http://docs.electrum.org/en/latest/protocol.html#server-peers-subscribe
func (n *Node) ServerPeersSubscribe() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-numblocks-subscribe
func (n *Node) BlockchainNumblocksSubscribe() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-headers-subscribe
func (n *Node) BlockchainHeadersSubscribe() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-subscribe
func (n *Node) BlockchainAddressSubscribe() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-history
func (n *Node) BlockchainAddressGetHistory() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-mempool
func (n *Node) BlockchainAddressGetMempool() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-balance
func (n *Node) BlockchainAddressGetBalance() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-proof
func (n *Node) BlockchainAddressGetProof() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-listunspent
func (n *Node) BlockchainAddressListunspent() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-utxo-get-address
func (n *Node) BlockchainUtxoGetAddress() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-block-get-header
func (n *Node) BlockchainBlockGetHeader() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-block-get-chunk
func (n *Node) BlockchainBlockGetChunk() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-broadcast
func (n *Node) BlockchainTransactionBroadcast() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-get-merkle
func (n *Node) BlockchainTransactionGetMerkle() error { return ErrNotImplemented }

// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-get
func (n *Node) BlockchainTransactionGet() error { return ErrNotImplemented }
