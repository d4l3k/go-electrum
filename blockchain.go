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
	PrevBlockHash []byte `json:"prev_block_hash"`
	Timestamp     uint64 `json:"timestamp"`
	MerkleRoot    []byte `json:"merkle_root"`
	BlockHeight   uint64 `json:"block_height"`
	UtxoRoot      []byte `json:"utxo_root"`
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

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-subscribe
func (n *Node) BlockchainAddressSubscribe() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-history
func (n *Node) BlockchainAddressGetHistory() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-mempool
func (n *Node) BlockchainAddressGetMempool() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-balance
func (n *Node) BlockchainAddressGetBalance() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-get-proof
func (n *Node) BlockchainAddressGetProof() error { return ErrNotImplemented }

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-address-listunspent
func (n *Node) BlockchainAddressListunspent() error { return ErrNotImplemented }

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

// TODO(d4l3k) implement
// http://docs.electrum.org/en/latest/protocol.html#blockchain-transaction-get
func (n *Node) BlockchainTransactionGet() error { return ErrNotImplemented }
