package netdata

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata/texttype"
	"strconv"

	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
)

// TextPack has text info
type TextPack struct {
	Xtype texttype.TextType
	Hash int32
	Text string
}

func NewTextPack() *TextPack {
	pack := new(TextPack)
	return pack
}

func (p *TextPack) Write(out *DataOutputX) error {
	_, err := out.WriteString(string(p.Xtype))
	_, err = out.WriteInt32(p.Hash)
	_, err = out.WriteString(p.Text)
	return err
}

func (p *TextPack) Read(in *DataInputX) (Pack, error) {
	//TODO not yet implemented
	return p, nil
}

func (pack *TextPack) ToString() string {
	var str string
	str += string(pack.Xtype)
	str += " hash: " + strconv.FormatInt(int64(pack.Hash), 10)
	str += " text: " + pack.Text
	return str
}

//GetPackType returns pack type
func (pack *TextPack) GetPackType() byte {
	return packconstants.TEXT
}


