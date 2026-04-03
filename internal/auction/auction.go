package auction

import (
	"auction-simulation/internal/bidder"
	"auction-simulation/internal/config"
	"auction-simulation/internal/types"
	"context"
	"log/slog"
	"time"
)

type Auction struct {
	ID         string
	Request    types.BidRequest
	Config     *config.Config
	BidderPool *bidder.BidderPool
}

// when a auction runs, it basically sends requests to bidders and wait for their bids
// before the timeout is reached. Once the timeout is reached, it will have to
// decide the winner and then end the auction
//
// in real it is like website loading and publisher asking the adExchange for adds
func (a *Auction) StartAuction() types.AuctionResult {
	ctx, cancel := context.WithTimeout(context.Background(), a.Config.AuctionDuration)
	defer cancel()

	startTime := time.Now()

	// place bids in parallel — each bidder receives the full BidRequest with attributes
	bidResponses := make(chan types.BidResponse, a.Config.NumBidders)
	for _, b := range a.BidderPool.Bidders {
		go func(bidder bidder.Bidder) {
			bidResponses <- bidder.PlaceBid(ctx, a.Request, a.Config.AuctionDuration)
		}(b)
	}

	var winningBid *types.BidResponse
	var bidsReceived []types.BidResponse
	for {
		select {
		case bidresponse, ok := <-bidResponses:
			// we will keep receiving bids until the auction duration is over or context is done
			if ok {
				if bidresponse.Amount != 0 {
					bidsReceived = append(bidsReceived, bidresponse)
					if winningBid == nil || bidresponse.Amount > winningBid.Amount {
						winningBid = &bidresponse
					}
				}
			}
		case <-ctx.Done():
			// context timeout, do not wait for more bids
			slog.Debug("auction closed", "auction_id", a.ID, "reason", "timeout", "bids_collected", len(bidsReceived))
			if winningBid != nil {
				slog.Debug("winner declared", "auction_id", a.ID, "bidder", winningBid.BidderID, "amount", winningBid.Amount)
			} else {
				slog.Warn("no valid bids", "auction_id", a.ID)
			}

			endTime := time.Now()

			return types.AuctionResult{
				AuctionID:    a.ID,
				Attributes:   a.Request.Attributes,
				TotalBidders: len(a.BidderPool.Bidders),
				BidsReceived: bidsReceived,
				BidWinner:    winningBid,
				Timeout:      a.Config.AuctionDuration,
				Duration:     endTime.Sub(startTime), // this tell us that auction got completed with the time frame
				StartTime:    startTime,
				EndTime:      endTime,
			}
		}
	}
}
