package tctxmanager

import (
	"context"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/structure/cachemap"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
)

type ctxKey struct{}
var tctxKey = ctxKey{}

var ac = conf.GetInstance()
var txidMap = cachemap.New(10000)

func Size() int {
	//TODO Size
	return 0
}

func GetActiveCount() []int {
	return nil
}

func GetTraceContext(ctx context.Context) *netio.TraceContext {
	traceContext, ok := ctx.Value(tctxKey).(*netio.TraceContext)
	if ok {
		return traceContext
	} else {
		return nil
	}
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

