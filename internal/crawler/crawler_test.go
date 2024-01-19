package crawler

import (
	"context"
	"testing"
	"time"

	"github.com/TrustWallet/tx-parser/internal/mocks"
	"github.com/TrustWallet/tx-parser/internal/types"
	"github.com/TrustWallet/tx-parser/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_getLatestBlock_success(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	repo.On("GetCurrentBlock", ctx).Return(uint64(12), nil)

	cli := mocks.NewClient(t)
	cli.On("BlockNumber", ctx).Return(uint64(13), nil)
	fakeBlock := types.Block{
		Number: utils.HexUint64(13),
	}
	cli.On("GetBlockByNumber", ctx, uint64(13)).Return(&fakeBlock, nil)

	crawler := ethereumCrawler{
		repo: repo,
		cli:  cli,
	}

	block, err := crawler.getLatestBlock(ctx)
	assert.NoError(t, err)
	assert.NotEqual(t, fakeBlock, block)
}

func Test_getLatestBlock_errDuplicate(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	repo.On("GetCurrentBlock", ctx).Return(uint64(13), nil)

	cli := mocks.NewClient(t)
	cli.On("BlockNumber", ctx).Return(uint64(13), nil)

	crawler := ethereumCrawler{
		repo: repo,
		cli:  cli,
	}

	block, err := crawler.getLatestBlock(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicateParsed)
	assert.Nil(t, block)
}

func Test_extractTransactions(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	repo.On("GetAddresses", ctx).Return([]string{"test1", "TEST2", "test3"}, nil)

	cli := mocks.NewClient(t)
	var fakeBlock = types.Block{
		Number: utils.HexUint64(14),
		Transactions: []types.Transaction{
			{
				BlockNumber: utils.HexUint64(14),
				From:        "0x123",
				To:          "0x321",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
			{
				BlockNumber: utils.HexUint64(14),
				From:        "TEST1",
				To:          "test1",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
			{
				BlockNumber: utils.HexUint64(14),
				From:        "TEST2",
				To:          "test1",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
			{
				BlockNumber: utils.HexUint64(14),
				From:        "TEst1",
				To:          "test2",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
		},
		Timestamp: utils.HexUint64(time.Now().Unix()),
	}

	crawler := ethereumCrawler{
		repo: repo,
		cli:  cli,
	}

	txns, err := crawler.extractTransactions(ctx, &fakeBlock)
	assert.NoError(t, err)
	assert.Len(t, txns, 3)

	for _, tx := range txns {
		assert.NotEqual(t, "0x123", tx.From)
	}
}

func TestEthereumCrawler_Run(t *testing.T) {
	ctx := context.TODO()
	repo := mocks.NewRepository(t)
	repo.On("GetCurrentBlock", ctx).Return(uint64(13), nil)
	repo.On("GetAddresses", ctx).Return([]string{"test1", "TEST2", "test3"}, nil)
	repo.On("SaveTransactions", ctx, uint64(14), mock.Anything).Return(nil)

	cli := mocks.NewClient(t)
	cli.On("BlockNumber", ctx).Return(uint64(14), nil)
	var fakeBlock = types.Block{
		Number: utils.HexUint64(14),
		Transactions: []types.Transaction{
			{
				BlockNumber: utils.HexUint64(14),
				From:        "0x123",
				To:          "0x321",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
			{
				BlockNumber: utils.HexUint64(14),
				From:        "TEST1",
				To:          "test1",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
			{
				BlockNumber: utils.HexUint64(14),
				From:        "TEST2",
				To:          "test1",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
			{
				BlockNumber: utils.HexUint64(14),
				From:        "TEst1",
				To:          "test2",
				Timestamp:   utils.HexUint64(time.Now().Unix()),
			},
		},
		Timestamp: utils.HexUint64(time.Now().Unix()),
	}
	cli.On("GetBlockByNumber", ctx, uint64(14)).Return(&fakeBlock, nil)

	crawler := NewEthereumCrawler(repo, cli)
	err := crawler.Run(ctx)
	assert.NoError(t, err)

	repo.AssertNumberOfCalls(t, "GetCurrentBlock", 1)
	repo.AssertNumberOfCalls(t, "GetAddresses", 1)
	repo.AssertNumberOfCalls(t, "SaveTransactions", 1)
	cli.AssertNumberOfCalls(t, "BlockNumber", 1)
	cli.AssertNumberOfCalls(t, "GetBlockByNumber", 1)
}
