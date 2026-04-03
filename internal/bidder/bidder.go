package bidder

import (
	"auction-simulation/internal/types"
	"context"
	"math/rand/v2"
	"time"
)

type Bidder struct {
	ID        string
	MaxBudget float64 // each bidder has a different max budget (1-15)
	BidChance float64 // probability that this bidder participates (0.3 to 1.0)
}

// PlaceBid simulates a bidder deciding whether and how much to bid.
// The bidder receives the full BidRequest (with 20 attributes) and uses
// attributes to influence the bid amount — mimicking real Prebid behavior.
func (b *Bidder) PlaceBid(ctx context.Context, req types.BidRequest, auctionTimeout time.Duration) types.BidResponse {
	// some bidders choose not to participate in certain auctions
	if rand.Float64() > b.BidChance {
		return types.BidResponse{BidderID: b.ID, Amount: 0}
	}

	// random latency between 0 and auctionTimeout+200ms (some bidders may respond after timeout)
	maxMs := int(auctionTimeout.Milliseconds()) + 200
	randLatency := time.Duration(rand.IntN(maxMs)) * time.Millisecond

	// base bid varies based on the bidder's budget
	bidAmount := rand.Float64() * b.MaxBudget

	// adjust bid based on auction attributes (simulates real bidder logic)
	bidAmount *= attributeMultiplier(req.Attributes)

	// cap bid at max budget
	if bidAmount > b.MaxBudget {
		bidAmount = b.MaxBudget
	}

	// use NewTimer instead of time.After to avoid timer leak —
	// if ctx fires first, we stop the timer immediately
	timer := time.NewTimer(randLatency)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return types.BidResponse{
			BidderID: b.ID,
			Amount:   0,
		}
	case <-timer.C:
		return types.BidResponse{
			BidderID: b.ID,
			Amount:   bidAmount,
		}
	}
}

// attributeMultiplier returns a multiplier (0.5–1.5) based on auction attributes.
// Premium inventory, high viewability, and purchase intent increase bid value.
func attributeMultiplier(attrs map[string]string) float64 {
	multiplier := 1.0

	switch attrs["inventory_type"] {
	case "premium":
		multiplier += 0.2
	case "remnant":
		multiplier -= 0.2
	}

	switch attrs["viewability"] {
	case "high":
		multiplier += 0.15
	case "low":
		multiplier -= 0.15
	}

	switch attrs["user_intent"] {
	case "purchase":
		multiplier += 0.15
	case "entertainment":
		multiplier -= 0.1
	}

	switch attrs["ad_format"] {
	case "video":
		multiplier += 0.1
	case "native":
		multiplier += 0.05
	}

	switch attrs["device"] {
	case "mobile":
		multiplier += 0.05
	case "desktop":
		multiplier += 0.1
	}

	return multiplier
}
