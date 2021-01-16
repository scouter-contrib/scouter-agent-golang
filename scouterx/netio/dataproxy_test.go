package netio

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio/udpsender"
	util2 "github.com/scouter-contrib/scouter-agent-golang/scouterx/util"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	util "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util/keygen"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestSendXlog(t *testing.T) {
	service := "/test-service/0"

	RegisterObj()
	xlog := netdata.NewXlogPack()
	xlog.Txid = keygen.Next()
	xlog.XType = netdata.XTYPE_WEB_SERVICE
	xlog.Service = SendServiceName(service)
	xlog.Elapsed = 100
	xlog.DiscardType = netdata.XLOG_DISCARD_NONE
	xlog.Ipaddr = util2.IpToBytes("127.0.0.1")

	SendXlog(xlog)
}

func TestSendProfileAndXlog(t *testing.T) {
	RegisterObj()
	service := "/test-service/witprofile/0"
	serviceHash := SendServiceName(service)
	txid := keygen.Next()

	context := NewTraceContext()
	context.XType = netdata.XTYPE_WEB_SERVICE
	context.Txid = txid
	context.ServiceHash = serviceHash

	pushStep := netdata.NewMessageStep("test-push0", 0)
	context.Profile.Push(pushStep)
	context.Profile.Add(netdata.NewMessageStep("test-message0", 0))
	context.Profile.Pop(pushStep)
	context.Profile.Close(true)

	xlog := netdata.NewXlogPack()
	xlog.Txid = context.Txid
	xlog.XType = context.XType
	xlog.Service = context.ServiceHash
	xlog.Elapsed = 100
	xlog.DiscardType = netdata.XLOG_DISCARD_NONE
	xlog.Ipaddr = util2.IpToBytes("127.0.0.1")
	SendXlog(xlog)
}

func RegisterObj() *netdata.ObjectPack {
	objPack := netdata.NewObjectPack()
	objPack.ObjName = "node-testcase0"
	objPack.ObjHash = util.HashString(objPack.ObjName)
	objPack.ObjType = "java"
	SendPackDirect(objPack)
	conf.GetInstance().ObjHash = objPack.ObjHash

	return objPack
}

func TestSendObjectPack(t *testing.T) {

	os.Setenv("NET_COLLECTOR_IP", "127.0.0.1")
	os.Setenv("NET_COLLECTOR_UDP_PORT", "6002")
	os.Setenv("NET_COLLECTOR_TCP_PORT", "6002")
	/*
		os.Setenv("NET_COLLECTOR_IP", "127.0.0.1")
		os.Setenv("NET_COLLECTOR_UDP_PORT", "6100")
		os.Setenv("NET_COLLECTOR_TCP_PORT", "6100")
	*/

	os.Setenv("UDP_MAX_BYTES", "60000")

	sender := udpsender.GetInstance()
	for true {
		for i := 0; i < 200; i++ {
			objPack := netdata.NewObjectPack2()
			objPack.ObjName = "node" + strconv.Itoa(i)
			objPack.ObjHash = util.HashString(objPack.ObjName)
			objPack.ObjType = "host"
			objPack.Family = 2
			sender.AddPack(objPack)
			//time.Sleep(3 * time.Second)
		}
		time.Sleep(1 * time.Second)
	}
}
