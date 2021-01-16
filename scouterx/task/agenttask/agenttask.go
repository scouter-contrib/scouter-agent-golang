package agenttask

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"time"
)

func Start() {
	go start()
}

func start() {
	for {
		time.Sleep(2000 * time.Millisecond)
		SendObjPack()
	}
}

func SendObjPack() {
	ac := conf.GetInstance()
	objName := ac.ObjName
	objHash := netio.SendObjName(objName)

	objPack := netdata.NewObjectPack()
	objPack.ObjName = objName
	objPack.ObjHash = objHash
	objPack.ObjType = "go"
	objPack.Version = "0.0.0"
	netio.SendPackDirect(objPack)
}
