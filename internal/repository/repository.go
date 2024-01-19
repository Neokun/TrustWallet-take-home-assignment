package repository

import (
	"context"

	"github.com/TrustWallet/tx-parser/internal/types"
)

//go:generate mockery --name Repository --output ../mocks
type Repository interface {
	// GetCurrentBlock return last parsed block number
	GetCurrentBlock(ctx context.Context) (uint64, error)

	// GetAddresses get list of subscribed addresses
	GetAddresses(ctx context.Context) ([]string, error)

	// GetTransactions return transactions for an address
	GetTransactions(ctx context.Context, address string) ([]types.Transaction, error)

	// AddAddress add an address to list of subscription
	AddAddress(ctx context.Context, address string) error

	// SaveTransactions save the list of transactions
	SaveTransactions(ctx context.Context, blockNumber uint64, txns []*types.Transaction) error
}
