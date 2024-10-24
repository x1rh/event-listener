package eventlistener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/x1rh/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventListener struct {
	client *ethclient.Client

	Config   ChainConfig
	Contract *Contract
	logChan  chan types.Log

	opts *EventListenerOptions
}

func New(c ChainConfig, options ...Option) (*EventListener, error) {
	opts := &EventListenerOptions{}
	for _, option := range options {
		option(opts)
	}

	var client *ethclient.Client
	var err error

	if opts.Client != nil {
		client = opts.Client
	} else if opts.URL != "" {
		client, err = ethclient.Dial(opts.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
		}
	} else {
		return nil, fmt.Errorf("either URL or ethclient.Client must be provided")
	}

	el := &EventListener{
		Config:  c,
		client:  client,
		logChan: make(chan types.Log, 256),
		opts:    opts,
	}

	if opts.Contract != nil {
		el.Contract = opts.Contract
	}
	if el.Contract == nil {
		return nil, errors.New("invali Contract")
	}

	return el, nil
}

func (el *EventListener) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(el.Contract.Address)},
		}
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ticker.C:
				// todo: check if it filters in interval [FromBlock, ToBlock] or [FromBlock, ToBlock)
				toBlock := big.NewInt(0).Add(el.Contract.BlockNumber, el.Contract.Step)
				mostRecentBlockNumber, err := el.client.BlockNumber(ctx)
				if err != nil {
					slog.Error("failed to get newest block number", "err", err)
					continue
				}
				if toBlock.Uint64() > mostRecentBlockNumber {
					toBlock = big.NewInt(int64(mostRecentBlockNumber))
				}

				query.FromBlock = el.Contract.BlockNumber
				query.ToBlock = toBlock

				if query.FromBlock.Cmp(query.ToBlock) > 0 {
					continue
				}

				slog.Info("handle block", slog.Any("fromBlock", query.FromBlock), slog.Any("toBlock", query.ToBlock))

				logList, err := el.client.FilterLogs(ctx, query)
				if err != nil {
					slog.Error(
						"Failed to query logs",
						"error", err,
					)
					continue
				}

				ok := true
				for _, txLog := range logList {
					if err := el.Contract.Handle(ctx, &txLog); err != nil {
						slog.Error("fail to handle event", slog.Any("err", err))
						ok = false
						break
					}
				}
				if ok {
					el.Contract.BlockNumber = big.NewInt(0).Add(toBlock, big.NewInt(1))
				}
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("received shutdown signal")
	cancel()
	el.Stop()

	log.Println("All goroutines have finished, exiting main function")
}

func (el *EventListener) Stop() {

}
