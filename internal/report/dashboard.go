package report

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Auction Simulation - Concurrency Timeline</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Courier New', monospace; background: #0f172a; color: #e2e8f0; padding: 24px; }
        h1 { text-align: center; margin-bottom: 4px; font-size: 1.5rem; color: #38bdf8; }
        .subtitle { text-align: center; color: #64748b; margin-bottom: 20px; font-size: 0.85rem; }

        .stats-bar { display: flex; gap: 24px; justify-content: center; flex-wrap: wrap; margin-bottom: 24px; padding: 12px; background: #1e293b; border-radius: 8px; border: 1px solid #334155; }
        .stat { text-align: center; }
        .stat .val { font-size: 1.3rem; font-weight: 700; color: #38bdf8; }
        .stat .lbl { font-size: 0.65rem; text-transform: uppercase; color: #64748b; letter-spacing: 1px; }

        .timeline-container { background: #1e293b; border-radius: 8px; border: 1px solid #334155; padding: 20px; overflow-x: auto; }
        .timeline-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
        .timeline-header h2 { font-size: 1rem; color: #cbd5e1; }
        .legend { display: flex; gap: 16px; font-size: 0.75rem; color: #94a3b8; }
        .legend span { display: flex; align-items: center; gap: 4px; }
        .legend .dot { width: 10px; height: 10px; border-radius: 2px; }

        .timeline { position: relative; min-width: 800px; }
        .time-axis { position: relative; height: 24px; border-bottom: 1px solid #475569; margin-bottom: 4px; }
        .time-tick { position: absolute; top: 0; font-size: 0.65rem; color: #94a3b8; transform: translateX(-50%); }
        .time-tick::after { content: ''; position: absolute; left: 50%; bottom: -4px; width: 1px; height: 6px; background: #475569; }

        .auction-row { display: flex; align-items: center; height: 22px; margin: 2px 0; }
        .auction-label { width: 80px; flex-shrink: 0; font-size: 0.7rem; color: #94a3b8; text-align: right; padding-right: 8px; }
        .auction-track { position: relative; flex: 1; height: 100%; }
        .auction-bar { position: absolute; height: 16px; top: 3px; border-radius: 3px; cursor: pointer; transition: opacity 0.15s; min-width: 2px; }
        .auction-bar:hover { opacity: 0.85; }
        .tooltip { display: none; position: absolute; bottom: 24px; left: 50%; transform: translateX(-50%); background: #0f172a; border: 1px solid #475569; border-radius: 6px; padding: 8px 12px; font-size: 0.7rem; white-space: nowrap; z-index: 10; color: #e2e8f0; pointer-events: none; }
        .auction-bar:hover .tooltip { display: block; }

        .resource-section { margin-top: 20px; background: #1e293b; border-radius: 8px; border: 1px solid #334155; padding: 16px; }
        .resource-section h2 { font-size: 1rem; color: #cbd5e1; margin-bottom: 10px; }
        .res-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 8px; }
        .res-item { display: flex; justify-content: space-between; padding: 6px 10px; background: #0f172a; border-radius: 4px; font-size: 0.75rem; }
        .res-item .res-label { color: #94a3b8; }
        .res-item .res-val { color: #22d3ee; font-weight: 600; }

        #loading { text-align: center; padding: 60px; color: #94a3b8; }
    </style>
</head>
<body>
    <h1>Auction Concurrency Timeline</h1>
    <p class="subtitle">Each bar shows when an auction was actively collecting bids (start to end)</p>
    <div id="loading">Loading summary.json...</div>
    <div id="app" style="display:none;"></div>

    <script>
    async function init() {
        let d;
        try {
            const r = await fetch('./summary.json');
            d = await r.json();
        } catch(e) {
            document.getElementById('loading').textContent = 'Failed to load summary.json';
            return;
        }
        document.getElementById('loading').style.display = 'none';
        document.getElementById('app').style.display = 'block';

        const auctions = d.auctions.sort((a,b) => a.start_ms - b.start_ms);
        const maxEnd = Math.max(...auctions.map(a => a.end_ms));

        // stats bar
        let html = '<div class="stats-bar">';
        const stats = [
            ['Auctions', d.num_auctions],
            ['Concurrent', d.max_concurrent_auctions],
            ['Bidders', d.num_bidders],
            ['Timeout', d.auction_timeout_ms + 'ms'],
            ['Total', d.total_duration],
            ['Avg Bids', d.avg_bids_per_auction.toFixed(1)],
            ['Min Bid', '$' + d.min_winning_bid.toFixed(2)],
            ['Avg Winner', '$' + d.avg_winning_bid.toFixed(2)],
            ['Max Bid', '$' + d.max_winning_bid.toFixed(2)],
            ['Mem Delta', d.memory_delta_mb.toFixed(2) + 'MB'],
        ];
        stats.forEach(([lbl, val]) => {
            html += '<div class="stat"><div class="val">' + val + '</div><div class="lbl">' + lbl + '</div></div>';
        });
        html += '</div>';

        // timeline
        html += '<div class="timeline-container">';
        html += '<div class="timeline-header"><h2>Auction Execution Timeline</h2>';
        html += '<div class="legend"><span><div class="dot" style="background:#22c55e"></div> Has winner</span>';
        html += '<span><div class="dot" style="background:#ef4444"></div> No winner</span></div></div>';

        html += '<div class="timeline">';

        // time axis
        html += '<div class="time-axis">';
        const tickCount = 10;
        for (let i = 0; i <= tickCount; i++) {
            const ms = (maxEnd * i / tickCount);
            const pct = (i / tickCount * 100);
            html += '<span class="time-tick" style="left:calc(80px + (100% - 80px) * ' + pct/100 + ')">' + ms.toFixed(1) + 'ms</span>';
        }
        html += '</div>';

        // auction bars
        auctions.forEach(a => {
            const leftPct = (a.start_ms / maxEnd * 100);
            const widthPct = ((a.end_ms - a.start_ms) / maxEnd * 100);
            const color = a.winner_id ? '#22c55e' : '#ef4444';
            const num = a.auction_id.replace('auction_', '');

            html += '<div class="auction-row">';
            html += '<div class="auction-label">#' + num + '</div>';
            html += '<div class="auction-track">';
            html += '<div class="auction-bar" style="left:' + leftPct + '%;width:' + widthPct + '%;background:' + color + '">';
            html += '<div class="tooltip">';
            html += '<strong>' + a.auction_id + '</strong><br>';
            html += 'Start: ' + a.start_ms.toFixed(3) + 'ms<br>';
            html += 'End: ' + a.end_ms.toFixed(3) + 'ms<br>';
            html += 'Duration: ' + a.duration_ms.toFixed(3) + 'ms<br>';
            html += 'Bids: ' + a.bids_received + '<br>';
            if (a.winner_id) {
                html += 'Winner: ' + a.winner_id + ' ($' + a.winner_amount.toFixed(4) + ')';
            } else {
                html += 'Winner: none';
            }
            html += '</div></div></div></div>';
        });

        html += '</div></div>';

        // resource section
        html += '<div class="resource-section"><h2>Resource Usage</h2><div class="res-grid">';
        const res = [
            ['Heap Before', d.resource_before.alloc_mb.toFixed(2) + ' MB'],
            ['Heap After', d.resource_after.alloc_mb.toFixed(2) + ' MB'],
            ['Heap Delta', d.memory_delta_mb.toFixed(2) + ' MB'],
            ['Total Alloc', d.resource_after.total_alloc_mb.toFixed(2) + ' MB'],
            ['Sys Memory', d.resource_after.sys_mb.toFixed(2) + ' MB'],
            ['GC Cycles', (d.resource_after.num_gc - d.resource_before.num_gc)],
            ['Goroutines (before)', d.resource_before.num_goroutines],
            ['Goroutines (after)', d.resource_after.num_goroutines],
            ['CPUs', d.resource_before.num_cpu],
            ['Max CPU Limit', d.max_cpu],
            ['Max Memory', d.max_memory_mb + ' MB'],
        ];
        res.forEach(([lbl, val]) => {
            html += '<div class="res-item"><span class="res-label">' + lbl + '</span><span class="res-val">' + val + '</span></div>';
        });
        html += '</div></div>';

        document.getElementById('app').innerHTML = html;
    }
    init();
    </script>
</body>
</html>`
