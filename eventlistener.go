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

	_ "github.com/x1rh/event-listener/logger"

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
				query.FromBlock = el.Contract.BlockNumber
				query.ToBlock = toBlock
				slog.Info("handle block", slog.Any("fromBlock", query.FromBlock), slog.Any("toBlock", query.ToBlock))

				logList, err := el.client.FilterLogs(ctx, query)
				if err != nil {
					slog.Error(
						"Failed to query logs",
						"error", err,
					)
					continue
				}

				for _, vLog := range logList {
					// TODO: only handle topic
					// topic[0] is always a signature when a event is topic
					eventSig := vLog.Topics[0]
					el.Contract.Abi.EventByID(eventSig)

					event, err := el.Contract.Abi.EventByID(vLog.Topics[0])
					if err != nil {
						slog.Error("fail to get even", slog.Any("err", err))
						continue
					}

					eventInfo := &Event{
						Name:          event.Name,
						IndexedParams: make([]common.Hash, len(vLog.Topics)-1),
						Data:          vLog.Data,
						Outputs:       nil,
					}
					slog.Debug("event", slog.Any("event", event))

					// topic[1:] is other indexed params in event
					if len(vLog.Topics) > 1 {
						for i, param := range vLog.Topics[1:] {
							eventInfo.IndexedParams[i] = param
							slog.Debug("", event.Inputs[i].Name, common.HexToAddress(param.Hex()))
						}
					}
					if len(vLog.Data) > 0 {
						outputDataMap := make(map[string]interface{})
						err = el.Contract.Abi.UnpackIntoMap(outputDataMap, event.Name, vLog.Data)
						if err != nil {
							slog.Error("fail to unpack", slog.Any("err", err))
							continue
						}
						eventInfo.Outputs = outputDataMap
					}

					slog.Debug(
						"hanle",
						slog.String("chainName", el.Config.ChainName),
						slog.String("contractName", el.Contract.Name),
						slog.String("ContractAddress", el.Contract.Address),
						slog.Any("block number", vLog.BlockNumber),
					)

					if err := el.Contract.Handle(ctx, eventInfo); err != nil {
						slog.Error("fail to handle event", slog.Any("err", err))
						continue
					}
				}
				el.Contract.BlockNumber = toBlock // todo: check it
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
