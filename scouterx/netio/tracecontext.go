package netio

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util/keygen"
	"strconv"
	"time"
)

type TraceContext struct {
	Closed       bool
	LastMethod   string
	IsStream     bool

	Inherit bool
	InheritStartTime time.Time

	Goid         int
	Parent       *TraceContext
	Profile      *ProfileCollector
	ProfileCount int
	ProfileSize  int

	Txid  int64
	Gxid  int64
	XType netdata.XlogType

	StartTime time.Time
	//TODO discardType XLogDiscardTypes.XLogDiscard

	ServiceHash int32
	ServiceName string
	RemoteIp    string

	Error           int32
	HttpMethod      string
	HttpQuery       string
	HttpContentType string

	SqlCount int32
	SqlTime  int32
	Sqltext  string

	ApicallName   string
	ApicallCount  int32
	ApicallTime   int32
	ApicallTarget string

	//rsStart int64
	//rsCount int32
	//TODO sql = new SqlParameter()
	//sqlActiveArgs SqlParameter

	Userid          int64
	UserAgent       int32
	UserAgentString string
	Referer         int32

	IsChildTx     bool
	Caller        int64
	Callee        int64
	CallerObjHash int32

	Login string
	Desc  string

	Text1 string
	Text2 string
	Text3 string
	Text4 string
	Text5 string

	//TODO public Queue<DumpStep> temporaryDumpSteps = new LinkedBlockingQueue<DumpStep>(5);
	//TODO public boolean hasDumpStack;
}

func NewTraceContext() *TraceContext {
	context := new(TraceContext)
	context.Txid = keygen.Next()
	context.StartTime = time.Now()
	context.Profile = NewProfileCollector(context)
	return context
}

func (tctx (TraceContext)) ToXlog(discardType netdata.XlogDiscardType, elapsed int32) *netdata.XlogPack {
	xlog := netdata.NewXlogPack()
	xlog.Elapsed = elapsed
	xlog.Service = tctx.ServiceHash
	xlog.XType = tctx.XType
	xlog.Txid = tctx.Txid
	xlog.Gxid = tctx.Gxid
	xlog.Caller = tctx.Caller
	xlog.ThreadNameHash = SendHashedMessage(strconv.Itoa(int(tctx.Gxid)))
	xlog.SqlCount = tctx.SqlCount
	xlog.SqlTime = tctx.SqlTime
	xlog.Ipaddr = util.IpToBytes(tctx.RemoteIp)
	xlog.Userid = tctx.Userid
	xlog.HasDump = 0
	xlog.Error = tctx.Error

	xlog.DiscardType = discardType
	xlog.ProfileSize = int32(tctx.ProfileSize)
	xlog.ProfileCount = int32(tctx.ProfileCount)
	xlog.UserAgent = tctx.UserAgent
	xlog.Referer = tctx.Referer

	xlog.ApicallCount = tctx.ApicallCount
	xlog.ApicallTime = tctx.ApicallTime

	if tctx.Login != "" {
		xlog.Login = SendLogin(tctx.Login)
	}
	if tctx.Desc != "" {
		xlog.Desc = SendDesc(tctx.Desc)
	}
	xlog.Text1 = tctx.Text1
	xlog.Text2 = tctx.Text2
	xlog.Text4 = tctx.Text3
	xlog.Text4 = tctx.Text4
	xlog.Text5 = tctx.Text5

	return xlog
}
