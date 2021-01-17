package strace

import (
	"context"
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/counter"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio/tcpclient"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/strace/tctxmanager"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var ac = conf.GetInstance()

func StartTracingMode() {
	tctxmanager.RegisterEndStuckServiceForciblyFunc(endStuckServiceForcibly)
	tcpclient.StartTcp()
}

func goid() int {
	defer common.ReportScouterPanic()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func IsStream(ctx context.Context) bool {
	common.ReportScouterPanic()
	if ctx == nil {
		return false
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return false
	}
	return tctx.IsStream
}

func SetAsStream(ctx context.Context) bool {
	common.ReportScouterPanic()
	if ctx == nil {
		return false
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return false
	}
	tctx.IsStream = true
	return true
}

func MarkAsError(ctx context.Context, errorMessage string) {
	common.ReportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return
	}
	if tctx.Error == 0 {
		tctx.Error = netio.SendError(errorMessage)
	}
	AddPMessageStep(ctx, netdata.PMSG_ERROR, errorMessage, 0)
}

func MarkAsErrorForcibly(ctx context.Context, errorMessage string) {
	common.ReportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return
	}
	tctx.Error = netio.SendError(errorMessage)
}

func SetServiceNameForcibly(ctx context.Context, serviceName string) {
	common.ReportScouterPanic()
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
	defer common.ReportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	tctx.Profile.Add(step)
}

func AddMessageStep(ctx context.Context, message string) {
	defer common.ReportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)

	step := netdata.NewMessageStep(message, util.MillisToNow(tctx.StartTime))
	tctx.Profile.Add(step)
}

func AddHashedMessageStep(ctx context.Context, message string, value, elapsed int32) {
	defer common.ReportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)

	step := netdata.NewHashedMessageStep(netio.SendHashedMessage(message), util.MillisToNow(tctx.StartTime))
	step.Value = value
	step.Time = elapsed
	tctx.Profile.Add(step)
}

func AddPMessageStep(ctx context.Context, level netdata.PMessageLevel, message string, elapsed int32, params ...string) {
	defer common.ReportScouterPanic()
	if ctx == nil {
		return
	}
	tctx := tctxmanager.GetTraceContext(ctx)

	step := netdata.NewPMessageStep(util.MillisToNow(tctx.StartTime))
	step.SetMessage(netio.SendHashedMessage(message), params...)
	step.Elapsed = elapsed
	step.Level = level
	tctx.Profile.Add(step)
}

func StartHttpService(ctx context.Context, req *http.Request) (newCtx context.Context) {
	defer common.ReportScouterPanic()
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

func getRemoteIp(req *http.Request) string {
	ip := req.RemoteAddr
	if ac.TraceHttpClientIpHeaderKey != "" {
		headerIp := req.Header.Get(ac.TraceHttpClientIpHeaderKey)
		if headerIp != "" {
			ip = headerIp
		}
	}
	return strings.Split(ip, ":")[0]
}

func EndHttpService(ctx context.Context, req *http.Request, res *http.Response) {
	common.ReportScouterPanic()
	//TODO body (of specific service) profile from req.body

	if res != nil {
		if ctx == nil {
			return
		}
		tctx := tctxmanager.GetTraceContext(ctx)
		if tctx == nil || tctx.Closed {
			return
		}
		tctx.Status = int32(res.StatusCode)
	}
	endAnyService(ctx)
}

func StartService(ctx context.Context, serviceName, remoteIp string) (newCtx context.Context) {
	common.ReportScouterPanic()
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
	common.ReportScouterPanic()
	endAnyService(ctx)
}

func StartNewInheritanceService(ctx context.Context, parentTctx *netio.TraceContext) (newCtx context.Context, newTctx *netio.TraceContext) {
	common.ReportScouterPanic()
	if ctx == nil {
		ctx = context.Background()
	}
	newCtx, newTctx = startService(ctx, parentTctx.ServiceName, parentTctx.RemoteIp)
	newTctx = inheritTctx(newTctx, parentTctx)

	return newCtx, newTctx
}

func inheritTctx(newTctx *netio.TraceContext, parentTctx *netio.TraceContext) *netio.TraceContext {
	newTctx.Inherit = true
	newTctx.Gxid = parentTctx.Gxid
	newTctx.XType = parentTctx.XType
	newTctx.Profile.Add(netdata.NewMessageStep("scouter inheritance step", 0))
	newTctx.IsStream = parentTctx.IsStream

	newTctx.Error = parentTctx.Error
	newTctx.HttpMethod = parentTctx.HttpMethod
	newTctx.HttpQuery = parentTctx.HttpQuery
	newTctx.HttpContentType = parentTctx.HttpContentType

	newTctx.SqlCount = parentTctx.SqlCount
	newTctx.SqlTime = parentTctx.SqlTime
	newTctx.Sqltext = parentTctx.Sqltext

	newTctx.ApicallName = parentTctx.ApicallName
	newTctx.ApicallCount = parentTctx.ApicallCount
	newTctx.ApicallTime = parentTctx.ApicallTime
	newTctx.ApicallTarget = parentTctx.ApicallTarget

	newTctx.Userid = parentTctx.Userid
	newTctx.UserAgent = parentTctx.UserAgent
	newTctx.UserAgentString = parentTctx.UserAgentString
	newTctx.Referer = parentTctx.Referer

	newTctx.IsChildTx = true
	newTctx.Caller = parentTctx.Txid
	newTctx.CallerObjHash = ac.ObjHash

	newTctx.Login = parentTctx.Login
	newTctx.Desc = parentTctx.Desc

	newTctx.Text1 = parentTctx.Text1
	newTctx.Text2 = parentTctx.Text2
	newTctx.Text3 = parentTctx.Text3
	newTctx.Text4 = parentTctx.Text4
	newTctx.Text5 = parentTctx.Text5

	return newTctx
}

//<usage> for chained goroutine tracing
//
//GoWithTrace(ctx, "myFuncName()", func(cascadeGoCtx context.Context) {
//	myFunc(cascadeGoCtx)
//})
func GoWithTrace(ctx context.Context, serviceName string, func4Goroutine func(cascadeGoCtx context.Context)) {
	common.ReportScouterPanic()
	newCtx, childTctx := startChildGoroutineService(ctx, serviceName)
	go func() {
		if childTctx != nil {
			childTctx.StartTime = time.Now()
		}
		defer endChildGoroutineService(newCtx)
		func4Goroutine(newCtx)
	}()
}

func startChildGoroutineService(ctx context.Context, serviceName string) (ctx4Goroutine context.Context, childTctx *netio.TraceContext) {
	if ctx == nil {
		return ctx, nil
	}
	parentTctx := tctxmanager.GetTraceContext(ctx)
	if parentTctx == nil {
		return ctx, nil
	}

	ctx4Goroutine, childTctx = startService(ctx, serviceName, parentTctx.RemoteIp)
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

	return ctx4Goroutine, childTctx
}

func endChildGoroutineService(ctx context.Context) {
	common.ReportScouterPanic()
	endAnyService(ctx)
}

func startService(ctx context.Context, serviceName, remoteAddr string) (context.Context, *netio.TraceContext) {
	newCtx, tctx := tctxmanager.NewTraceContext(ctx)
	tctxmanager.Start(tctx)

	tctx.Gxid = tctx.Txid
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
	if tctx == nil || tctx.Closed {
		return
	}
	endAnyServiceOfTraceContext(tctx)
}

func endStuckServiceForcibly(tctx *netio.TraceContext) {
	if ac.StuckServiceRemoveEnabled {
		step := netdata.NewPMessageStep(util.MillisToNow(tctx.StartTime))
		step.SetMessage(netio.SendHashedMessage("Service currently may running, not finished!"))
		step.Level = netdata.PMSG_ERROR
		tctx.Profile.Add(step)

		if tctx.Error == 0 {
			tctx.Error = netio.SendError("This stuck service currently may running, not finished!")
		}
		endAnyServiceOfTraceContext(tctx)
	}
}

func endAnyServiceOfTraceContext(tctx *netio.TraceContext) {
	if tctx.Closed {
		return
	}
	tctx.Closed = true

	tctx.ServiceHash = netio.SendServiceName(tctx.ServiceName)
	tctxmanager.End(tctx)

	elapsed := util.MillisToNow(tctx.StartTime)
	discardType := findXLogDiscard(tctx, elapsed)
	xlog := tctx.ToXlog(discardType, elapsed)

	writeProfile := discardType == netdata.XLOG_DISCARD_NONE //TODO consequence sampling
	tctx.Profile.Close(writeProfile)

	counter.GetServiceMeter().Add(int(xlog.Elapsed), xlog.Error != 0)
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
	defer common.ReportScouterPanic()
	return startMethodWithParam(ctx) //do not call StartMethodWithParam (because of method name finding depth)
}

func StartMethodWithParam(ctx context.Context, params ...interface{}) *netdata.MethodStep {
	defer common.ReportScouterPanic()
	return startMethodWithParam(ctx, params)
}

func StartCustomMethod(ctx context.Context, methodName string) *netdata.MethodStep {
	return StartCustomMethodWithParam(ctx, methodName)
}

func StartCustomMethodWithParam(ctx context.Context, methodName string, params ...interface{}) *netdata.MethodStep {
	defer common.ReportScouterPanic()
	if ctx == nil {
		return nil
	}
	tctx := tctxmanager.GetTraceContext(ctx)
	if tctx == nil {
		return nil
	}

	return startMethodWithParam0(tctx, methodName, methodName, params)
}

func startMethodWithParam(ctx context.Context, params ...interface{}) *netdata.MethodStep {
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
	methodName := split[len(split)-1] + "()"

	return startMethodWithParam0(tctx, funcName, methodName, params)
}

func startMethodWithParam0(tctx *netio.TraceContext, funcName string, methodName string, params ...interface{}) *netdata.MethodStep {
	addMessageStepIfParamExist(tctx, params)
	tctx.LastMethod = funcName

	step := netdata.NewMethodStep()
	step.Hash = netio.SendMethod(methodName)
	step.StartTime = util.MillisToNow(tctx.StartTime)
	tctx.Profile.Push(step)
	return step
}

func addMessageStepIfParamExist(tctx *netio.TraceContext, params []interface{}) {
	if params == nil || len(params) == 0 {
		return
	}
	for _, param := range params {
		if param == nil {
			continue
		}
		step := netdata.NewMessageStep(fmt.Sprintf("param: %v", param), util.MillisToNow(tctx.StartTime))
		step.StartTime = util.MillisToNow(tctx.StartTime)
		tctx.Profile.Add(step)
	}
}

func EndMethod(ctx context.Context, step *netdata.MethodStep) {
	defer common.ReportScouterPanic()

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

func profileHttpHeaders(r *http.Request, tctx *netio.TraceContext) {
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
