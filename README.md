## This is a repository for TrustWallet's take home assignment

## Assignment

The details of the assignment [here](https://trustwallet.notion.site/Backend-Homework-Tx-Parser-abd431fca950427db75d73d90a0244a8)

For the sake of simplicity, here are assumptions based on the requirements of the assignment:

* The project runs in a single thread, parsing blocks one by one. No re-org handling.
* `GetTransactions` return transactions of the latest parsed block, no pagination.
    * Only store the latest block's transaction data for each address. No historical data is saved.
    * If the latest parsed block does not have transactions for the subscribed address -> return empty list.
* Crawl latest block on Ethereum Blockchain, meaning the block time is 12s (approximately) -> let the job run interval every 4s.
* Avoid usage of external libraries: gin-gonic/gin, go-ethereum, etc.
    * Use Go's `net/http` package to make requests to the Ethereum JSONRPC API.
    * Use Go's built-in data structures to store in memory data.
    * Use simple REST API for expose public interface.
    * Let `testify` is the exception because it is used to support write tests and mock. 


## Main components

### Crawler
Interval (every 4 seconds) running job to get latest block from ethereum and extract transactions filtered by subscribed addresses.

Using the *Repository* to get addresses and save block data. 

Call [RPC node]([https://cloudflare-eth.com](https://cloudflare-eth.com/)) to get block data: 
- [eth_blockNumber](https://ethereum.org/en/developers/docs/apis/json-rpc#eth_blocknumber)
- [eth_getBlockByNumber](https://ethereum.org/en/developers/docs/apis/json-rpc#eth_getblockbynumber)

### Parser
Handle and expose public interface for biz logic:

```go
type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}
```

### Repository
Handle storage queries:

```go
type Repository interface {
	// last parsed block 
	GetCurrentBlock() (int, error)

	// get list of subscribed addresses 
	GetAddresses() ([]string, error)

	// list of transactions for an address 
	GetTransactions(address string) ([]Transaction, error)
	
	// add a address to list of subscription
	AddAddress(address string) error
	
	// save the list of transactions
	SaveTransactions(blockNum int, txn []Transaction) error
}
```

## How to run

When run main, the code will start both the crawler job and the APIs server

```bash
go run cmd/tx-parser/main.go
```
