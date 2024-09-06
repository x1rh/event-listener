package eventlistener

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventListener struct {
	client *ethclient.Client

	Config          ChainConfig
	ContractMap     map[string]Contract
	subscriptionMap map[string]ethereum.Subscription
	logChanMap      map[string]chan types.Log

	opts *EventListenerOptions
}

func New(c ChainConfig, options ...Option) (*EventListener, error) {
	opts := &EventListenerOptions{ContractMap: make(map[string]Contract)}
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

	return &EventListener{
		Config:      c,
		client:      client,
		logChanMap:  make(map[string]chan types.Log),
		ContractMap: opts.ContractMap,
		opts:        opts,
	}, nil
}

func (el *EventListener) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for address, contract := range el.ContractMap {
		// log.Printf("address=%+v, contract=%+v\n", address, contract)
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(address)},
		}

		logChan := make(chan types.Log, 1024)
		sub, err := el.client.SubscribeFilterLogs(ctx, query, logChan)
		if err != nil {
			panic(fmt.Sprintf("Failed to subscribe to logs: %v", err))
		}
		el.logChanMap[address] = logChan
		el.subscriptionMap[address] = sub

		wg.Add(1)
		go func(contract Contract) {
			defer wg.Done()
			for {
				select {
				case err := <-sub.Err():
					log.Printf("Subscription error: %v\n", err)
				case vLog := <-logChan:
					eventSig := vLog.Topics[0]
					eventData, ok := contract.EventMap[eventSig]
					if !ok {
						log.Printf("Unknown event signature: %s", eventSig.Hex())
						continue
					}

					fmt.Printf("Event: %s\n", eventData.Name)

					eventDataValues, err := contract.Abi.Unpack(eventData.Name, vLog.Data)
					if err != nil {
						log.Printf("Failed to unpack log data: %v", err)
						continue
					}

					for i, input := range eventData.Inputs {
						if input.Indexed {
							fmt.Printf("%s: %s\n", input.Name, common.HexToAddress(vLog.Topics[i+1].Hex()).Hex())
						} else {
							fmt.Printf("%s: %v\n", input.Name, eventDataValues[i])
						}
					}
					log.Printf(
						"chainName=%s contractName=%s contractAddress=%s block number=%v, done\n",
						el.Config.ChainName, contract.Name, contract.Address, vLog.BlockNumber,
					)
				}
			}
		}(contract)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("received shutdown signal")
		cancel()
		el.Stop()
	}()

	wg.Wait()
	log.Println("All goroutines have finished, exiting main function")
}

func (el *EventListener) Stop() {
	for _, sub := range el.subscriptionMap {
		sub.Unsubscribe()
	}
	for _, ch := range el.logChanMap {
		close(ch)
	}
}
