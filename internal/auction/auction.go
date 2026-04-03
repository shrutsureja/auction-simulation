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
	Id         string
	Request    types.BidRequest
	Config     config.Config
	BidderPool *bidder.BidderPool
}

// when a auction runs, it basically sends requests to bidders and wait for their bids
// before the timeout is reached. Once the timeout is reached, it will have to
// decide the winner and then end the auction
//
// in real it is like website loading and publisher asking the adexchange for add
func (a *Auction) StartAuction() types.AuctionResult {
	ctx, cancel := context.WithTimeout(context.Background(), a.Config.AuctionDuration)
	defer cancel()

	startTime := time.Now()

	// Is place the bids in parallel and wait for the context to timeout
	bidResponses := make(chan types.BidResponse, len(a.BidderPool.Bidders))
	for _, b := range a.BidderPool.Bidders {
		go func(bidder bidder.Bidder) {
			bidResponses <- bidder.PlaceBid(ctx, a.Config.AuctionDuration)
		}(b)
	}

	var winningBid *types.BidResponse
	var bidsReceived []types.BidResponse
	for {
		select {
		case bidresponse, ok := <-bidResponses:
			// we will keep receiving bids until the auction duration is over or context is done
			if ok {
				bidsReceived = append(bidsReceived, bidresponse)
				if bidresponse.Amount != 0 {
					if winningBid == nil || bidresponse.Amount > winningBid.Amount {
						winningBid = &bidresponse
					}
				}
			}
		case <-ctx.Done():
			// context timeout, do not wait for more bids
			slog.Debug("auction closed", "auction_id", a.Id, "reason", "timeout", "bids_collected", len(bidsReceived))
			if winningBid != nil {
				slog.Debug("winner declared", "auction_id", a.Id, "bidder", winningBid.BidderId, "amount", winningBid.Amount)
			} else {
				slog.Warn("no valid bids", "auction_id", a.Id)
			}

			endTime := time.Now()

			return types.AuctionResult{
				AuctionId:    a.Id,
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
