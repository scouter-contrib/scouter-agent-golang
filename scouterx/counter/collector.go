package counter

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/counter/gocounter"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"time"
)

var prevGoCounter = gocounter.GoCounter{}
var activeCounter = NewValueMeter()

func Start() {
	go start()
}

func start() {
	time.Sleep(2000 * time.Millisecond)
	SendServiceCounter()
	SendGoCounter()
}

func SendGoCounter() {
	ac := conf.GetInstance()
	pack := netdata.NewPerfCounterPack()
	pack.ObjName = ac.ObjName
	pack.Timetype = netdata.TimeTypeRealTime

	c := gocounter.GetGoCounter(prevGoCounter)
	prevGoCounter = c

	pack.Put(GO_GOROUTINE, c.GoroutineNum)
	pack.Put(GO_CGO_CALL, c.CgoCallNum)
	pack.Put(GO_GC_COUNT, c.GcPerSec)
	pack.Put(GO_GC_PAUSE, c.GcPausePerSec)
	pack.Put(GO_HEAP_USED, c.HeapAlloc)

	netio.SendPackDirect(pack)
}

func SendServiceCounter() {
	ac := conf.GetInstance()
	pack := netdata.NewPerfCounterPack()
	pack.ObjName = ac.ObjName
	pack.Timetype = netdata.TimeTypeRealTime

	counter := GetServiceMeter().GetAllCounter(30)

	//TODO activeCounter
	//TODO WAS_ACTIVE_SERVICE, WAS_ACTIVE_SPEED
	pack.Put(TPS, counter.ErrorRate)
	pack.Put(ELAPSED_TIME, counter.Elapsed)
	pack.Put(ERROR_RATE, counter.Tps)

	netio.SendPackDirect(pack)
}
