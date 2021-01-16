package netdata

import (
	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
	"strconv"
)

const (
	XTYPE_WEB_SERVICE XlogType = 0
	XTYPE_APP_SERVICE XlogType = 1
	XTYPE_BACK_THREAD XlogType = 2
	XTYPE_BACK_THREAD2 XlogType = 4
	XTYPE_ZIPKIN_SPAN XlogType = 5
	XTYPE_UNKNOWN XlogType = 99

	XLOG_DISCARD_NONE XlogDiscardType = 1
	XLOG_DISCARD_ALL XlogDiscardType = 2
	XLOG_DISCARD_PROFILE XlogDiscardType = 3
	XLOG_DISCARD_ALL_FORCE XlogDiscardType = 4
	XLOG_DISCARD_PROFILE_FORCE XlogDiscardType = 5
)

type XlogType uint8
type XlogDiscardType uint8

// XlogPack has xlog info
type XlogPack struct {
	EndTime int64
	ObjHash int32

	Service int32

	Txid int64
	ThreadNameHash int32
	Caller int64
	Gxid int64

	Elapsed int32

	Error int32
	Cpu int32
	SqlCount int32
	SqlTime int32
	Ipaddr []byte
	Kbytes int32
	Status int32
	Userid int64

	UserAgent int32
	Referer int32
	Group int32

	ApicallCount int32
	ApicallTime int32

	CountryCode string
	City int32

	XType XlogType

	Login int32
	Desc int32

	WebHash int32
	WebTime int32

	HasDump byte

	Text1 string
	Text2 string

	QueuingHostHash int32
	QueuingTime int32
	Queuing2ndHostHash int32
	Queuing2ndTime int32

	Text3 string
	Text4 string
	Text5 string

	ProfileCount int32
	B3Mode bool
	ProfileSize int32
	DiscardType XlogDiscardType
	IgnoreGlobalConsequentSampling bool
}

func NewXlogPack() *XlogPack {
	pack := new(XlogPack)
	return pack
}

func (p *XlogPack) IsDriving() bool {
	return (p.Gxid == p.Txid) || p.Gxid == 0
}

func (p *XlogPack) Write(out *DataOutputX) error {
	o := NewDataOutputX(nil)

	o.WriteDecimal(p.EndTime)
	o.WriteDecimal32(p.ObjHash)
	o.WriteDecimal32(p.Service)

	o.WriteInt64(p.Txid)
	o.WriteInt64(p.Caller)
	o.WriteInt64(p.Gxid)
	o.WriteDecimal32(p.Elapsed)

	o.WriteDecimal32(p.Error)

	o.WriteDecimal32(p.Cpu)

	o.WriteDecimal32(p.SqlCount)

	o.WriteDecimal32(p.SqlTime)

	o.WriteBlob(p.Ipaddr)
	o.WriteDecimal32(p.Kbytes)
	o.WriteDecimal32(p.Status)
	o.WriteDecimal(p.Userid)
	o.WriteDecimal32(p.UserAgent)
	o.WriteDecimal32(p.Referer)
	o.WriteDecimal32(p.Group)
	o.WriteDecimal32(p.ApicallCount)
	o.WriteDecimal32(p.ApicallTime)
	o.WriteString(p.CountryCode)
	o.WriteDecimal32(p.City)
	o.WriteUInt8(uint8(p.XType))
	o.WriteDecimal32(p.Login)
	o.WriteDecimal32(p.Desc)
	o.WriteDecimal32(p.WebHash)
	o.WriteDecimal32(p.WebTime)
	o.WriteUInt8(p.HasDump)
	o.WriteDecimal32(p.ThreadNameHash)
	o.WriteString(p.Text1)
	o.WriteString(p.Text2)
	o.WriteDecimal32(p.QueuingHostHash)
	o.WriteDecimal32(p.QueuingTime)
	o.WriteDecimal32(p.Queuing2ndHostHash)
	o.WriteDecimal32(p.Queuing2ndTime)
	o.WriteString(p.Text3)
	o.WriteString(p.Text4)
	o.WriteString(p.Text5)
	o.WriteDecimal32(p.ProfileCount)
	o.WriteBoolean(p.B3Mode)
	o.WriteDecimal32(p.ProfileSize)
	o.WriteUInt8(uint8(p.DiscardType))
	o.WriteBoolean(p.IgnoreGlobalConsequentSampling)

	err := out.WriteBlob(o.Bytes())
	return err
}

func (p *XlogPack) Read(in *DataInputX) (Pack, error) {
	//TODO not yet implemented
	return p, nil
}

func (pack *XlogPack) ToString() string {
	var str string
	str += "XLOG: "
	str += " objHash: " + strconv.FormatInt(int64(pack.ObjHash), 10)
	str += " service: " + strconv.FormatInt(int64(pack.Service), 10)
	str += " txid: " + strconv.FormatInt(pack.Txid, 10)
	str += " elapsed: " + strconv.FormatInt(int64(pack.Elapsed), 10)
	str += " error: " + strconv.FormatInt(int64(pack.Error), 10)
	return str
}

//GetPackType returns pack type
func (pack *XlogPack) GetPackType() byte {
	return packconstants.XLOG
}


