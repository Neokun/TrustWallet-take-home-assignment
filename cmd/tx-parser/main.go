package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/TrustWallet/tx-parser/internal/api/v1"
	"github.com/TrustWallet/tx-parser/internal/crawler"
	"github.com/TrustWallet/tx-parser/internal/parser"
	"github.com/TrustWallet/tx-parser/internal/repository"
)

type runFn func(ctx context.Context) error

func main() {
	repo := repository.NewInMemRepo()
	cli := crawler.NewEthereumClient(crawler.EthNodeUrl)
	crawler := crawler.NewEthereumCrawler(repo, cli)
	parser := parser.NewParserService(repo)
	register := api.NewRegister(parser)

	// Run the interval job
	_ = runner(crawler.Run)

	// Run the APIs
	http.HandleFunc("/subscribe", register.SubscribeHandler)
	http.HandleFunc("/current-block", register.GetCurrentBlockHandler)
	http.HandleFunc("/transactions", register.GetTransactionsHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func runner(fn runFn) chan<- struct{} {
	// Create a ticker that ticks every 4 seconds.
	ticker := time.NewTicker(4 * time.Second)
	// Use a channel to signal the stop of the program.
	quit := make(chan struct{})

	// Start a goroutine that executes job.
	go func() {
		for {
			select {
			case <-ticker.C:
				err := fn(context.Background())
				if err != nil {
					log.Printf("Error executing job: %v", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
