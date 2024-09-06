package eventlistener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-errors/errors"
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
				toBlock := big.NewInt(0).Add(el.Contract.BlockNumber, el.Contract.Step)
				query.FromBlock = el.Contract.BlockNumber
				query.ToBlock = toBlock

				logList, err := el.client.FilterLogs(ctx, query)
				if err != nil {
					log.Printf("Failed to query logs: %v", err)
					continue
				}

				for _, vLog := range logList {
					eventSig := vLog.Topics[0]
					el.Contract.Abi.EventByID(eventSig)

					event, err := el.Contract.Abi.EventByID(vLog.Topics[0])
					if err != nil {
						log.Println(err, "get event fail")
						continue
					}

					eventInfo := &Event{
						Name:          event.Name,
						IndexedParams: make([]common.Hash, len(vLog.Topics)-1),
						Data:          vLog.Data,
						Outputs:       nil,
					}
					fmt.Printf("event: %+v\n", event)
					// fmt.Printf("eventInfo: %+v\n", eventInfo)

					// topic[1:] is other indexed params in event
					if len(vLog.Topics) > 1 {
						for i, param := range vLog.Topics[1:] {
							// fmt.Printf("Indexed params %d in hex: %s\n", i, param)
							// fmt.Printf("Indexed params %d decoded %s\n", i, common.HexToAddress(param.Hex()))
							fmt.Printf("%s = %s\n", event.Inputs[i].Name, common.HexToAddress(param.Hex()))
							eventInfo.IndexedParams[i] = param
						}
					}
					if len(vLog.Data) > 0 {
						//fmt.Printf("Log Data in Hex: %s\n", hex.EncodeToString(vLog.Data))
						outputDataMap := make(map[string]interface{})
						err = el.Contract.Abi.UnpackIntoMap(outputDataMap, event.Name, vLog.Data)
						if err != nil {
							log.Println(err, "uppack fail")
							continue
						}
						//fmt.Printf("Event outputs: %v\n", outputDataMap)
						eventInfo.Outputs = outputDataMap
						for k, v := range outputDataMap {
							fmt.Println(k, v)
						}
					}

					// fmt.Printf("eventInfo: %+v\n", eventInfo)

					log.Printf(
						"chainName=%s contractName=%s contractAddress=%s block number=%v, done\n\n",
						el.Config.ChainName, el.Contract.Name, el.Contract.Address, vLog.BlockNumber,
					)
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
