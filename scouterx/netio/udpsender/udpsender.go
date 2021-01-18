package udpsender

import (
	"container/list"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio/udpclient"
	"sync"
	"time"

	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/factory/channelfactory"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/logger"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
)

var once sync.Once
var ac = conf.GetInstance()

var udpObjHash int32
var serverAddr string
var udpServerPort int
var udpMaxBytes int

type UDPSender struct {
	udpChannel chan []byte
	running    bool
	udpClient  *udpclient.UDPClient
}

var udpSender *UDPSender

func GetInstance() *UDPSender {
	once.Do(func() {
		serverAddr = ac.NetCollectorIP
		udpServerPort = ac.NetCollectorUDPPort
		udpMaxBytes = ac.UDPMaxBytes

		udpSender = new(UDPSender)
		udpSender.udpChannel = channelfactory.GetUDPChannel()
		udpSender.running = true
		udpSender.udpClient = udpclient.New(serverAddr, udpServerPort)
		udpSender.udpClient.SetUDPMaxBytes(udpMaxBytes)
		go udpSender.run()
		go reloadUdpSender()
	})
	return udpSender
}

func reloadUdpSender() {
	for {
		time.Sleep(1000)
		if serverAddr != ac.NetCollectorIP || udpServerPort != ac.NetCollectorUDPPort || udpMaxBytes != ac.UDPMaxBytes {
			serverAddr = ac.NetCollectorIP
			udpServerPort = ac.NetCollectorUDPPort
			udpMaxBytes = ac.UDPMaxBytes

			udpSender.running = true
			prevClient := udpSender.udpClient
			udpClient := udpclient.New(serverAddr, udpServerPort)
			udpClient.SetUDPMaxBytes(udpMaxBytes)
			udpSender.udpClient = udpClient

			if prevClient.Conn != nil {
				prevClient.Conn.Close()
			}
		}
	}
}

func (udpSender *UDPSender) AddPack(pack netdata.Pack) {
	writePack, _ := netdata.NewDataOutputX(nil).WritePack(pack)
	bytes := writePack.Bytes()
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
	if ac.TraceObjSend {
		if p, ok := pack.(*netdata.ObjectPack); ok {
			logger.Info.Printf("[scouter] SendPackDirect[ObjPack], to:%s, pack:%s", udpSender.udpClient.Conn.RemoteAddr(), p.ToString())
		}
	}
	writePack, _ := netdata.NewDataOutputX(nil).WritePack(pack)
	bytes := writePack.Bytes()
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
