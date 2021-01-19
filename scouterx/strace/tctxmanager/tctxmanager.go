package tctxmanager

import (
	"context"
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/structure/cachemap"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"time"
)

type ctxKey struct{}
var tctxKey = ctxKey{}

var ac = conf.GetInstance()
var txidMap = cachemap.New(10000)

var fEndStuckServiceForcibly func(tctx *netio.TraceContext)
var noopTctx = netio.NewNoopTraceContext()

func RegisterEndStuckServiceForciblyFunc(f func(tctx *netio.TraceContext)) {
	fEndStuckServiceForcibly = f
}

func Size() int {
	//TODO Size
	return 0
}

func GetTraceContextFallbackNoop(ctx context.Context) *netio.TraceContext {
	tctx := GetTraceContext(ctx)
	if tctx == nil {
		return noopTctx
	}
	return tctx
}

func GetTraceContext(ctx context.Context) *netio.TraceContext {
	traceContext, ok := ctx.Value(tctxKey).(*netio.TraceContext)
	if ok {
		return traceContext
	} else {
		return nil
	}
}

func GetTraceContextByTxid(txid int64) *netio.TraceContext {
	v := txidMap.Get(txid)
	if v == nil {
		return nil
	}
	return v.(*netio.TraceContext)
}

func NewTraceContext(ctx context.Context) (context.Context, *netio.TraceContext) {
	tctx := netio.NewTraceContext()
	newCtx := context.WithValue(ctx, tctxKey, tctx)
	return newCtx, tctx
}

func Start(tctx *netio.TraceContext) {
	txidMap.Add(tctx.Txid, tctx)
}

func End(tctx *netio.TraceContext) {
	txidMap.Remove(tctx.Txid)
}

func GetAllTctx() []interface{} {
	return txidMap.GetValues()
}

func GetActiveCount() [3]int64 {
	max := 2000
	count := 0

	values := txidMap.GetValues()
	active := [3]int64{}
	now := time.Now()
	for _, v := range values {
		count++
		if  count > max {
			return active
		}
		tctx := v.(*netio.TraceContext)
		elapsed := int(now.Sub(tctx.StartTime) / time.Millisecond)
		if elapsed <= ac.TraceActiveserviceYellowTime {
			active[0]++
		} else if elapsed <= ac.TraceActiveserviceRedTime {
			active[1]++
		} else {
			active[2]++
		}
		if elapsed > ac.StuckServiceBaseTimeMs {
			netio.SendAlert(netdata.AlertError, "STUCK_SERVICE", fmt.Sprintf("service: %s, elapsed: %d", tctx.ServiceName, elapsed))
			fEndStuckServiceForcibly(tctx)
		}
	}
	return active
}
