package countertask

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/counter"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/counter/gocounter"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/strace/tctxmanager"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"time"
)

var prevGoCounter = gocounter.GoCounter{}
var activeCounter = counter.NewValueMeter()

func Start() {
	go startRealTimeCounter()
	go start5MinCounter()
}

func startRealTimeCounter() {
	for {
		time.Sleep(2000 * time.Millisecond)
		pack := getAllCounterOnPack()
		netio.SendPackDirect(pack)
	}
}

func start5MinCounter() {
	for {
		time.Sleep(10000 * time.Millisecond)
		pack := getAll5MinCounterOnPack()
		netio.SendPackDirect(pack)
	}
}

func getAllCounterOnPack() netdata.Pack {
	ac := conf.GetInstance()
	pack := netdata.NewPerfCounterPack()
	pack.ObjName = ac.ObjName
	pack.Timetype = netdata.TimeTypeRealTime

	c := gocounter.GetGoCounter(prevGoCounter)
	prevGoCounter = c

	pack.Put(counter.GO_GOROUTINE, c.GoroutineNum)
	pack.Put(counter.GO_CGO_CALL, c.CgoCallNum)
	pack.Put(counter.GO_GC_COUNT, c.GcPerSec)
	pack.Put(counter.GO_GC_PAUSE, c.GcPausePerSec)
	pack.Put(counter.GO_HEAP_USED, c.HeapAlloc)

	sc := counter.GetServiceMeter().GetAllCounter(30)
	pack.Put(counter.TPS, sc.Tps)
	pack.Put(counter.ELAPSED_TIME, sc.Elapsed)
	pack.Put(counter.ERROR_RATE, sc.ErrorRate)

	activeCounts := tctxmanager.GetActiveCount()
	activeSum := activeCounts[0] + activeCounts[1] + activeCounts[2]
	activeCounter.Add(float64(activeSum))
	active := activeCounter.GetAllCounter(30).Avg
	pack.Put(counter.ACTIVE_SERVICE, active)

	activeSpeed := netdata.NewListValue()
	activeSpeed.AddInt64(activeCounts[0])
	activeSpeed.AddInt64(activeCounts[1])
	activeSpeed.AddInt64(activeCounts[2])

	pack.Put(counter.ACTIVE_SPEED, activeSpeed)

	return pack
}

func getAll5MinCounterOnPack() netdata.Pack {
	ac := conf.GetInstance()
	pack := netdata.NewPerfCounterPack()
	pack.ObjName = ac.ObjName
	pack.Timetype = netdata.TimeTypeFiveMin

	sc := counter.GetServiceMeter().GetAllCounter(300)
	pack.Put(counter.TPS, sc.Tps)
	pack.Put(counter.ELAPSED_TIME, sc.Elapsed)
	pack.Put(counter.ERROR_RATE, sc.ErrorRate)

	active := activeCounter.GetAllCounter(300).Avg
	pack.Put(counter.ACTIVE_SERVICE, active)

	return pack
}
