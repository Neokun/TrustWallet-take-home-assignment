package types

import (
	"github.com/TrustWallet/tx-parser/internal/utils"
)

// Block contains information of block, an add more fields if needed
type Block struct {
	Number       utils.HexUint64 `json:"number"`
	Hash         string          `json:"hash"`
	ParentHash   string          `json:"parentHash"`
	Transactions []Transaction   `json:"transactions"`
	Timestamp    utils.HexUint64 `json:"timestamp"`
}
