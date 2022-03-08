package transactions

import (
	"log"
	"time"
	"util/client"

	"go.vocdoni.io/dvote/crypto/ethereum"
)

func MintTokens(client *client.Client, treasurer, to *ethereum.SignKeys, amount uint64) {
	// If account does not exist, create account
	if _, _, _, err := client.GetAccount(to.Address().Bytes()); err != nil {
		if err := client.SetAccountInfo(to, nil, "faucetUri", 0); err != nil {
			log.Fatalf("could not mint tokens to %x: could not set account info: %s", to.Address().Bytes(), err.Error())
		}
		time.Sleep(time.Second * 20)
	}
	treasurerAcct, err := client.GetTreasurer(to)
	if err != nil {
		log.Fatalf("could not get treasurer: %s", err.Error())
	}
	if _, _, _, err := client.GetAccount(to.Address().Bytes()); err != nil {
		log.Fatalf("target account does not exist, could not be created: %s", err.Error())
	}
	if err := client.MintTokens(treasurer, to.Address().Bytes(), amount, treasurerAcct.Nonce); err != nil {
		log.Fatalf("could not mint tokens to %x: %s", to.Address().Bytes(), err.Error())
	}
}
