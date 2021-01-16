package udpsender

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	timeconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/timeconstant"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
)

func TestSendPerfPack(t *testing.T) {
	sender := GetInstance()
	perfPack := netdata.NewPerfCounterPack()
	perfPack.Put("abd", 123)
	perfPack.ObjName = "testObj"
	perfPack.Time = time.Now().Unix()
	perfPack.Timetype = timeconstants.REALTIME
	pack, _ := netdata.NewDataOutputX(nil).WritePack(perfPack)
	buffer := pack.Bytes()
	sender.AddBuffer(buffer)
	fmt.Printf("queue size: %d", sender.getQueueSize())
	for true {
		time.Sleep(1 * time.Second)
	}
}

func TestSendObjectPack1(t *testing.T) {
	sender := GetInstance()
	objPack := netdata.NewObjectPack()
	objPack.ObjName = "node-test0"
	objPack.ObjHash = util.HashString(objPack.ObjName)
	objPack.ObjType = "host"
	sender.AddPack(objPack)
}

func TestSendObjectPack(t *testing.T) {
	sender := GetInstance()
	for true {
		for i := 0; i < 200; i++ {
			objPack := netdata.NewObjectPack()
			objPack.ObjName = "node" + strconv.Itoa(i)
			objPack.ObjHash = util.HashString(objPack.ObjName)
			objPack.ObjType = "host"
			sender.AddPack(objPack)
			//time.Sleep(3 * time.Second)
		}
		time.Sleep(1 * time.Second)
	}
}
