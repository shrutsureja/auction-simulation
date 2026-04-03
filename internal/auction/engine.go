package auction

import (
	"auction-simulation/internal/bidder"
	"auction-simulation/internal/config"
	"auction-simulation/internal/resource"
	"auction-simulation/internal/types"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"
)

type Engine struct {
	Config     *config.Config
	BidderPool *bidder.BidderPool
}

func (e *Engine) RunAll() ([]types.AuctionResult, time.Duration, resource.Snapshot, resource.Snapshot) {
	// start of the simulation
	start := time.Now()

	// capturing the initial resource memory
	before := resource.TakeSnapshot()
	slog.Info("engine starting", "auctions", e.Config.NumAuctions, "bidders", e.Config.NumBidders, "resources", before.String())

	resultsChan := make(chan types.AuctionResult, e.Config.NumAuctions)
	var wg sync.WaitGroup
	sem := make(chan struct{}, e.Config.MaxConcurrentAuctions) // semaphore to limit concurrent auctions

	for i := 0; i < e.Config.NumAuctions; i++ {
		wg.Add(1)
		go func(index int) {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()

			id := fmt.Sprintf("auction_%d", index+1)
			auction := Auction{
				Id:         id,
				Request:    types.BidRequest{AuctionId: id, Attributes: generateAttributes()},
				Config:     *e.Config,
				BidderPool: e.BidderPool,
			}
			resultsChan <- auction.StartAuction()
		}(i)
	}

	// capture peak resource usage while all auctions are running
	// time.Sleep(50 * time.Millisecond) // brief pause to let goroutines spin up
	// peak := resource.TakeSnapshot()
	// slog.Info("peak resources (mid-run)", "resources", peak.String())

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

	after := resource.TakeSnapshot()
	slog.Info("engine finished",
		"total_duration", end.Sub(start).String(),
		"resources", after.String(),
		"memory_delta_mb", fmt.Sprintf("%.2f", after.AllocMB-before.AllocMB),
	)

	return results, end.Sub(start), before, after
}

// generateAttributes creates 20 random attributes simulating an ad auction object.
func generateAttributes() map[string]string {
	pick := func(options []string) string {
		return options[rand.IntN(len(options))]
	}

	return map[string]string{
		"category":       pick([]string{"electronics", "fashion", "sports", "automotive", "food", "travel"}),
		"subcategory":    pick([]string{"phones", "laptops", "shoes", "watches", "bikes", "snacks"}),
		"region":         pick([]string{"us-east", "us-west", "eu-west", "eu-east", "ap-south", "ap-east"}),
		"device":         pick([]string{"mobile", "desktop", "tablet"}),
		"os":             pick([]string{"android", "ios", "windows", "macos", "linux"}),
		"browser":        pick([]string{"chrome", "firefox", "safari", "edge"}),
		"ad_format":      pick([]string{"banner", "video", "native", "interstitial"}),
		"ad_size":        pick([]string{"300x250", "728x90", "160x600", "320x50"}),
		"language":       pick([]string{"en", "es", "fr", "de", "hi", "zh"}),
		"gender":         pick([]string{"male", "female", "unknown"}),
		"age_group":      pick([]string{"18-24", "25-34", "35-44", "45-54", "55+"}),
		"time_of_day":    pick([]string{"morning", "afternoon", "evening", "night"}),
		"day_of_week":    pick([]string{"weekday", "weekend"}),
		"connection":     pick([]string{"wifi", "4g", "5g", "3g"}),
		"publisher":      pick([]string{"news_site", "social_app", "game_app", "video_platform", "blog"}),
		"content_rating": pick([]string{"G", "PG", "PG-13", "R"}),
		"user_intent":    pick([]string{"purchase", "research", "entertainment", "education"}),
		"session_depth":  pick([]string{"1", "2-3", "4-6", "7+"}),
		"viewability":    pick([]string{"high", "medium", "low"}),
		"inventory_type": pick([]string{"premium", "standard", "remnant"}),
	}
}
