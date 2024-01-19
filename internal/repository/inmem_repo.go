package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/TrustWallet/tx-parser/internal/types"
)

type addressTransactionsDict map[string][]*types.Transaction

type inMemRepo struct {
	mu              sync.RWMutex
	txnDict         addressTransactionsDict
	currentBlockNum uint64
}

func NewInMemRepo() *inMemRepo {
	txnDict := make(addressTransactionsDict)
	return &inMemRepo{
		txnDict:         txnDict,
		currentBlockNum: 0,
	}
}

// GetCurrentBlock return last parsed block number
func (r *inMemRepo) GetCurrentBlock(ctx context.Context) (uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.currentBlockNum, nil
}

// GetAddresses get list of subscribed addresses
func (r *inMemRepo) GetAddresses(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	addresses := make([]string, 0, len(r.txnDict))
	for addr := range r.txnDict {
		addresses = append(addresses, addr)
	}

	return addresses, nil
}

// GetTransactions return transactions for an address
func (r *inMemRepo) GetTransactions(ctx context.Context, address string) ([]types.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	address = strings.ToLower(address)
	if txns, ok := r.txnDict[address]; ok {
		transactions := make([]types.Transaction, len(txns))
		for i, tx := range txns {
			transactions[i] = *tx
		}

		return transactions, nil
	}

	return nil, fmt.Errorf("address not found")
}

// AddAddress add an address to list of subscription
func (r *inMemRepo) AddAddress(ctx context.Context, address string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	address = strings.ToLower(address)
	if _, ok := r.txnDict[address]; ok {
		return fmt.Errorf("address already exists")
	}

	r.txnDict[address] = []*types.Transaction{}
	return nil
}

// SaveTransactions save the list of transactions
func (r *inMemRepo) SaveTransactions(ctx context.Context, blockNumber uint64, txns []*types.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.currentBlockNum = blockNumber
	newTxnDict := make(addressTransactionsDict)
	for address := range r.txnDict {
		newTxnDict[address] = []*types.Transaction{}
	}

	for _, tx := range txns {
		if _, ok := newTxnDict[strings.ToLower(tx.To)]; ok {
			newTxnDict[strings.ToLower(tx.To)] = append(newTxnDict[strings.ToLower(tx.To)], tx)
		}

		// avoid duplicate on self send transactions
		if strings.EqualFold(tx.To, tx.From) {
			continue
		}

		if _, ok := newTxnDict[strings.ToLower(tx.From)]; ok {
			newTxnDict[strings.ToLower(tx.From)] = append(newTxnDict[strings.ToLower(tx.From)], tx)
		}
	}

	r.txnDict = newTxnDict
	return nil
}
