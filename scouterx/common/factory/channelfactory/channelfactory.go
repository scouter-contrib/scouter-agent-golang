package channelfactory

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"sync"
)

var once sync.Once
var udpChannel chan []byte

//GetUDPChannel returns  channel which stores pack data.
func GetUDPChannel() chan []byte {
	once.Do(func() {
		udpChannel = make(chan []byte, conf.GetInstance().SendQueueSize)
	})
	return udpChannel
}
