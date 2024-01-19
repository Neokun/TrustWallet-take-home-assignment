package crawler

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/TrustWallet/tx-parser/internal/repository"
	"github.com/TrustWallet/tx-parser/internal/types"
)

type Crawler interface {
	Run(ctx context.Context) error
}

type ethereumCrawler struct {
	repo repository.Repository
	cli  Client
}

func NewEthereumCrawler(repo repository.Repository, cli Client) *ethereumCrawler {
	return &ethereumCrawler{repo: repo, cli: cli}
}

func (c *ethereumCrawler) Run(ctx context.Context) error {
	block, err := c.getLatestBlock(ctx)
	if err != nil {
		if errors.Is(err, ErrDuplicateParsed) {
			return nil
		}

		log.Printf("error getting latest block %v", err)
		return err
	}

	txns, err := c.extractTransactions(ctx, block)
	if err != nil {
		log.Printf("extract transactions from block with err: %v", err)
		return err
	}

	return c.saveData(ctx, uint64(block.Number), txns)

}

func (c *ethereumCrawler) getLatestBlock(ctx context.Context) (*types.Block, error) {
	blockNumber, err := c.cli.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	parsedBlockNum, err := c.repo.GetCurrentBlock(ctx)
	if err != nil {
		return nil, err
	}
	if blockNumber == parsedBlockNum {
		log.Printf("duplicate: block %d already parsed", blockNumber)
		return nil, ErrDuplicateParsed
	}

	block, err := c.cli.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (c *ethereumCrawler) extractTransactions(ctx context.Context, block *types.Block) ([]types.Transaction, error) {
	addresses, err := c.repo.GetAddresses(ctx)
	if err != nil {
		return nil, err
	}
	addressDict := make(map[string]struct{}, len(addresses))
	for _, address := range addresses {
		addressDict[strings.ToLower(address)] = struct{}{}
	}

	var txns []types.Transaction
	for _, txn := range block.Transactions {
		if _, ok := addressDict[strings.ToLower(txn.From)]; ok {
			txns = append(txns, txn)
		} else if _, ok = addressDict[strings.ToLower(txn.To)]; ok {
			txns = append(txns, txn)
		}
	}

	return txns, nil
}

func (c *ethereumCrawler) saveData(ctx context.Context, blockNumber uint64, txns []types.Transaction) error {
	err := c.repo.SaveTransactions(ctx, blockNumber, txns)
	if err != nil {
		log.Printf("error saving %d transactions", len(txns))
		return err
	}

	return nil
}
