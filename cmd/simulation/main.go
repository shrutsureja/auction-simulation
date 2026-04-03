package main

import (
	"auction-simulation/internal/auction"
	"auction-simulation/internal/bidder"
	"auction-simulation/internal/config"
	"auction-simulation/internal/output"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
)

func main() {
	// logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.DefaultConfig()

	// Resource standardization: limit CPU and memory
	// As mentioned in mail
	runtime.GOMAXPROCS(cfg.MaxCPU)
	debug.SetMemoryLimit(cfg.MaxMemoryBytes)

	// creating the bidder pool
	pool := bidder.NewBidderPool(cfg.NumBidders)

	// setting up engine
	engine := auction.Engine{
		Config:     &cfg,
		BidderPool: pool,
	}

	// running the simulation
	results, duration, resBefore, resAfter := engine.RunAll()
	slog.Info("all auctions completed", "total_duration", duration.String(), "total_auctions", len(results))

	// building summary and generating output files
	output.GenerateAll("output", cfg, results, duration, resBefore, resAfter)

	for _, r := range results {
		if r.BidWinner != nil {
			slog.Info("result",
				"auction", r.AuctionID,
				"winner", r.BidWinner.BidderID,
				"amount", r.BidWinner.Amount,
				"bids_received", len(r.BidsReceived),
				"duration", r.Duration.String(),
			)
		} else {
			slog.Info("result",
				"auction", r.AuctionID,
				"winner", "none",
				"bids_received", len(r.BidsReceived),
				"duration", r.Duration.String(),
			)
		}
	}
}
