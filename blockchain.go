package electrum

// BlockchainNumBlocksSubscribe returns the current number of blocks.
// http://docs.electrum.org/en/latest/protocol.html#blockchain-numblocks-subscribe
func (n *Node) BlockchainNumBlocksSubscribe() (int, error) {
	resp := &struct {
		Result int `json:"result"`
	}{}
	err := n.request("blockchain.numblocks.subscribe", nil, resp)
	return resp.Result, err
}

type BlockchainHeader struct {
	Nonce         uint64 `json:"nonce"`
	PrevBlockHash string `json:"prev_block_hash"`
	Timestamp     uint64 `json:"timestamp"`
	MerkleRoot    string `json:"merkle_root"`
	BlockHeight   uint64 `json:"block_height"`
	UtxoRoot      string `json:"utxo_root"`
	Version       int    `json:"version"`
	Bits          uint64 `json:"bits"`
}

// BlockchainHeadersSubscribe request client notifications about new blocks in
// form of parsed blockheaders and returns the current block header.
// TODO (d4l3k) handle header updates
// http://docs.electrum.org/en/latest/protocol.html#blockchain-headers-subscribe
func (n *Node) BlockchainHeadersSubscribe() (*BlockchainHeader, error) {
	resp := &struct {
		Result *BlockchainHeader `json:"result"`
	}{}
	err := n.request("blockchain.headers.subscribe", nil, resp)
	return resp.Result, err
}

// BlockchainAddressSubscribe subscribes to transactions on an address and
// returns the hash of the transaction history.
// TODO (d4l3k) handle address updates
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-subscribe
func (n *Node) BlockchainAddressSubscribe(address string) (string, error) {
	resp := &basicResp{}
	err := n.request("blockchain.address.subscribe", []string{address}, resp)
	return resp.Result, err
}

type Transaction struct {
	Hash   string `json:"tx_hash"`
	Height int    `json:"height"`
	Value  int    `json:"value"`
	Pos    int    `json:"tx_pos"`
}

// BlockchainAddressGetHistory returns the history of an address.
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-history
func (n *Node) BlockchainAddressGetHistory(address string) ([]*Transaction, error) {
	resp := &struct {
		Result []*Transaction `json:"result"`
	}{}
	err := n.request("blockchain.address.get_history", []string{address}, resp)
	return resp.Result, err
}

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-mempool
func (n *Node) BlockchainAddressGetMempool() error { return ErrNotImplemented }

type Balance struct {
	Confirmed   int `json:"confirmed"`
	Unconfirmed int `json:"unconfirmed"`
}

// BlockchainAddressGetBalance returns the balance of an address.
// TODO (d4l3k) investigate `error from server: "'Node' object has no attribute '__getitem__'"`
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-balance
func (n *Node) BlockchainAddressGetBalance(address string) (*Balance, error) {
	resp := &struct {
		Result *Balance `json:"result"`
	}{}
	err := n.request("blockchain.address.get_balance", []string{address}, resp)
	return resp.Result, err
}

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-proof
func (n *Node) BlockchainAddressGetProof() error { return ErrNotImplemented }

// BlockchainAddressListUnspent lists the unspent transactions for the given address.
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-listunspent
func (n *Node) BlockchainAddressListUnspent(address string) ([]*Transaction, error) {
	resp := &struct {
		Result []*Transaction `json:"result"`
	}{}
	err := n.request("blockchain.address.listunspent", []string{address}, resp)
	return resp.Result, err
}

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-utxo-get-address
func (n *Node) BlockchainUtxoGetAddress() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-block-get-header
func (n *Node) BlockchainBlockGetHeader() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-block-get-chunk
func (n *Node) BlockchainBlockGetChunk() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-broadcast
func (n *Node) BlockchainTransactionBroadcast() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-get-merkle
func (n *Node) BlockchainTransactionGetMerkle() error { return ErrNotImplemented }

// BlockchainTransactionGet returns the raw transaction (hex-encoded) for the given txid. If transaction doesn't exist, an error is returned.
// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-get
func (n *Node) BlockchainTransactionGet(txid string) (string, error) {
	resp := &basicResp{}
	err := n.request("blockchain.transaction.get", []string{txid}, resp)
	return resp.Result, err
}
