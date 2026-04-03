package main

import (
	"auction-simulation/internal/auction"
	"auction-simulation/internal/bidder"
	"auction-simulation/internal/config"
	"auction-simulation/internal/output"
	"auction-simulation/internal/resource"
	"auction-simulation/internal/types"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	// logger
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.DefaultConfig()

	// Resource standardization: limit CPU and memory
	// As mentioned in mail
	runtime.GOMAXPROCS(cfg.MaxCPU)
	debug.SetMemoryLimit(cfg.MaxMemoryBytes)
	snap := resource.TakeSnapshot()
	slog.Info("system resources", "snapshot", snap.String())

	// creating the bidder pool
	pool := bidder.NewBidderPool(cfg.NumBidders)

	engine := auction.Engine{
		Config:     &cfg,
		BidderPool: pool,
	}

	// runs the simulation
	results, duration, resBefore, resAfter := engine.RunAll()
	slog.Info("all auctions completed", "total_duration", duration.String(), "total_auctions", len(results))

	// building summary and generating output files
	generateOutputFiles(results, buildSummary(cfg, results, duration, resBefore, resAfter))

	for _, r := range results {
		if r.BidWinner != nil {
			slog.Info("result",
				"auction", r.AuctionId,
				"winner", r.BidWinner.BidderId,
				"amount", r.BidWinner.Amount,
				"bids_received", len(r.BidsReceived),
				"duration", r.Duration.String(),
			)
		} else {
			slog.Info("result",
				"auction", r.AuctionId,
				"winner", "none",
				"bids_received", len(r.BidsReceived),
				"duration", r.Duration.String(),
			)
		}
	}
}

func buildSummary(cfg config.Config, results []types.AuctionResult, duration time.Duration, before, after resource.Snapshot) types.SimulationSummary {
	summary := types.SimulationSummary{
		NumAuctions:           cfg.NumAuctions,
		MaxConcurrentAuctions: cfg.MaxConcurrentAuctions,
		NumBidders:            cfg.NumBidders,
		AuctionTimeoutMs:      cfg.AuctionDuration.Milliseconds(),
		MaxCPU:                cfg.MaxCPU,
		MaxMemoryMB:           cfg.MaxMemoryBytes / 1024 / 1024,
		TotalDurationMs:       duration.Milliseconds(),
		TotalDuration:         duration.String(),
		ResourceBefore: types.ResourceSnapshot{
			AllocMB: before.AllocMB, TotalAllocMB: before.TotalAllocMB,
			SysMB: before.SysMB, NumGC: before.NumGC,
			NumGoroutine: before.NumGoroutine, NumCPU: before.NumCPU,
		},
		ResourceAfter: types.ResourceSnapshot{
			AllocMB: after.AllocMB, TotalAllocMB: after.TotalAllocMB,
			SysMB: after.SysMB, NumGC: after.NumGC,
			NumGoroutine: after.NumGoroutine, NumCPU: after.NumCPU,
		},
		MemoryDeltaMB: after.AllocMB - before.AllocMB,
	}

	// find the earliest auction start as the reference point (T=0)
	simStart := results[0].StartTime
	for _, r := range results {
		if r.StartTime.Before(simStart) {
			simStart = r.StartTime
		}
	}

	var totalBids int
	var totalDurationMs float64
	var winnerCount int
	var winSum, maxWin, minWin float64
	minWin = 999

	for _, r := range results {
		totalBids += len(r.BidsReceived)
		durationMs := float64(r.Duration.Microseconds()) / 1000.0
		totalDurationMs += durationMs

		as := types.AuctionSummary{
			AuctionId:    r.AuctionId,
			BidsReceived: len(r.BidsReceived),
			DurationMs:   durationMs,
			StartMs:      float64(r.StartTime.Sub(simStart).Microseconds()) / 1000.0,
			EndMs:        float64(r.EndTime.Sub(simStart).Microseconds()) / 1000.0,
		}
		if r.BidWinner != nil {
			winnerCount++
			winSum += r.BidWinner.Amount
			as.WinnerId = r.BidWinner.BidderId
			as.WinnerAmount = r.BidWinner.Amount
			if r.BidWinner.Amount > maxWin {
				maxWin = r.BidWinner.Amount
			}
			if r.BidWinner.Amount < minWin {
				minWin = r.BidWinner.Amount
			}
		}
		summary.Auctions = append(summary.Auctions, as)
	}

	summary.TotalBidsReceived = totalBids
	summary.AuctionsWithWinner = winnerCount
	summary.AuctionsNoWinner = len(results) - winnerCount
	if len(results) > 0 {
		summary.AvgBidsPerAuction = float64(totalBids) / float64(len(results))
		summary.AvgAuctionDurationMs = int64(totalDurationMs / float64(len(results)))
	}
	if winnerCount > 0 {
		summary.AvgWinningBid = winSum / float64(winnerCount)
		summary.MaxWinningBid = maxWin
		summary.MinWinningBid = minWin
	}

	return summary
}

func generateOutputFiles(results []types.AuctionResult, summary types.SimulationSummary) {
	// clean old output and write new files
	if err := output.CleanOutputDir("output"); err != nil {
		slog.Error("failed to clean output dir", "error", err)
	}
	if err := output.WriteResults("output", results); err != nil {
		slog.Error("failed to write output files", "error", err)
	}
	if err := output.WriteSummary("output", summary); err != nil {
		slog.Error("failed to write summary", "error", err)
	}
	if err := output.WriteDashboard("output"); err != nil {
		slog.Error("failed to write dashboard", "error", err)
	}
}
