package config

type BlockChain struct {
	Broker *BrokerConnectOptions
}

var BlockChainConfig = new(BlockChain)
