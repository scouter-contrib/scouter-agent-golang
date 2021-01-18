package tcpclient

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/netcafeconstant"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/tcpflag"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/logger"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"net"
	"strconv"
	"sync"
	"time"
)

var ac = conf.GetInstance()
var client *TCPClient
var onceTcpClient sync.Once

type TCPClient struct {
	conn              *net.TCPConn
	host              string
	port              int
	connectionTimeout int
	soTimeout         int
	localAddr         string
	objHash			  int32
}

func GetTcpClient() *TCPClient {
	onceTcpClient.Do(func() {
		client = new(TCPClient)
	})
	return client
}

func (c *TCPClient) Close() {
	c.conn.Close()
}

func (c *TCPClient) Prepare() bool {
	c.host = ac.NetCollectorIP
	c.port = ac.NetCollectorTCPPort
	c.connectionTimeout = ac.NetCollectorTcpConnectionTimeoutMs
	c.soTimeout = ac.NetCollectorTcpSoTimeoutMs
	c.objHash = ac.ObjHash

	logger.Info.Printf("[scouter]tcp prepare %s, %d", c.host, c.port)
	conn0, err := net.DialTimeout("tcp", c.host+":"+strconv.Itoa(c.port), time.Duration(c.connectionTimeout)*time.Millisecond)
	if err != nil {
		logger.Error.Printf("[scouter][err]%v\n", err)
		c.conn = nil
		return false
	}
	tcpConn, ok := conn0.(*net.TCPConn)
	if !ok {
		c.conn = nil
		conn0.Close()
		return false
	}

	c.conn = tcpConn
	c.localAddr = "127.0.0.1" //TODO

	return true
}

func (c *TCPClient) Process() error {
	if c.conn == nil {
		return nil
	}
	out := netdata.NewDataOutputX(c.conn)
	defer c.conn.Close()

	serverAddr := ac.NetCollectorIP
	serverPort := ac.NetCollectorTCPPort

	_, err := out.WriteInt32(int32(netcafeconstant.TCP_AGENT_V2))
	_, err = out.WriteInt32(c.objHash)
	if err != nil {
		return err
	}

	for c.objHash == ac.ObjHash && serverAddr == ac.NetCollectorIP && serverPort == ac.NetCollectorTCPPort {
		err := c.conn.SetReadDeadline(time.Now().Add(time.Duration(ac.NetCollectorTcpSoTimeoutMs) * time.Millisecond))
		if err != nil {
			return err
		}
		in := netdata.NewDataInputX(c.conn)
		buff, err := in.ReadIntBytes()
		if err != nil {
			return err
		}

		in0 := netdata.NewDataInputX(buff)
		cmd, err := in0.ReadString()
		parameter, err := in0.ReadPack()
		if err != nil {
			return err
		}

		out0 := netdata.NewDataOutputX(nil)
		res := handle(cmd, parameter, in, out0)
		if res != nil {
			out0.WriteUInt8(tcpflag.HasNEXT)
			var pack *netdata.DataOutputX
			pack, err = netdata.NewDataOutputX(nil).WritePack(res)
			if err != nil {
				return err
			}
			packBytes := pack.Bytes()
			out0.WriteIntBytes(packBytes)
		}
		out0.WriteUInt8(tcpflag.NoNEXT)
		err = out.Write(out0.Bytes())
		if err != nil {
			return err
		}
	}
	return err
}
