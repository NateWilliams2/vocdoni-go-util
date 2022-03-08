package transactions

type Config struct {
	GatewayUrl *string
	LogLevel   *string
	Faucet     *string
	Treasurer  *string
	Operation  *string
	Amount     *uint64
}
