package udpclient

import (
	"container/list"
	"github.com/scouter-project/scouter-go-lib/common/constants/netcafeconstant"
	"github.com/scouter-project/scouter-go-lib/common/logger"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"github.com/scouter-project/scouter-go-lib/common/util"
	"github.com/scouter-project/scouter-go-lib/common/util/keygen"
	"net"
	"strconv"
)

//UDPClient is a upp socket client
type UDPClient struct {
	Conn          *net.UDPConn
	remoteAddress string
	remotePort    int
	udpMaxBytes   int
}

//NewUDPClient returns new udpclient instance
func New(addr string , port int) *UDPClient {
	udpclient := &UDPClient{remoteAddress: addr, remotePort: port}
	udpclient.open()
	udpclient.udpMaxBytes = 60000
	return udpclient
}

func (udpClient *UDPClient) SetUDPMaxBytes (max int) {
	udpClient.udpMaxBytes = max
}


func (udpClient *UDPClient) open() error {
	if udpClient.Conn != nil {
		udpClient.close()
	}

	address := udpClient.remoteAddress + ":" + strconv.Itoa(udpClient.remotePort)
	s, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		logger.Error.Printf("can't initialize udp client. %s\n", err.Error())
		return err
	}
	udpClient.Conn, err = net.DialUDP("udp", nil, s)
	if err != nil {
		logger.Error.Printf("can't initialize udp client. %s\n", err.Error())
		return err
	}
	return nil
}

func (udpClient *UDPClient) close() {
	udpClient.Conn.Close()
}


func (udpClient *UDPClient) writeMTU(data []byte, packetSize int) bool {
	if udpClient.Conn == nil {
		return false
	}
	pkid := keygen.Next()
	total := len(data) / packetSize
	remainder := len(data) % packetSize
	if remainder > 0 {
		total++
	}
	var num int

	for num = 0; num < len(data) / packetSize; num++ {
		udpClient.writeMTUSub(pkid, total, num,  util.CopyArray(data, num*packetSize, packetSize))
	}
	if remainder > 0 {
		udpClient.writeMTUSub(pkid, total, num,  util.CopyArray(data, len(data)-remainder, remainder))
	}
	return true
}

func (udpClient *UDPClient) writeMTUSub(pkid int64, total int, num int, data []byte) {
	out := netdata.NewDataOutputX(nil)
	out.Write(netcafeconstant.CAFE_MTU)
	out.WriteInt32(0) // In deed, this value is objhash.. but i don't use objhash anymore.
	out.WriteInt64(pkid)
	out.WriteInt16(int16(total))
	out.WriteInt16(int16(num))
	out.WriteBlob(data)
	buff := out.Bytes()
	udpClient.Conn.Write(buff)

}

func (udpClient *UDPClient) sendBufferList(bufferCount int16, data []byte) {
	out := netdata.NewDataOutputX(nil)
	out.Write(netcafeconstant.CAFE_N)
	out.WriteInt16(bufferCount)
	out.Write(data)
	sendData := out.Bytes()
	udpClient.Conn.Write(sendData)

}

// Send will send data via udp socket
func (udpClient *UDPClient) WriteBuffer(buff []byte) bool {
	if udpClient.Conn == nil {
		return false
	}
	if len(buff) > udpClient.udpMaxBytes {
		return udpClient.writeMTU(buff, udpClient.udpMaxBytes)
	}
	out := netdata.NewDataOutputX(nil)
	out.Write(netcafeconstant.CAFE)
	out.Write(buff)

	udpClient.Conn.Write(out.Bytes())
	return true
}

func (udpClient *UDPClient) WriteBufferList(bufferList *list.List) bool {
	if udpClient.Conn == nil {
		return false
	}
	out := netdata.NewDataOutputX(nil)
	var outCount int16
	for buffer:= bufferList.Front(); buffer != nil; buffer = buffer.Next() {
		b := buffer.Value.([]byte)
		buffLen := len(b)
		if buffLen > udpClient.udpMaxBytes {
			udpClient.writeMTU(b, udpClient.udpMaxBytes)
		} else if buffLen+ int(out.GetWriteSize()) > udpClient.udpMaxBytes {
			udpClient.sendBufferList(outCount, out.Bytes())
			out = netdata.NewDataOutputX(nil)
			outCount = 1
			out.Write(b)
		} else {
			outCount++
			out.Write(b)
		}
	}
	if out.GetWriteSize() > 0 {
		udpClient.sendBufferList(outCount, out.Bytes())
	}

	return true
}

