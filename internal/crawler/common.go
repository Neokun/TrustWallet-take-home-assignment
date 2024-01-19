package crawler

import (
	"encoding/json"
)

type method string

const (
	blockNumberMethod      method = "eth_blockNumber"
	getBlockByNumberMethod method = "eth_getBlockByNumber"
)

const EthNodeUrl = "https://cloudflare-eth.com"

// A value of this type can a JSON-RPC request, notification, successful response or
// error response. Which one it is depends on the fields.
type jsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}
