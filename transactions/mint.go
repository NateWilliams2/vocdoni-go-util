package transactions

import (
	"time"
	"util/client"

	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/dvote/log"
)

func MintTokens(client *client.Client, treasurer, to *ethereum.SignKeys, amount uint64) {
	// If account does not exist, create account
	if _, _, _, err := client.GetAccount(to.Address().Bytes()); err != nil {
		log.Infof("to account does not yet exist: %s", err.Error())
		if err := client.SetAccountInfo(to, nil, "ipfs://QmTfquDp9puFXSpSyGyieGCxKUvFEr8UHXof421R5pTSHh", 0); err != nil {
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
