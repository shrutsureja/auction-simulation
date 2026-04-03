package auction

import (
	"auction-simulation/internal/bidder"
	"auction-simulation/internal/config"
	"auction-simulation/internal/types"
	"fmt"
	"sync"
	"time"
)

type Engine struct {
	Config     *config.Config
	BidderPool *bidder.BidderPool
}

func (e *Engine) RunAll() ([]types.AuctionResult, time.Duration) {
	start := time.Now()

	resultsChan := make(chan types.AuctionResult, e.Config.NumAuctions)
	var wg sync.WaitGroup

	for i := 0; i < e.Config.NumAuctions; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			id := fmt.Sprintf("auction_%d", index+1)
			auction := Auction{
				Id:         id,
				Request:    types.BidRequest{AuctionId: id, Attributes: map[string]string{"category": "electronics"}},
				Config:     *e.Config,
				BidderPool: e.BidderPool,
			}
			resultsChan <- auction.StartAuction()
		}(i)
	}

	// close channel once all auctions finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []types.AuctionResult
	for result := range resultsChan {
		results = append(results, result)
	}

	end := time.Now()
	return results, end.Sub(start)
}
