package strace

import (
	"context"
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/strace/tctxmanager"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"github.com/scouter-project/scouter-go-lib/common/util"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var ac = conf.GetInstance()

func goid() int {
	defer reportScouterPanic()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func ForceErrorMark(ctx context.Context, errorMessage string) {
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return
	}
	tctx.Error = netio.SendError(errorMessage)
}

func ForceSetServiceName(ctx context.Context, serviceName string) {
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return
	}
	tctx.ServiceName = serviceName
	tctx.ServiceHash = netio.SendServiceName(serviceName)
}

func AddStep(ctx context.Context, step netdata.Step) {
	defer reportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	tctx.Profile.Push(step)
}

func AddMessageStep(ctx context.Context, message string) {
	defer reportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)

	step := netdata.NewMessageStep(message, util.MillisToNow(tctx.StartTime))
	tctx.Profile.Push(step)
}

func AddHashedMessageStep(ctx context.Context, message string, value, elapsed int32) {
	defer reportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)

	step := netdata.NewHashedMessageStep(netio.SendHashedMessage(message), util.MillisToNow(tctx.StartTime))
	step.Value = value
	step.Time = elapsed
	tctx.Profile.Push(step)
}

func AddPMessageStep(ctx context.Context, level netdata.PMessageLevel, message string, elapsed int32, params ...string) {
	defer reportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)

	step := netdata.NewPMessageStep(util.MillisToNow(tctx.StartTime))
	step.SetMessage(netio.SendHashedMessage(message), params...)
	step.Elapsed = elapsed
	step.Level = level
	tctx.Profile.Push(step)
}

func StartHttpService(ctx context.Context, req http.Request) (newCtx context.Context) {
	defer reportScouterPanic()
	if ctx == nil {
		return context.Background()
	}
	if tctxmanager.GetTraceContext(ctx) != nil {
		return ctx
	}

	serviceName := strings.Join([]string{req.URL.Path, "<", req.Method, ">"}, "")

	//TODO propagation request (gxid, caller)
	//TODO query profile
	//TODO body (of specific service) profile

	newCtx, tctx := startService(ctx, serviceName, getRemoteIp(req))
	tctx.XType = netdata.XTYPE_WEB_SERVICE
	tctx.UserAgent = netio.SendUserAgent(req.UserAgent())
	tctx.HttpMethod = req.Method
	tctx.Referer = netio.SendReferer(req.Referer())
	profileHttpHeaders(req, tctx)

	return newCtx
}

func getRemoteIp(req http.Request) string {
	ip := req.RemoteAddr
	if ac.TraceHttpClientIpHeaderKey != "" {
		headerIp := req.Header.Get(ac.TraceHttpClientIpHeaderKey)
		if headerIp != "" {
			ip = headerIp
		}
	}
	return strings.Split(ip, ":")[0]
}

func EndHttpService(ctx context.Context, req http.Request, res http.Response) {
	//TODO body (of specific service) profile from req.body
	defer reportScouterPanic()
	endAnyService(ctx)
}

func StartService(ctx context.Context, serviceName, remoteIp string) (newCtx context.Context) {
	defer reportScouterPanic()
	if ctx == nil {
		ctx = context.Background()
	}
	if tctxmanager.GetTraceContext(ctx) != nil {
		return ctx
	}
	newCtx, tctx := startService(ctx, serviceName, remoteIp)
	tctx.XType = netdata.XTYPE_APP_SERVICE

	return newCtx
}

func EndService(ctx context.Context) {
	defer reportScouterPanic()
	endAnyService(ctx)
}

//<usage> for chained tracing
//go func() {
//	ctxForGoroutine := StartChildGoroutineService(ctx, "myGoroutineService1()")
//	EndChildGoroutineService(ctxForGoroutine)
//	myGoroutineService1(ctxForGoroutine)
//}()
func StartChildGoroutineService(ctx context.Context, serviceName string) (ctx4Goroutine context.Context) {
	defer reportScouterPanic()
	if ctx == nil {
		return ctx
	}
	parentTctx := tctxmanager.GetTraceContext(ctx)
	if parentTctx == nil {
		return ctx
	}

	ctx4Goroutine, childTctx := startService(ctx, serviceName, parentTctx.RemoteIp)
	childTctx.XType = netdata.XTYPE_BACK_THREAD2
	childTctx.Caller = parentTctx.Txid
	childTctx.Gxid = parentTctx.Gxid
	if childTctx.Gxid == 0 {
		childTctx.Gxid = parentTctx.Txid
		parentTctx.Gxid = parentTctx.Txid
	}

	asyncStep := netdata.NewAsyncServiceStep()
	asyncStep.Txid = childTctx.Txid
	asyncStep.StartTime = util.MillisToNow(parentTctx.StartTime)
	asyncStep.Hash = netio.SendApicall(serviceName)
	parentTctx.Profile.Add(asyncStep)

	return ctx4Goroutine
}

func EndChildGoroutineService(ctx context.Context) {
	defer reportScouterPanic()
	endAnyService(ctx)
}

func startService(ctx context.Context, serviceName, remoteAddr string) (context.Context, *netio.TraceContext) {
	newCtx, tctx := tctxmanager.NewTraceContext(ctx)
	tctxmanager.Start(tctx)

	tctx.Goid = goid()
	tctx.Profile.Add(netdata.NewMessageStep(fmt.Sprintf("goroutine:%d", tctx.Goid), 0))

	tctx.ServiceName = serviceName
	tctx.ServiceHash = netio.SendServiceName(serviceName)
	tctx.RemoteIp = remoteAddr
	return newCtx, tctx
}

func endAnyService(ctx context.Context) {
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return
	}
	tctx.ServiceHash = netio.SendServiceName(tctx.ServiceName)
	tctxmanager.End(tctx)

	elapsed := util.MillisToNow(tctx.StartTime)
	discardType := findXLogDiscard(tctx, elapsed)
	xlog := tctx.ToXlog(discardType, elapsed)

	writeProfile := discardType == netdata.XLOG_DISCARD_NONE //TODO consequence sampling
	tctx.Profile.Close(writeProfile)

	//TODO metering
	//TODO meteringInteraction

	if (xlog.DiscardType != netdata.XLOG_DISCARD_ALL && xlog.DiscardType != netdata.XLOG_DISCARD_ALL_FORCE) ||
		(!xlog.IsDriving() && xlog.DiscardType == netdata.XLOG_DISCARD_ALL) {
		netio.SendXlog(xlog)
	} else {
		//TODO send Dropped XLog
	}
}

func findXLogDiscard(tctx *netio.TraceContext, elapsed int32) netdata.XlogDiscardType {
	var discardMode netdata.XlogDiscardType
	if tctx.Error != 0 {
		discardMode = netdata.XLOG_DISCARD_NONE
	} else {
		discardMode = GetXlogSampler().eEvaluateXLogDiscard(elapsed, tctx.ServiceName)
	}
	//check xlog discard pattern
	if GetXlogSampler().isDiscardServicePattern(tctx.ServiceName) {
		discardMode = netdata.XLOG_DISCARD_ALL_FORCE
		if tctx.Error != 0 {
			discardMode = netdata.XLOG_DISCARD_NONE
		}
	}
	return discardMode
}

func StartMethod(ctx context.Context) *netdata.MethodStep {
	return StartMethodWithParam(ctx, nil)
}

func StartMethodWithParam(ctx context.Context, params []interface{}) *netdata.MethodStep {
	defer reportScouterPanic()

	if ctx == nil {
		return nil
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return nil
	}

	pc, _, _, _ := runtime.Caller(2)
	funcName := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	split := strings.Split(funcName, "/")
	shortName := split[len(split)-1] + "()"

	addMessageStepIfParamExist(tctx, params)

	step := netdata.NewMethodStep()
	step.Hash = netio.SendMethod(shortName)
	step.StartTime = util.MillisToNow(tctx.StartTime)
	tctx.Profile.Push(step)
	return step
}

func addMessageStepIfParamExist(tctx *netio.TraceContext, params []interface{}) {
	if params == nil || len(params) == 0 {
		return
	}
	for _, param := range params {
		step := netdata.NewMessageStep(fmt.Sprintf("param: %v", param), util.MillisToNow(tctx.StartTime))
		step.StartTime = util.MillisToNow(tctx.StartTime)
		tctx.Profile.Add(step)
	}
}

func EndMethod(ctx context.Context, step *netdata.MethodStep) {
	defer reportScouterPanic()

	if ctx == nil || step == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return
	}
	step.Elapsed = util.MillisToNow(tctx.StartTime) - step.StartTime
	tctx.Profile.Pop(step)
}

func profileHttpHeaders(r http.Request, tctx *netio.TraceContext) {
	startTime := util.MillisToNow(tctx.StartTime)
	if ac.ProfileHttpHeaderEnabled {
		notAll := len(ac.ProfileHttpHeaderKeys) > 0
		if notAll {
			split := strings.Split(ac.ProfileHttpHeaderKeys, ",")
			for _, k := range split {
				values := r.Header.Values(strings.TrimSpace(k))
				if values != nil && len(values) > 0 {
					v := strings.Join(values, ",")
					tctx.Profile.Add(netdata.NewMessageStep(fmt.Sprintf("header: %s: %s", k, v), startTime))
				}
			}
		} else {
			for k, v := range r.Header {
				vs := strings.Join(v, ",")
				tctx.Profile.Add(netdata.NewMessageStep(fmt.Sprintf("header: %s: %s", k, vs), startTime))
			}
		}
	}

	if ac.ProfileHttpQuerystringEnabled {
		tctx.Profile.Add(netdata.NewMessageStep(fmt.Sprintf("query: %s", r.URL.RawQuery), startTime))
	}
}

func reportScouterPanic() {
	if r := recover(); r != nil {
		fmt.Printf("[scouter][panic]%+v\n", r)
	}
}
