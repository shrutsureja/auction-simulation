package report

import (
	"github.com/shrutsureja/auction-simulation/internal/config"
	"github.com/shrutsureja/auction-simulation/internal/resource"
	"github.com/shrutsureja/auction-simulation/internal/types"

	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"time"
)

// GenerateAll builds the summary and writes all output files (auction JSONs, summary, dashboard).
func GenerateAll(dir string, cfg config.Config, results []types.AuctionResult, duration time.Duration, before, after resource.Snapshot) {
	summary := buildSummary(cfg, results, duration, before, after)

	if err := CleanOutputDir(dir); err != nil {
		slog.Error("failed to clean output dir", "error", err)
	}
	if err := WriteResults(dir, results); err != nil {
		slog.Error("failed to write output files", "error", err)
	}
	if err := WriteSummary(dir, summary); err != nil {
		slog.Error("failed to write summary", "error", err)
	}
	if err := WriteDashboard(dir); err != nil {
		slog.Error("failed to write dashboard", "error", err)
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
		ResourceBefore:        before,
		ResourceAfter:         after,
		MemoryDeltaMB:         after.AllocMB - before.AllocMB,
	}

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
	minWin = math.MaxFloat64

	for _, r := range results {
		totalBids += len(r.BidsReceived)
		durationMs := float64(r.Duration.Microseconds()) / 1000.0
		totalDurationMs += durationMs

		as := types.AuctionSummary{
			AuctionID:    r.AuctionID,
			BidsReceived: len(r.BidsReceived),
			DurationMs:   durationMs,
			StartMs:      float64(r.StartTime.Sub(simStart).Microseconds()) / 1000.0,
			EndMs:        float64(r.EndTime.Sub(simStart).Microseconds()) / 1000.0,
		}
		if r.BidWinner != nil {
			winnerCount++
			winSum += r.BidWinner.Amount
			as.WinnerID = r.BidWinner.BidderID
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

// CleanOutputDir removes all files in the output directory before writing new results.
func CleanOutputDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("clean output dir: %w", err)
	}
	slog.Info("cleaned output directory", "dir", dir)
	return nil
}

// WriteSummary writes the simulation summary as a JSON file.
func WriteSummary(dir string, summary types.SimulationSummary) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal summary: %w", err)
	}
	filename := filepath.Join(dir, "summary.json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}
	slog.Info("summary written", "file", filename)
	return nil
}

// WriteResults writes each auction result as a separate JSON file in the given directory.
func WriteResults(dir string, results []types.AuctionResult) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	for _, r := range results {
		filename := filepath.Join(dir, r.AuctionID+".json")
		data, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			slog.Error("failed to marshal result", "auction_id", r.AuctionID, "error", err)
			continue
		}
		if err := os.WriteFile(filename, data, 0644); err != nil {
			slog.Error("failed to write result file", "auction_id", r.AuctionID, "error", err)
			continue
		}
		slog.Debug("output written", "file", filename)
	}

	slog.Info("output files written", "dir", dir, "count", len(results))
	return nil
}

// WriteDashboard writes an HTML dashboard file that visualizes summary.json.
func WriteDashboard(dir string) error {
	filename := filepath.Join(dir, "dashboard.html")
	if err := os.WriteFile(filename, []byte(dashboardHTML), 0644); err != nil {
		return fmt.Errorf("write dashboard: %w", err)
	}
	slog.Info("dashboard written", "file", filename)
	return nil
}
