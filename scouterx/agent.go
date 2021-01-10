package scouterx

import (
	"sync"
)

var agentOnce sync.Once

func Init() {
	agentOnce.Do(func() {
		//confDir := "./nil"
		//dir, err := os.Getwd()
		//if err == nil {
		//	confDir = dir
		//}
		//UdpSender = udpsender.GetInstance()
	})
}

