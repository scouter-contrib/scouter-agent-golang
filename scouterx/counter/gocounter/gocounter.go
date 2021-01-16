package gocounter

import (
	"runtime"
	"sync"
	"time"
)

var statsMux sync.Mutex

type GoCounter struct {
	Timestamp    int64
	GoroutineNum int
	CgoCallNum   int64

	// heap
	HeapAlloc    int64

	// garbage collection
	GcPerSec      float64 //gc count per sec
	GcPausePerSec float64 //pause ms per sec
	mPauseTotalNs  int64
	mNumGC         int32
}

func GetGoCounter(prev GoCounter) GoCounter {
	statsMux.Lock()
	defer statsMux.Unlock()

	nowNano := time.Now().UnixNano()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	counter := GoCounter{Timestamp: nowNano}

	counter.GoroutineNum = runtime.NumGoroutine()
	counter.CgoCallNum = runtime.NumCgoCall()
	counter.HeapAlloc = int64(mem.HeapAlloc)

	counter.mPauseTotalNs = int64(mem.PauseTotalNs)
	counter.mNumGC = int32(mem.NumGC)

	if prev.Timestamp != 0 {
		counter.GcPausePerSec = float64(int64(mem.PauseTotalNs) - prev.mPauseTotalNs) / float64(nowNano - prev.Timestamp) * 1000
		counter.GcPerSec = float64(int32(mem.NumGC) - prev.mNumGC) / float64(nowNano - prev.Timestamp) * float64(time.Second)
	}

	return counter
}


//type StatsTest struct {
//	Time         int64
//	GoVersion    string
//	GoOs         string
//	GoArch       string
//	CpuNum       int
//	GoroutineNum int
//	Gomaxprocs   int
//	CgoCallNum   int64
//
//	// memory
//	MemoryAlloc      uint64
//	MemoryTotalAlloc uint64
//	MemorySys        uint64
//	MemoryLookups    uint64
//	MemoryMallocs    uint64
//	MemoryFrees      uint64
//
//	// stack
//	StackInUse uint64
//	// heap
//	HeapAlloc    uint64
//	HeapSys      uint64
//	HeapIdle     uint64
//	HeapInuse    uint64
//	HeapReleased uint64
//	HeapObjects  uint64
//
//	// garbage collection
//	GcNext           uint64
//	GcLast           uint64
//	GcNum            uint32
//	GcPerSecond      float64
//	GcPausePerSecond float64
//	GcPause          []float64
//}
//
//var lastSampleTime time.Time
//var lastPauseNs uint64 = 0
//var lastNumGc uint32 = 0
//
//var nsInMs float64 = float64(time.Millisecond)
//

//func GetStatsTest() *StatsTest {
//	statsMux.Lock()
//	defer statsMux.Unlock()
//
//	var mem runtime.MemStats
//	runtime.ReadMemStats(&mem)
//
//	now := time.Now()
//
//	var gcPausePerSecond float64
//
//	if lastPauseNs > 0 {
//		pauseSinceLastSample := mem.PauseTotalNs - lastPauseNs
//		gcPausePerSecond = float64(pauseSinceLastSample) / nsInMs
//	}
//
//	lastPauseNs = mem.PauseTotalNs
//
//	countGc := int(mem.NumGC - lastNumGc)
//
//	var gcPerSecond float64
//
//	if lastNumGc > 0 {
//		diff := float64(countGc)
//		diffTime := now.Sub(lastSampleTime).Seconds()
//		gcPerSecond = diff / diffTime
//	}
//
//	if countGc > 256 {
//		// lagging GC pause times
//		countGc = 256
//	}
//
//	gcPause := make([]float64, countGc)
//
//	for i := 0; i < countGc; i++ {
//		idx := int((mem.NumGC-uint32(i))+255) % 256
//		pause := float64(mem.PauseNs[idx])
//		gcPause[i] = pause / nsInMs
//	}
//
//	lastNumGc = mem.NumGC
//	lastSampleTime = time.Now()
//
//	return &StatsTest{
//		Time:         now.UnixNano(),
//		GoVersion:    runtime.Version(),
//		GoOs:         runtime.GOOS,
//		GoArch:       runtime.GOARCH,
//		CpuNum:       runtime.NumCPU(),
//		GoroutineNum: runtime.NumGoroutine(),
//		Gomaxprocs:   runtime.GOMAXPROCS(0),
//		CgoCallNum:   runtime.NumCgoCall(),
//		// memory
//		MemoryAlloc:      mem.Alloc,
//		MemoryTotalAlloc: mem.TotalAlloc,
//		MemorySys:        mem.Sys,
//		MemoryLookups:    mem.Lookups,
//		MemoryMallocs:    mem.Mallocs,
//		MemoryFrees:      mem.Frees,
//		// stack
//		StackInUse: mem.StackInuse,
//		// heap
//		HeapAlloc:    mem.HeapAlloc,
//		HeapSys:      mem.HeapSys,
//		HeapIdle:     mem.HeapIdle,
//		HeapInuse:    mem.HeapInuse,
//		HeapReleased: mem.HeapReleased,
//		HeapObjects:  mem.HeapObjects,
//		// garbage collection
//		GcNext:           mem.NextGC,
//		GcLast:           mem.LastGC,
//		GcNum:            mem.NumGC,
//		GcPerSecond:      gcPerSecond,
//		GcPausePerSecond: gcPausePerSecond,
//		GcPause:          gcPause,
//	}
//}
