package bidder

import "fmt"

type BidderPool struct {
	Bidders []Bidder
}

func NewBidderPool(numBidders int) *BidderPool {
	bidders := make([]Bidder, numBidders)
	for i := 0; i < numBidders; i++ {
		bidders[i] = Bidder{ID: fmt.Sprintf("bidder_%d", i+1)}
	}
	return &BidderPool{Bidders: bidders}
}
