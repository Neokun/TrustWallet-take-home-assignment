package types

import (
	"github.com/TrustWallet/tx-parser/internal/utils"
)

// Transaction example of a transaction model, can add more fields if needed
type Transaction struct {
	BlockNumber      utils.HexUint64 `json:"blockNumber"`
	BlockHash        string          `json:"blockHash"`
	From             string          `json:"from"`
	To               string          `json:"to"`
	Value            string          `json:"value"`
	Gas              string          `json:"gas"`
	GasPrice         string          `json:"gasPrice"`
	Hash             string          `json:"hash"`
	TransactionIndex string          `json:"transactionIndex"`
	Timestamp        utils.HexUint64 `json:"timestamp"`
}
