package types

import "time"

type BidRequest struct {
	AuctionId  string
	Attributes map[string]string
}

type BidResponse struct {
	BidderId string
	Amount   float64
}

type AuctionResult struct {
	AuctionId    string            `json:"auction_id"`
	Attributes   map[string]string `json:"attributes"`
	TotalBidders int               `json:"total_bidders"`
	BidsReceived []BidResponse     `json:"bids_received"`
	BidWinner    *BidResponse      `json:"bid_winner,omitempty"`
	Timeout      time.Duration     `json:"timeout"`
	Duration     time.Duration     `json:"duration"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
}
