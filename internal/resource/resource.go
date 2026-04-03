package resource

import (
	"fmt"
	"runtime"
)

type Snapshot struct {
	AllocMB      float64 // currently allocated heap memory
	TotalAllocMB float64 // cumulative allocated (includes freed)
	SysMB        float64 // total memory obtained from OS
	NumGC        uint32  // number of GC cycles
	NumGoroutine int     // active goroutines
	NumCPU       int     // available logical CPUs
}

func TakeSnapshot() Snapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Snapshot{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       2,
	}
}

func (s Snapshot) String() string {
	return fmt.Sprintf("Alloc=%.2fMB | TotalAlloc=%.2fMB | Sys=%.2fMB | GC=%d | Goroutines=%d | CPUs=%d",
		s.AllocMB, s.TotalAllocMB, s.SysMB, s.NumGC, s.NumGoroutine, s.NumCPU)
}
