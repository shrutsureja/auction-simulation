package main

import (
	"auction-simulation/internal/auction"
	"auction-simulation/internal/bidder"
	"auction-simulation/internal/config"
)

func main() {

	cfg := config.DefaultConfig()

	pool := bidder.NewBidderPool(cfg.NumBidders)

	engine := auction.Engine{
		Config:     &cfg,
		BidderPool: pool,
	}
	responses, duration := engine.RunAll()
	println("Auction completed in", duration.String())

	for _, response := range responses {
		println("Auction Result:")
		println("Auction ID:", response.AuctionId)
		println("Total Bidders:", response.TotalBidders)
		println("Bids Received:")
		for _, bid := range response.BidsReceived {
			println("- Bidder ID:", bid.BidderId, "Amount:", bid.Amount)
		}
		if response.BidWinner != nil {
			println("Bid Winner: Bidder ID:", response.BidWinner.BidderId, "Amount:", response.BidWinner.Amount)
		} else {
			println("No valid bids received.")
		}
	}
}
