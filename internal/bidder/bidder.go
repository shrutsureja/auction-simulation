package bidder

import (
	"auction-simulation/internal/types"
	"context"
	"math/rand/v2"
	"time"
)

type Bidder struct {
	ID string
}

// place bid will get a channel where bidder can pass the bid
func (b *Bidder) PlaceBid(ctx context.Context, auctionTimeout time.Duration) types.BidResponse {
	// random latency between 0 and auctionTimeout+200ms (some bidders may respond after timeout)
	maxMs := int(auctionTimeout.Milliseconds()) + 200
	randLatency := time.Duration(rand.IntN(maxMs)) * time.Millisecond

	// this is to simulate the bid amount
	bidAmount := rand.Float64() * 10

	select {
	case <-ctx.Done():
		// context timeout, do not place the bid
		return types.BidResponse{
			BidderId: b.ID,
			Amount:   0,
		}
	case <-time.After(randLatency):
		// latency is over, place the bid
		// return the bid amount and the latency
		return types.BidResponse{
			BidderId: b.ID,
			Amount:   bidAmount,
		}
	}

}
