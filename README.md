# Auction Simulation

A concurrent real-time bidding (RTB) simulation in Go — modeled after how Prebid-style ad auctions work.

## How it runs

```mermaid
flowchart TD
    A[Start: 100 Bidders created\neach with random budget + participation chance] --> B

    B[Engine: launch 100 auctions\nmax 40 running at once via semaphore] --> C

    C[Each auction sends BidRequest\n20 attributes: device, format, viewability...] --> D

    D[100 Bidders bid in parallel goroutines\nadjust bid based on attributes + simulate latency] --> E

    E{500ms timeout fires} --> F
    E --> G

    F[Bids collected before timeout\n→ highest bid wins] --> H
    G[Bid arrives after timeout\n→ dropped] --> H

    H[Results written to output/]
```

## Concurrency model

```mermaid
sequenceDiagram
    participant Engine
    participant Sem as Semaphore (size 40)
    participant Auction
    participant Bidders as 100 Bidders (goroutines)

    Engine->>Sem: acquire slot
    Sem-->>Engine: slot granted
    Engine->>Auction: StartAuction()
    Auction->>Bidders: PlaceBid(ctx, BidRequest) × 100
    Note over Bidders: each bidder simulates latency,<br/>adjusts bid from attributes
    Bidders-->>Auction: BidResponse (or timeout)
    Auction-->>Engine: AuctionResult (winner + all bids)
    Engine->>Sem: release slot
```

## Bidder decision logic

Each bidder has a random budget ($1–$15) and participation rate (30–100%). It also adjusts its bid based on the auction attributes:

| Signals that raise the bid | Signals that lower it |
|---|---|
| Premium inventory (+20%) | Remnant inventory (-20%) |
| High viewability (+15%) | Low viewability (-15%) |
| Purchase intent (+15%) | Entertainment intent (-10%) |
| Video ad format (+10%) | |
| Desktop device (+10%) | |

## Run it

```bash
go run ./cmd/simulation/
```

Output lands in `output/`. Serve the dashboard:

```bash
cd output && python3 -m http.server 8080
# open http://localhost:8080/dashboard.html
```

The dashboard shows a live concurrency timeline — every auction as a bar, green = winner, red = no bids, hover for sub-ms timing.

## Config

All parameters are in `internal/config/config.go`:

| Parameter | Default |
|---|---|
| Auctions | 100 |
| Max concurrent | 40 |
| Bidders | 100 |
| Auction timeout | 500ms |
| CPU limit | 2 vCPUs |
| Memory limit | 512 MB |

## Project layout

```
cmd/simulation/main.go       # entrypoint
internal/
  auction/engine.go          # semaphore, runs all auctions
  auction/auction.go         # single auction, for/select bid collection
  bidder/bidder.go           # bid logic, attribute multiplier, timer
  bidder/pool.go             # creates pool with random budgets
  config/config.go           # all parameters
  resource/resource.go       # runtime.MemStats snapshots
  types/types.go             # shared types
  output/writer.go           # summary builder + file writer
  output/dashboard.go        # embedded HTML timeline
```

## References

- https://www.youtube.com/watch?v=ylhKJSrxutM
- https://www.youtube.com/watch?v=Cqki_mlQmkI
