package crawler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test get latest block number no error
func TestEthereumClient_BlockNumber(t *testing.T) {
	cli := NewEthereumClient(EthNodeUrl)

	blockNumber, err := cli.BlockNumber(context.TODO())
	assert.NoError(t, err)
	assert.Greater(t, blockNumber, uint64(0))
	t.Logf("block number %d", blockNumber)
}

func TestEthereumClient_GetBlockByNumber(t *testing.T) {
	ctx := context.TODO()

	cli := NewEthereumClient(EthNodeUrl)

	blockNumber, err := cli.BlockNumber(ctx)
	assert.NoError(t, err)

	block, err := cli.GetBlockByNumber(ctx, blockNumber)
	assert.NoError(t, err)
	assert.Equal(t, blockNumber, uint64(block.Number))
}
