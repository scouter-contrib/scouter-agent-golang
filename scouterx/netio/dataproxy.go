package netio

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/structure/cacheset"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio/udpsender"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"github.com/scouter-project/scouter-go-lib/common/netdata/texttype"
	"github.com/scouter-project/scouter-go-lib/common/util"
	"time"
)

var ac = conf.GetInstance()

var serviceNameSent = cacheset.New(10000)
var objNameSent = cacheset.New(100)
var refererSent = cacheset.New(10000)
var userAgentSent = cacheset.New(10000)
var methodSent = cacheset.New(10000)
var apicallSent = cacheset.New(10000)
var errorSent = cacheset.New(10000)
var loginSent = cacheset.New(10000)
var descSent = cacheset.New(10000)
var stackElementSent = cacheset.New(10000)
var hashMessageSent = cacheset.New(10000)

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
}

func SendServiceName(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if serviceNameSent.Contains(hash) {
		return hash
	}
	serviceNameSent.Add(hash)
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
	hashMessageSent.Add(hash)
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
	objNameSent.Add(hash)
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
	refererSent.Add(hash)
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
	userAgentSent.Add(hash)
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
	methodSent.Add(hash)
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
	apicallSent.Add(hash)
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
	errorSent.Add(hash)
	udpsender.GetInstance().AddPack(&netdata.TextPack{texttype.ERROR, hash, name})
	return hash
}

func SendLogin(name string) int32 {
	if name == "" {
		return 0
	}
	hash := util.HashString(name)
	if loginSent.Contains(hash) {
		return hash
	}
	loginSent.Add(hash)
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
	descSent.Add(hash)
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
	stackElementSent.Add(hash)
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
