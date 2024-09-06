package eventlistener

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EventListenerOptions struct {
	URL         string
	Client      *ethclient.Client
	ContractMap map[string]Contract // address -> Contract

	FromBlock *big.Int
	ToBlock   *big.Int
}

type Option func(*EventListenerOptions)

func WithURL(url string) Option {
	return func(opts *EventListenerOptions) {
		opts.URL = url
	}
}

func WithClient(client *ethclient.Client) Option {
	return func(opts *EventListenerOptions) {
		opts.Client = client
	}
}

func WithContract(c Contract) Option {
	return func(opts *EventListenerOptions) {
		_, found := opts.ContractMap[c.Address]
		if !found {
			opts.ContractMap[c.Address] = c
		}
	}
}

func WtihFromBlock(fromBlock *big.Int) Option {
	return func(o *EventListenerOptions) {
		if fromBlock != nil {
			o.FromBlock = fromBlock
		}
	}
}

func WithToBlock(toBlock *big.Int) Option {
	return func(o *EventListenerOptions) {
		if toBlock != nil {
			o.ToBlock = toBlock
		}
	}
}
