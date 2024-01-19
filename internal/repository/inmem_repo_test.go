package repository

import (
	"context"
	"testing"

	"github.com/TrustWallet/tx-parser/internal/types"
	"github.com/TrustWallet/tx-parser/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestInMemRepo_GetAddAddress(t *testing.T) {
	repo := NewInMemRepo()
	addresses, err := repo.GetAddresses(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, addresses, 0)

	err = repo.AddAddress(context.TODO(), "test1")
	assert.NoError(t, err)
	err = repo.AddAddress(context.TODO(), "test2")
	assert.NoError(t, err)

	addresses, err = repo.GetAddresses(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, addresses, 2)
	assert.Equal(t, addresses[0], "test1")
	assert.Equal(t, addresses[1], "test2")

	err = repo.AddAddress(context.TODO(), "test1")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "address already exists")
}

func TestInMemRepo_GetCurrentBlock(t *testing.T) {
	repo := NewInMemRepo()
	repo.currentBlockNum = 13

	blockNum, err := repo.GetCurrentBlock(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, uint64(13), blockNum)
}

func TestInMemRepo_SaveGetTransactions(t *testing.T) {
	repo := NewInMemRepo()
	err := repo.AddAddress(context.TODO(), "test1")
	assert.NoError(t, err)
	err = repo.AddAddress(context.TODO(), "test2")
	assert.NoError(t, err)

	txns, err := repo.GetTransactions(context.TODO(), "test2")
	assert.NoError(t, err)
	assert.Len(t, txns, 0)

	_, err = repo.GetTransactions(context.TODO(), "test3")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "address not found")

	blockNumber := uint64(14)
	fakeTxns := []*types.Transaction{
		{
			BlockNumber: utils.HexUint64(blockNumber),
			From:        "Test2",
			To:          "tesT1",
		},
		{
			BlockNumber: utils.HexUint64(blockNumber),
			From:        "Test2",
			To:          "tesT2",
		},
		{
			BlockNumber: utils.HexUint64(blockNumber),
			From:        "Test1",
			To:          "tesT2",
		},
		{
			BlockNumber: utils.HexUint64(blockNumber),
			From:        "Test1",
			To:          "tesT3",
		},
		{
			BlockNumber: utils.HexUint64(blockNumber),
			From:        "Test3",
			To:          "tesT1",
		},
	}
	err = repo.SaveTransactions(context.TODO(), blockNumber, fakeTxns)
	assert.NoError(t, err)

	txns, err = repo.GetTransactions(context.TODO(), "test2")
	assert.NoError(t, err)
	assert.Len(t, txns, 3)

	txns, err = repo.GetTransactions(context.TODO(), "test1")
	assert.NoError(t, err)
	assert.Len(t, txns, 4)

	_, err = repo.GetTransactions(context.TODO(), "test3")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "address not found")

	num, err := repo.GetCurrentBlock(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, blockNumber, num)
}
