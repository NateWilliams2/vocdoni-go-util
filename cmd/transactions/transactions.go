package main

import (
	"util/client"
	"util/transactions"

	flag "github.com/spf13/pflag"
	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/dvote/log"
)

func newTransactionsConfig() (*transactions.Config, error) {
	cfg := transactions.Config{}
	// flags
	cfg.LogLevel = flag.String("logLevel", "info", "Log level (debug, info, warn, error, fatal)")
	cfg.Operation = flag.String("operation", "", "Operation to perform (mintTokens)")
	cfg.Amount = flag.Uint64("amount", 10000000, "number of tokens to send")
	cfg.Faucet = flag.String("faucet", "", "faucet")
	cfg.Treasurer = flag.String("treasurer", "", "treasurer")
	cfg.GatewayUrl = flag.String("gatewayUrl",
		"https://gw1.vocdoni.net", "url to use as gateway api endpoint")
	flag.CommandLine.SortFlags = false

	// parse flags
	flag.Parse()

	return &cfg, nil
}

func main() {
	cfg, err := newTransactionsConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Init(*cfg.LogLevel, "stdout")

	// Faucet
	faucet := ethereum.NewSignKeys()
	if *cfg.Faucet != "" {
		if err := faucet.AddHexKey(*cfg.Faucet); err != nil {
			log.Fatal(err)
		}
		pub, _ := faucet.HexString()
		log.Infof("faucet public key: %s, address %s", pub, faucet.Address().String())
	}
	// Treasurer
	treasurer := ethereum.NewSignKeys()
	if *cfg.Treasurer != "" {
		if err := treasurer.AddHexKey(*cfg.Treasurer); err != nil {
			log.Fatal(err)
		}
		pub, _ := treasurer.HexString()
		log.Infof("treasurer public key: %s, address %s", pub, treasurer.Address().String())
	}

	client, err := client.New(*cfg.GatewayUrl, treasurer)
	if err != nil {
		log.Fatal(err)
	}

	transactions.MintTokens(client, treasurer, faucet, *cfg.Amount)
}
