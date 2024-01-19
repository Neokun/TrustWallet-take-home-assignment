package parser

import (
	"context"
	"log"

	"github.com/TrustWallet/tx-parser/internal/repository"
	"github.com/TrustWallet/tx-parser/internal/types"
)

type Parser interface {
	// GetCurrentBlock return last parsed block
	GetCurrentBlock() int

	// Subscribe add address to observer
	Subscribe(address string) bool

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []types.Transaction
}

type parserService struct {
	repo repository.Repository
}

func NewParserService(repo repository.Repository) *parserService {
	return &parserService{repo: repo}
}

// GetCurrentBlock return last parsed block
func (p *parserService) GetCurrentBlock() int {
	blockNum, err := p.repo.GetCurrentBlock(context.Background())
	if err != nil {
		log.Printf("Error getting current block: %v", err)
	}

	return int(blockNum)
}

// Subscribe add address to observer
func (p *parserService) Subscribe(address string) bool {
	err := p.repo.AddAddress(context.Background(), address)
	if err != nil {
		log.Printf("Error subcribe address %s: %v", address, err)
		return false
	}

	return true
}

// GetTransactions list of inbound or outbound transactions for an address
func (p *parserService) GetTransactions(address string) []types.Transaction {
	txns, err := p.repo.GetTransactions(context.Background(), address)
	if err != nil {
		log.Printf("Error get transactions for address %s: %v", address, err)
		return nil
	}

	return txns
}
