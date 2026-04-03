package types

import (
	"time"

	"github.com/shrutsureja/auction-simulation/internal/resource"
)

// request that auction sends to bidders
type BidRequest struct {
	AuctionID  string
	Attributes map[string]string
}

// response the bidder sends back
type BidResponse struct {
	BidderID string  `json:"bidder_id"`
	Amount   float64 `json:"amount"`
}

// storing the result of an auction for reporting
type AuctionResult struct {
	AuctionID    string            `json:"auction_id"`
	Attributes   map[string]string `json:"attributes"`
	TotalBidders int               `json:"total_bidders"`
	BidsReceived []BidResponse     `json:"bids_received"`
	BidWinner    *BidResponse      `json:"bid_winner,omitempty"`
	Timeout      time.Duration     `json:"timeout"`
	Duration     time.Duration     `json:"duration"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
}

// summary of the entire simulation run for reporting
type SimulationSummary struct {
	// Config
	NumAuctions           int               `json:"num_auctions"`
	MaxConcurrentAuctions int               `json:"max_concurrent_auctions"`
	NumBidders            int               `json:"num_bidders"`
	AuctionTimeoutMs      int64             `json:"auction_timeout_ms"`
	MaxCPU                int               `json:"max_cpu"`
	MaxMemoryMB           int64             `json:"max_memory_mb"`
	TotalDurationMs       int64             `json:"total_duration_ms"`
	TotalDuration         string            `json:"total_duration"`
	ResourceBefore        resource.Snapshot `json:"resource_before"`
	ResourceAfter         resource.Snapshot `json:"resource_after"`
	MemoryDeltaMB         float64           `json:"memory_delta_mb"`

	// Auction stats
	AuctionsWithWinner   int     `json:"auctions_with_winner"`
	AuctionsNoWinner     int     `json:"auctions_no_winner"`
	TotalBidsReceived    int     `json:"total_bids_received"`
	AvgBidsPerAuction    float64 `json:"avg_bids_per_auction"`
	AvgWinningBid        float64 `json:"avg_winning_bid"`
	MaxWinningBid        float64 `json:"max_winning_bid"`
	MinWinningBid        float64 `json:"min_winning_bid"`
	AvgAuctionDurationMs int64   `json:"avg_auction_duration_ms"`

	// Per-auction summary
	Auctions []AuctionSummary `json:"auctions"`
}

// summary of each auction for reporting
type AuctionSummary struct {
	AuctionID    string  `json:"auction_id"`
	BidsReceived int     `json:"bids_received"`
	WinnerID     string  `json:"winner_id,omitempty"`
	WinnerAmount float64 `json:"winner_amount,omitempty"`
	DurationMs   float64 `json:"duration_ms"`
	StartMs      float64 `json:"start_ms"`
	EndMs        float64 `json:"end_ms"`
}
