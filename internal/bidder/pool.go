package bidder

import (
	"fmt"
	"math/rand/v2"
)

type BidderPool struct {
	Bidders []Bidder
}

func NewBidderPool(numBidders int) *BidderPool {
	bidders := make([]Bidder, numBidders)
	for i := 0; i < numBidders; i++ {
		bidders[i] = Bidder{
			ID:        fmt.Sprintf("bidder_%d", i+1),
			MaxBudget: 1 + rand.Float64()*14,    // budget between 1 and 15
			BidChance: 0.3 + rand.Float64()*0.7, // participation chance between 30% and 100%
		}
	}
	return &BidderPool{Bidders: bidders}
}
