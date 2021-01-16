package netdata

import (
	"strconv"

	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
)

// ObjectPack2 has object info
type ObjectPack2 struct {
	SiteID  string
	ObjType string
	ObjHash int32
	ObjName string
	Address string
	Version string
	Alive   bool
	Wakeup  int64
	Family  int8
	Tags    *MapValue
}

// ObjectPack has object info
type ObjectPack struct {
	ObjType string
	ObjHash int32
	ObjName string
	Address string
	Version string
	Alive   bool
	Wakeup  int64
	Tags    *MapValue
}

// NewObjectPack returns new object pack instance
func NewObjectPack() *ObjectPack {
	pack := new(ObjectPack)
	pack.Tags = NewMapValue()
	pack.Alive = true

	return pack
}

// NewObjectPack2 returns new object pack instance
func NewObjectPack2() *ObjectPack2 {
	pack := new(ObjectPack2)
	pack.Tags = NewMapValue()
	pack.SiteID = "Default"
	pack.Alive = true
	return pack
}

func (objectPack *ObjectPack2) Write(out *DataOutputX) error {
	var err error
	_, err = out.WriteString(objectPack.SiteID)
	_, err = out.WriteString(objectPack.ObjType)
	_, err = out.WriteDecimal32(objectPack.ObjHash)
	_, err = out.WriteString(objectPack.ObjName)
	_, err = out.WriteString(objectPack.Address)
	_, err = out.WriteString(objectPack.Version)
	_, err = out.WriteBoolean(objectPack.Alive)
	_, err = out.WriteDecimal(objectPack.Wakeup)
	_, err = out.WriteInt8(objectPack.Family)
	_, err = out.WriteValue(objectPack.Tags)
	return err
}

func (objectPack *ObjectPack) Write(out *DataOutputX) error {
	var err error
	_, err = out.WriteString(objectPack.ObjType)
	_, err = out.WriteDecimal32(objectPack.ObjHash)
	_, err = out.WriteString(objectPack.ObjName)
	_, err = out.WriteString(objectPack.Address)
	_, err = out.WriteString(objectPack.Version)
	_, err = out.WriteBoolean(objectPack.Alive)
	_, err = out.WriteDecimal(objectPack.Wakeup)
	_, err = out.WriteValue(objectPack.Tags)
	return err
}

func (objectPack *ObjectPack2) Read(in *DataInputX) (Pack, error){
	var err error
	objectPack.SiteID, err = in.ReadString()
	objectPack.ObjType, err = in.ReadString()
	objHash, err := in.ReadDecimal()
	objectPack.ObjHash = int32(objHash)
	objectPack.ObjName, err = in.ReadString()
	objectPack.Address, err = in.ReadString()
	objectPack.Version, err = in.ReadString()
	objectPack.Alive, err = in.ReadBoolean()
	objectPack.Wakeup, err = in.ReadDecimal()
	objectPack.Family, err = in.ReadInt8()
	value, err := in.ReadValue()
	objectPack.Tags = value.(*MapValue)

	return objectPack, err
}
func (objectPack *ObjectPack) Read(in *DataInputX) (Pack, error) {
	var err error
	objectPack.ObjType, err = in.ReadString()
	objHash, err := in.ReadDecimal()
	objectPack.ObjHash = int32(objHash)
	objectPack.ObjName, err = in.ReadString()
	objectPack.Address, err = in.ReadString()
	objectPack.Version, err = in.ReadString()
	objectPack.Alive, err = in.ReadBoolean()
	objectPack.Wakeup, err = in.ReadDecimal()
	value, err := in.ReadValue()
	objectPack.Tags = value.(*MapValue)

	return objectPack, err
}

// ToString returns objectpack2 info
func (objectPack *ObjectPack2) ToString() string {
	var str string
	str += "object siteID: " + objectPack.SiteID
	str += " name: " + objectPack.ObjName
	str += " type: " + objectPack.ObjType
	str += " hash: " + strconv.FormatInt(int64(objectPack.ObjHash), 10)
	str += " version: " + objectPack.Version
	str += " alive: " + strconv.FormatBool(objectPack.Alive)
	str += " familly: " + strconv.FormatInt(int64(objectPack.Family), 10)
	str += " tags: " + objectPack.Tags.ToString()
	return str
}

// ToString returns objectpack2 info
func (objectPack *ObjectPack) ToString() string {
	var str string
	str += "object name: " + objectPack.ObjName
	str += " type: " + objectPack.ObjType
	str += " hash: " + strconv.FormatInt(int64(objectPack.ObjHash), 10)
	str += " version: " + objectPack.Version
	str += " alive: " + strconv.FormatBool(objectPack.Alive)
	str += " tags: " + objectPack.Tags.ToString()
	return str
}

func (objectPack *ObjectPack) SetStatus(status int) {
	objectPack.Tags.Put("status", status)
}

func (objectPack *ObjectPack2) SetStatus(status int) {
	objectPack.Tags.Put("status", status)
}

//GetPackType returns pack type
func (objectPack *ObjectPack2) GetPackType() byte {
	return packconstants.OBJECT
}

//GetPackType returns pack type
func (objectPack *ObjectPack) GetPackType() byte {
	return packconstants.OBJECT
}
