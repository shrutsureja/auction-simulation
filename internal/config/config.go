package config

import "time"

type Config struct {
	NumAuctions     int
	NumBidders      int
	AuctionDuration time.Duration
}

func DefaultConfig() Config {
	return Config{
		NumAuctions:     40,
		NumBidders:      100,
		AuctionDuration: 500 * time.Millisecond,
	}
}

func LoadConfig() Config {
	// For keeping this simple, loading the default config.
	// other wise we can load it from env, args or config files...
	return DefaultConfig()
}
