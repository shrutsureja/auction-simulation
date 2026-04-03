package resource

import (
	"auction-simulation/internal/config"
	"fmt"
	"runtime"
)

type Snapshot struct {
	AllocMB      float64 `json:"alloc_mb"`      // currently allocated heap memory
	TotalAllocMB float64 `json:"total_alloc_mb"` // cumulative allocated (includes freed)
	SysMB        float64 `json:"sys_mb"`         // total memory obtained from OS
	NumGC        uint32  `json:"num_gc"`         // number of GC cycles
	NumGoroutine int     `json:"num_goroutines"` // active goroutines
	NumCPU       int     `json:"num_cpu"`        // available logical CPUs
}

func TakeSnapshot(cfg *config.Config) Snapshot {
	if cfg == nil {
		cfg = &config.Config{
			MaxCPU:         runtime.NumCPU(),
			MaxMemoryBytes: 0, // no limit
		}
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Snapshot{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       cfg.MaxCPU,
	}
}

func (s Snapshot) String() string {
	return fmt.Sprintf("Alloc=%.2fMB | TotalAlloc=%.2fMB | Sys=%.2fMB | GC=%d | Goroutines=%d | CPUs=%d",
		s.AllocMB, s.TotalAllocMB, s.SysMB, s.NumGC, s.NumGoroutine, s.NumCPU)
}
