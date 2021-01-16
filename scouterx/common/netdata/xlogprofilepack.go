package netdata

import (
	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
	"strconv"
)

type XlogProfilePack struct {
	Time int64
	ObjHash int32
	Service int32
	Txid int64
	Elapsed int32
	Profile []byte
}

func NewXlogProfilePack() *XlogProfilePack {
	pack := new(XlogProfilePack)
	return pack
}

func (p *XlogProfilePack) Write(out *DataOutputX) error {
	_, err := out.WriteDecimal(p.Time)
	_, err = out.WriteDecimal32(p.ObjHash)
	_, err = out.WriteDecimal32(p.Service)
	_, err = out.WriteInt64(p.Txid)
	err = out.WriteBlob(p.Profile)
	return err
}

func (p *XlogProfilePack) Read(in *DataInputX) (Pack, error) {
	//TODO not yet implemented
	return p, nil
}

func (pack *XlogProfilePack) ToString() string {
	var str string
	str += "Profile: "
	str += " objHash: " + strconv.FormatInt(int64(pack.ObjHash), 10)
	return str
}

func (pack *XlogProfilePack) GetPackType() byte {
	return packconstants.XLOG_PROFILE
}


