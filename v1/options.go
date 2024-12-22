package eventlistener

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventListenerOptions struct {
	URL      string
	Client   *ethclient.Client
	Contract *Contract
}

type Option func(*EventListenerOptions)

func WithURL(url string) Option {
	return func(opts *EventListenerOptions) {
		if url != "" {
			opts.URL = url
		}
	}
}

func WithClient(client *ethclient.Client) Option {
	return func(opts *EventListenerOptions) {
		if client != nil {
			opts.Client = client
		}
	}
}

func WithContract(c *Contract) Option {
	return func(opts *EventListenerOptions) {
		opts.Contract = c
	}
}
