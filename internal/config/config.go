package config

import "time"

type Config struct {
	NumAuctions           int
	MaxConcurrentAuctions int
	NumBidders            int
	AuctionDuration       time.Duration
	MaxCPU                int   // GOMAXPROCS limit
	MaxMemoryBytes        int64 // GOMEMLIMIT in bytes
}

func DefaultConfig() Config {
	return Config{
		NumAuctions:           40,
		MaxConcurrentAuctions: 40,
		NumBidders:            100,
		AuctionDuration:       500 * time.Millisecond,
		MaxCPU:                2,                 // standardize to 2 vCPUs
		MaxMemoryBytes:        512 * 1024 * 1024, // standardize to 512MB RAM
	}
}
