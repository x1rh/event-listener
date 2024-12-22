package eventlistener

import (
	"context"
	"log/slog"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/x1rh/logger"
)

// LogTokenCreated signature: TokenCreated(address,address,uint8,uint96,uint256)
type LogTokenCreated struct {
	Owner        common.Address // token owner
	Token        common.Address // token address
	TokenType    uint8
	TokenVersion *big.Int
	Level        *big.Int
}

type AppCtx struct {
	// db connection
}

func logEventHandleFunc(appctx *AppCtx, c *Contract) LogHandleFunc {
	return func(ctx context.Context, l *ParsedLog) error {
		slog.Info("eventInfo", slog.Any("event", l.Event))

		// use appctx do something
		_ = appctx

		// handle contract's different event by name
		switch l.Event.Name {
		case "TokenCreated":
			var data LogTokenCreated
			if err := c.Abi.UnpackIntoInterface(&data, l.Event.Name, l.Event.Data); err != nil {
				return errors.Wrap(err, "fail to unpack log")
			}

			// handle indexed topic
			data.Owner = HashToAddress(l.Event.IndexedParams[0])
			data.Token = HashToAddress(l.Event.IndexedParams[1])
			slog.Info("TokenCreated event", slog.Any("event", data))
		default:
			// do nothing
		}
		return nil
	}
}

func NewTokenFactoryEventListener() (*EventListener, error) {
	config := ChainConfig{
		ChainId:   11155111,
		ChainName: "ethereum-sepolia",
		URL:       "https://eth-sepolia.g.alchemy.com/v2/gOeoBV9mlFL1pWj7qbKEdlB6pXTfNum6",
	}
	client, err := ethclient.Dial(config.URL)
	if err != nil {
		return nil, err 
	}
	tokenFactorAddress := "0x822935C2240E6A0b5C96E3eA355446a83ed12C03"
	abiStr := `[{"inputs":[{"internalType":"address","name":"factoryManager_","type":"address"},{"internalType":"address","name":"implementation_","type":"address"},{"internalType":"address","name":"feeTo_","type":"address"},{"internalType":"uint256","name":"maxFee_","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"InsufficientFee","type":"error"},{"inputs":[{"internalType":"address","name":"implementation","type":"address"}],"name":"InvalidFactoryManager","type":"error"},{"inputs":[{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"InvalidFee","type":"error"},{"inputs":[{"internalType":"address","name":"receiver","type":"address"}],"name":"InvalidFeeReceiver","type":"error"},{"inputs":[{"internalType":"address","name":"factoryManager","type":"address"}],"name":"InvalidImplementation","type":"error"},{"inputs":[{"internalType":"uint256","name":"level","type":"uint256"}],"name":"InvalidLevel","type":"error"},{"inputs":[{"internalType":"uint256","name":"maxFee","type":"uint256"}],"name":"InvalidMaxFee","type":"error"},{"inputs":[],"name":"OnlyOwner","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"newFeeTo","type":"address"}],"name":"FeeToUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"level","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"newFee","type":"uint256"}],"name":"FeeUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256[]","name":"newLevels","type":"uint256[]"}],"name":"LevelsUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"token","type":"address"},{"indexed":false,"internalType":"uint8","name":"tokenType","type":"uint8"},{"indexed":false,"internalType":"uint96","name":"tokenVersion","type":"uint96"},{"indexed":false,"internalType":"uint256","name":"level","type":"uint256"}],"name":"TokenCreated","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"string","name":"description","type":"string"},{"internalType":"string","name":"logoLink","type":"string"},{"internalType":"string","name":"twitterLink","type":"string"},{"internalType":"string","name":"telegramLink","type":"string"},{"internalType":"string","name":"discordLink","type":"string"},{"internalType":"string","name":"websiteLink","type":"string"}],"indexed":false,"internalType":"struct TokenMetaData","name":"tokenMetaData","type":"tuple"}],"name":"TokenMetaDataUpdated","type":"event"},{"inputs":[],"name":"FACTORY_MANAGER","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"MAX_FEE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"level","type":"uint256"},{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"},{"internalType":"uint8","name":"decimals","type":"uint8"},{"internalType":"uint256","name":"totalSupply","type":"uint256"},{"internalType":"string","name":"description","type":"string"},{"internalType":"string","name":"logoLink","type":"string"},{"internalType":"string","name":"twitterLink","type":"string"},{"internalType":"string","name":"telegramLink","type":"string"},{"internalType":"string","name":"discordLink","type":"string"},{"internalType":"string","name":"websiteLink","type":"string"}],"internalType":"struct TokenInitializeParams","name":"tokenInitializeParams","type":"tuple"}],"name":"create","outputs":[{"internalType":"address","name":"token","type":"address"}],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"feeTo","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"fees","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getLevels","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"implementation","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"implementationVersion","outputs":[{"internalType":"uint96","name":"","type":"uint96"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"level","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"setFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"feeTo_","type":"address"}],"name":"setFeeTo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"implementation_","type":"address"}],"name":"setImplementation","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256[]","name":"_levels","type":"uint256[]"}],"name":"setLevels","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"level","type":"uint256"},{"internalType":"address","name":"token","type":"address"},{"components":[{"internalType":"string","name":"description","type":"string"},{"internalType":"string","name":"logoLink","type":"string"},{"internalType":"string","name":"twitterLink","type":"string"},{"internalType":"string","name":"telegramLink","type":"string"},{"internalType":"string","name":"discordLink","type":"string"},{"internalType":"string","name":"websiteLink","type":"string"}],"internalType":"struct TokenMetaData","name":"tokenMetaData_","type":"tuple"}],"name":"updateTokenMetaData","outputs":[],"stateMutability":"payable","type":"function"}]`
	blockNumber := big.NewInt(6641775)
	step := big.NewInt(100)
	
	tokenFactory, err := NewContract(tokenFactorAddress, abiStr, blockNumber, step)
	if err != nil {
		panic(err)
	}
	tokenFactory.SetLogHandler(logEventHandleFunc(&AppCtx{}, tokenFactory))

	el, err := New(
		config,
		WithClient(client),
		WithContract(tokenFactory),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to new an EventListener object")
	}

	return el, nil
}

func TestEventListen(t *testing.T) {
	logger.Init(slog.LevelInfo, false)
	el, err := NewTokenFactoryEventListener()
	if err != nil {
		t.Fatal(err)
	}
	el.Start()
}
