package netio

import (
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata/texttype"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/structure/lra"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio/udpsender"
	"time"
)

var ac = conf.GetInstance()

var serviceNameSent = lra.New(10000)
var objNameSent = lra.New(100)
var refererSent = lra.New(10000)
var userAgentSent = lra.New(10000)
var methodSent = lra.New(10000)
var apicallSent = lra.New(10000)
var errorSent = lra.New(10000)
var loginSent = lra.New(10000)
var descSent = lra.New(10000)
var stackElementSent = lra.New(10000)
var hashMessageSent = lra.New(10000)

func ResetTextSent() {
	serviceNameSent.Clear()
	objNameSent.Clear()
	refererSent.Clear()
	userAgentSent.Clear()
	methodSent.Clear()
	apicallSent.Clear()
	errorSent.Clear()
	loginSent.Clear()
	descSent.Clear()
	stackElementSent.Clear()
	hashMessageSent.Clear()
}
func reportScouterPanicFinal() {
	if r := recover(); r != nil {
		fmt.Printf("[scouter][panic][FATAL]%+v\n", r)
	}
}

func SendPanic(message string) {
	reportScouterPanicFinal()
	SendAlert(netdata.AlertFatal, "Panic", message)
}

func SendServiceName(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if serviceNameSent.Contains(hash) {
		return hash
	}
	serviceNameSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.SERVICE, hash, name})
	return hash
}

func SendHashedMessage(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if hashMessageSent.Contains(hash) {
		return hash
	}
	hashMessageSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.HASH_MSG, hash, name})
	return hash
}

func SendObjName(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if objNameSent.Contains(hash) {
		return hash
	}
	objNameSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.OBJECT, hash, name})
	return hash
}

func SendReferer(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if refererSent.Contains(hash) {
		return hash
	}
	refererSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.REFERER, hash, name})
	return hash
}

func SendUserAgent(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if userAgentSent.Contains(hash) {
		return hash
	}
	userAgentSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.USER_AGENT, hash, name})
	return hash
}

func SendMethod(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if methodSent.Contains(hash) {
		return hash
	}
	methodSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.METHOD, hash, name})
	return hash
}

func SendApicall(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if apicallSent.Contains(hash) {
		return hash
	}
	apicallSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.APICALL, hash, name})
	return hash
}

func SendError(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if errorSent.Contains(hash) {
		return hash
	}
	errorSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.ERROR, hash, name})
	return hash
}

func SendAlert(level netdata.AlertLevel, title, message string) {
	if title == "" {
		return
	}
	pack := netdata.NewAlertPack()
	pack.ObjType = ac.ObjType
	pack.ObjHash = ac.ObjHash
	pack.Level = level
	pack.Title = title
	pack.Message = message

	SendPackDirect(pack)
}

func SendLogin(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if loginSent.Contains(hash) {
		return hash
	}
	loginSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.LOGIN, hash, name})
	return hash
}

func SendDesc(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if descSent.Contains(hash) {
		return hash
	}
	descSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.DESC, hash, name})
	return hash
}

func SendStackElement(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if stackElementSent.Contains(hash) {
		return hash
	}
	stackElementSent.AddKey(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.STACK_ELEMENT, hash, name})
	return hash
}

func SendXlog(pack *netdata.XlogPack) {
	pack.ObjHash = ac.ObjHash
	udpsender.GetInstance().SendPackDirect(pack)
}

func SendProfile(steps []netdata.Step, tctx *TraceContext) {
	if steps == nil || len(steps) == 0 {
		return
	}
	bulkSize := ac.ProfileStepMaxCount
	count := len(steps) / bulkSize
	if count == 0 {
		sendProfile0(steps, tctx)
		return
	}
	remainder := len(steps) % bulkSize
	for i := 0; i < count; i++ {
		sendProfile0(steps[i*bulkSize:i*bulkSize+bulkSize], tctx)
	}
	if remainder > 0 {
		sendProfile0(steps[count*bulkSize:count*bulkSize+remainder], tctx)
	}
}

func sendProfile0(steps []netdata.Step, tctx *TraceContext) {
	if steps == nil || len(steps) == 0 {
		return
	}

	pack := netdata.NewXlogProfilePack()
	pack.Txid = tctx.Txid
	pack.ObjHash = ac.ObjHash
	pack.Profile = netdata.StepsToBytes(steps)
	pack.Service = tctx.ServiceHash
	pack.Elapsed = int32(time.Now().Sub(tctx.StartTime) / time.Millisecond)
	tctx.ProfileCount += len(steps)
	tctx.ProfileSize += len(pack.Profile)

	udpsender.GetInstance().SendPackDirect(pack)
}

func SendPackDirect(pack netdata.Pack) {
	udpsender.GetInstance().SendPackDirect(pack)
}
