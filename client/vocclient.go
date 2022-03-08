package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"go.vocdoni.io/dvote/api"
	"go.vocdoni.io/dvote/client"
	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	gw         *client.Client
	signingKey *ethereum.SignKeys
}

// New initializes a new gatewayPool with the gatewayUrls, in order of health
// returns the new Client
func New(gatewayUrl string, signingKey *ethereum.SignKeys) (*Client, error) {
	gw, err := DiscoverGateway(gatewayUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		gw:         gw,
		signingKey: signingKey,
	}, nil
}

// ActiveEndpoint returns the address of the current active endpoint, if one exists
func (c *Client) ActiveEndpoint() string {
	if c.gw == nil {
		return ""
	}
	return c.gw.Addr
}

func (c *Client) request(req api.APIrequest,
	signer *ethereum.SignKeys) (*api.APIresponse, error) {
	resp, err := c.gw.Request(req, signer)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf(resp.Message)
	}
	return resp, nil
}

// MintTokens mints tokens to a given address
func (c *Client) MintTokens(treasurer *ethereum.SignKeys, to []byte, amount uint64, nonce uint32) error {
	req := api.APIrequest{Method: "submitRawTx"}
	p := &models.MintTokensTx{
		Txtype: models.TxType_MINT_TOKENS,
		Nonce:  nonce,
		To:     to,
		Value:  amount,
	}
	var err error
	stx := &models.SignedTx{}
	stx.Tx, err = proto.Marshal(&models.Tx{Payload: &models.Tx_MintTokens{MintTokens: p}})
	if err != nil {
		return err
	}
	if stx.Signature, err = treasurer.SignVocdoniTx(stx.Tx); err != nil {
		return err
	}
	if req.Payload, err = proto.Marshal(stx); err != nil {
		return err
	}
	resp, err := c.request(req, c.signingKey)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

// SetAccountInfo submits a transaction to set an account with the given
//  ethereum wallet address and metadata URI on the vochain
func (c *Client) SetAccountInfo(signer *ethereum.SignKeys,
	faucet *ethereum.SignKeys, uri string, nonce uint32) error {
	req := api.APIrequest{Method: "submitRawTx"}
	tx := models.Tx_SetAccountInfo{SetAccountInfo: &models.SetAccountInfoTx{
		Txtype:  models.TxType_SET_ACCOUNT_INFO,
		Nonce:   nonce,
		InfoURI: uri,
	}}
	var err error
	stx := new(models.SignedTx)
	stx.Tx, err = proto.Marshal(&models.Tx{Payload: &tx})
	if err != nil {
		return fmt.Errorf("could not marshal set account info tx")
	}
	stx.Signature, err = signer.SignVocdoniTx(stx.Tx)
	if err != nil {
		return fmt.Errorf("could not sign account transaction: %v", err)
	}
	if req.Payload, err = proto.Marshal(stx); err != nil {
		return err
	}
	resp, err := c.request(req, c.signingKey)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

// GetAccount returns the metadata URI, token balance, and nonce for the
//  given account ID on the vochain
func (c *Client) GetAccount(entityId []byte) (string, uint64, uint32, error) {
	req := api.APIrequest{Method: "getAccount", EntityId: entityId}
	resp, err := c.request(req, c.signingKey)
	if err != nil {
		return "", 0, 0, err
	}
	if !resp.Ok {
		return "", 0, 0, fmt.Errorf("could not get account: %s", resp.Message)
	}
	if resp.Balance == nil {
		resp.Balance = new(uint64)
	}
	if resp.Nonce == nil {
		resp.Nonce = new(uint32)
	}
	if resp.InfoURI == "" {
		return "", 0, 0, fmt.Errorf("account info URI not yet set")
	}
	return resp.InfoURI, *resp.Balance, *resp.Nonce, nil
}

// GetTreasurer returns information about the treasurer
func (c *Client) GetTreasurer(signer *ethereum.SignKeys) (*models.Treasurer, error) {
	req := api.APIrequest{Method: "getTreasurer"}
	resp, err := c.request(req, signer)
	treasurer := &models.Treasurer{}
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot not get treasurer: %s", resp.Message)
	}
	if resp.EntityID != "" {
		treasurer.Address = common.HexToAddress(resp.EntityID).Bytes()
	}
	if resp.Nonce != nil {
		treasurer.Nonce = *resp.Nonce
	}
	return treasurer, nil
}
