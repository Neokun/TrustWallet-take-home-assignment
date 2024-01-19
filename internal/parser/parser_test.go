package parser

import (
	"fmt"
	"testing"

	"github.com/TrustWallet/tx-parser/internal/mocks"
	"github.com/TrustWallet/tx-parser/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParserService_GetCurrentBlock(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.On("GetCurrentBlock", mock.Anything).Return(uint64(10), nil)
	parser := NewParserService(repo)
	num := parser.GetCurrentBlock()
	assert.Equal(t, 10, num)
}

func TestParserService_GetCurrentBlock_error(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.On("GetCurrentBlock", mock.Anything).Return(uint64(0), fmt.Errorf("some error"))

	parser := NewParserService(repo)
	num := parser.GetCurrentBlock()
	assert.Equal(t, 0, num)
}

func TestParserService_GetTransactions(t *testing.T) {
	repo := mocks.NewRepository(t)

	fakeTxns := []types.Transaction{
		{
			BlockNumber: 10,
			From:        "test",
		},
		{
			BlockNumber: 10,
			To:          "test",
		},
	}
	repo.On("GetTransactions", mock.Anything, "test").Return(fakeTxns, nil)
	parser := NewParserService(repo)
	transactions := parser.GetTransactions("test")
	assert.Len(t, transactions, 2)
}

func TestParserService_GetTransactions_error(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.On("GetTransactions", mock.Anything, "test").Return(nil, nil)
	parser := NewParserService(repo)
	transactions := parser.GetTransactions("test")
	assert.Len(t, transactions, 0)
}

func TestParserService_Subscribe(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.On("AddAddress", mock.Anything, "test").Return(nil)
	parser := NewParserService(repo)
	ok := parser.Subscribe("test")
	assert.True(t, ok)

	repo.On("AddAddress", mock.Anything, "test1").Return(fmt.Errorf("address already exists"))
	ok = parser.Subscribe("test1")
	assert.False(t, ok)
}
