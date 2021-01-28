package netdata

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

// DataInputX is a byte buffer read stream struct
type DataInputX struct {
	offset int32
	reader io.Reader
}

// NewDataInputX returns DataInputX instace
func NewDataInputX(any interface{}) *DataInputX {
	in := new(DataInputX)
	switch v:= any.(type) {
	case []byte :
		in.reader = bytes.NewBuffer(v)
	case io.Reader:
		in.reader = v
	default:
		in.reader = nil
	}
	return in
}



// ReadInt8 returns int8 value
func (in *DataInputX) ReadInt8() (int8, error) {
	in.offset++
	var value int8
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}
func (in *DataInputX) ReadUInt8() (uint8, error) {
	in.offset++
	var value uint8
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}


// ReadInt16 returns int16 value
func (in *DataInputX) ReadInt16() (int16, error) {
	in.offset += 2
	var value int16
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}

// ReadInt32 returns int16 value
func (in *DataInputX) ReadInt32() (int32, error) {
	in.offset += 4
	var value int32
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}

// ReadInt64 returns int16 value
func (in *DataInputX) ReadInt64() (int64, error) {
	in.offset +=8
	var value int64
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}

// ReadFloat32 returns float32 value
func (in *DataInputX) ReadFloat32() (float32, error) {
	in.offset +=4
	var value float32
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}

func (in *DataInputX) ReadFloat64() (float64, error) {
	in.offset +=8
	var value float64
	err := binary.Read(in.reader, binary.BigEndian, &value)
	return value, err
}

// ReadString returns string value
func (in *DataInputX) ReadString() (string, error) {
	bytes, err := in.ReadBlob()
	return string(bytes), err

}

func (in *DataInputX) ReadBlob() ([]byte, error) {
	baseLen, err := in.ReadUInt8()
	var length int32
	switch baseLen {
	case 255:
		var v int16
		v, err = in.ReadInt16()
		length = int32(v)
	case 254:
		length, err = in.ReadInt32()
	case 0:
		return []byte{}, err
	default:
		length = int32(baseLen)
	}
	in.offset += length
	val := make([]byte, length)
	_, err = in.reader.Read(val)
	if err != nil {
		val = []byte{}
	}
	return val, err
}

// ReadDecimal returns number value
func (in *DataInputX) ReadDecimal() (int64, error) {
	len, err := in.ReadUInt8()
	switch len {
	case 0:
		return 0, err
	case 1:
		var v int8
		v, err = in.ReadInt8()
		return int64(v), err
	case 2:
		var v int16
		v, err = in.ReadInt16()
		return int64(v), err
	case 4:
		var v int32
		v, err = in.ReadInt32()
		return int64(v), err
	case 8:
		var v int64
		v, err = in.ReadInt64()
		return v, err
	default:
		return 0, err
	}
}

// ReadBoolean reads bool value
func (in *DataInputX) ReadBoolean() (bool, error) {
	value, err := in.ReadInt8()
	if value == 0 {
		return false, err
	} else {
		return true, err
	}
}

// ReadValue reads value from datainputx
func (in *DataInputX) ReadValue() (Value, error) {
	valueType, err := in.ReadInt8()
	value := CreateValue(byte(valueType))
	if value == nil {
		return NewNilValue(), errors.New("[scouter] Not defined value type:" + strconv.FormatUint(uint64(valueType), 10))
	}
	readValue, err := value.Read(in)
	return readValue, err

}

func (in *DataInputX) ReadIntBytes() ([]byte, error) {
	length, err := in.ReadInt32()
	b, err := in.Read(length)
	return b, err
	//val := make([]byte, length)
	//_, err = in.reader.Read(val)
	//if err != nil {
	//	val = []byte{}
	//}
	//return val, err
}

func (in *DataInputX) Read(len int32) ([]byte, error) {
	in.offset += len
	val := make([]byte, len)
	_, err := in.reader.Read(val)
	if err != nil {
		val = []byte{}
	}
	return val, err
}

func (in *DataInputX) ReadPack() (Pack, error) {
	packType, err := in.ReadInt8()
	if err != nil {
		return nil, err
	}
	pack := CreatePack(byte(packType))
	pack.Read(in)
	return pack, err

}

