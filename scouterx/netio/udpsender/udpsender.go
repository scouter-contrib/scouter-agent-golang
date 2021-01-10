package udpsender

import (
	"container/list"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio/udpclient"
	"sync"
	"time"

	"github.com/scouter-project/scouter-go-lib/common/factory/channelfactory"
	"github.com/scouter-project/scouter-go-lib/common/logger"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
)

var once sync.Once
var ac = conf.GetInstance()

type UDPSender struct {
	udpChannel chan []byte
	running    bool
	udpClient  *udpclient.UDPClient
}

var udpSender *UDPSender

func GetInstance() *UDPSender {
	once.Do(func() {
		udpSender = new(UDPSender)
		udpSender.udpChannel = channelfactory.GetUDPChannel()
		udpSender.running = true
		udpSender.udpClient = udpclient.New(ac.NetCollectorIP, ac.NetCollectorUDPPort)
		udpSender.udpClient.SetUDPMaxBytes(ac.UDPMaxBytes)
		udpSender.running = true
		go udpSender.run()
	})
	return udpSender
}

func (udpSender *UDPSender) AddPack(pack netdata.Pack) {
	bytes := netdata.NewDataOutputX(nil).WritePack(pack).Bytes()
	select {
	case udpSender.udpChannel <- bytes:
	default:
		logger.Warning.Println("udp channel is full.")
	}
}

func (udpSender *UDPSender) AddBuffer(buffer []byte) {
	select {
	case udpSender.udpChannel <- buffer:
	default:
		logger.Warning.Println("udp channel is full.")
	}
}

func (udpSender *UDPSender) SendPackDirect(pack netdata.Pack) {
	bytes := netdata.NewDataOutputX(nil).WritePack(pack).Bytes()
	go udpSender.udpClient.WriteBuffer(bytes)
}

func (udpSender *UDPSender) SendDirect(bufferList *list.List) {
	if bufferList == nil {
		return
	}
	go udpSender.sendList(bufferList)
}

func (udpSender *UDPSender) run() {
	var size int
	for udpSender.running {
		size = udpSender.getQueueSize()
		//logger.Info.Printf("queue size : %d", size)
		switch size {
		case 0:
			time.Sleep(100 * time.Millisecond)
		case 1:
			p := <-udpSender.udpChannel
			udpSender.udpClient.WriteBuffer(p)
		default:
			udpSender.send(size)
		}
	}
}


func (udpSender *UDPSender) getQueueSize() int {
	return len(udpSender.udpChannel)
}

func (udpSender *UDPSender) send(size int) {
	bufferList := list.New()
	bytes := 0
	for i := 0; i < size; i++ {
		buffer := <-udpSender.udpChannel
		if bytes+len(buffer) >= ac.UDPMaxBytes {
			udpSender.sendList(bufferList)
			bytes = 0
			bufferList.Init()
		}
		bytes += len(buffer)
		bufferList.PushBack(buffer)
	}
	udpSender.udpClient.WriteBufferList(bufferList)

}

func (udpSender *UDPSender) sendList(bufferList *list.List) {
	if bufferList.Len() == 0 {
		return
	} else if bufferList.Len() == 1 {
		udpSender.udpClient.WriteBuffer(bufferList.Front().Value.([]byte))
	} else {
		udpSender.udpClient.WriteBufferList(bufferList)
	}
}
