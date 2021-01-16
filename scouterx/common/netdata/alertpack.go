package netdata

import packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"

type AlertLevel uint8

const AlertInfo AlertLevel = 0
const AlertWarn AlertLevel = 1
const AlertError AlertLevel = 2
const AlertFatal AlertLevel = 3

// AlertPack has text info
type AlertPack struct {
	Time int64
	ObjType string
	ObjHash int32
	Level AlertLevel
	Title string
	Message string
	Tags *MapValue
}

func NewAlertPack() *AlertPack {
	pack := new(AlertPack)
	pack.Tags = NewMapValue()
	return pack
}

func (p *AlertPack) Write(out *DataOutputX) error {
	_, err := out.WriteInt64(p.Time)
	_, err = out.WriteUInt8(uint8(p.Level))
	_, err = out.WriteString(p.ObjType)
	_, err = out.WriteInt32(p.ObjHash)
	_, err = out.WriteString(p.Title)
	_, err = out.WriteString(p.Message)
	_, err = out.WriteValue(p.Tags)

	return err
}

func (p *AlertPack) Read(in *DataInputX) (Pack, error) {
	//TODO not yet implemented
	return p, nil
}

func (pack *AlertPack) ToString() string {
	var str string
	str += "AlertPack:"
	str += " title: " + pack.Title
	str += " message: " + pack.Message
	return str
}

//GetPackType returns pack type
func (pack *AlertPack) GetPackType() byte {
	return packconstants.ALERT
}


