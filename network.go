package electrum

import (
	"encoding/json"
	"errors"
	"fmt"
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

type respMetadata struct {
	Id    int    `json:"id"`
	Error string `json:"error"`
}

type request struct {
	Id     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type basicResp struct {
	Result string `json:"result"`
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
			msg := &respMetadata{}
			if err := json.Unmarshal(bytes, msg); err != nil {
				n.err(err)
				return
			}
			if len(msg.Error) > 0 {
				n.err(fmt.Errorf("error from server: %#v", msg.Error))
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

func (n *Node) request(method string, params []string, v interface{}) error {
	msg := request{
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
