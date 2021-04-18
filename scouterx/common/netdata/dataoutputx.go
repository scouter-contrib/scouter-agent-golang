package netdata

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

// A DataOutputX is a output stream which used write various kinds of data.
type DataOutputX struct {
	written int32 // the wrtten bytes.
	//buffer  *bytes.Buffer
	writer io.Writer
}

// NewDataOutputX returns DataOutputX object
func NewDataOutputX(any interface{}) *DataOutputX {
	out := new(DataOutputX)
	out.written = 0
	switch v := any.(type) {
	case nil:
		out.writer = new(bytes.Buffer)
	case io.Writer:
		out.writer = v
	default:
		out.writer = new(bytes.Buffer)
	}

	return out
}

func (out *DataOutputX) WriteInt32Array(values []int32) (*DataOutputX, error) {
	var err error
	if values == nil {
		_, err = out.WriteInt16(0)
	} else {
		_, err = out.WriteInt16(int16(len(values)))
		for _, v := range values {
			_, err = out.WriteInt32(v)
		}
	}
	return out, err
}

// WriteInt32 write int32 number to buffer.
func (out *DataOutputX) WriteInt32(value int32) (*DataOutputX, error){
	out.written += 4
	err := binary.Write(out.writer, binary.BigEndian, value)
	if err != nil {
		return out, err
	}
	return out, err
}

// WriteInt16 write int16 number to buffer.
func (out *DataOutputX) WriteInt16(value int16) (*DataOutputX, error) {
	out.written += 2
	err := binary.Write(out.writer, binary.BigEndian, value)
	if err != nil {
		return out, err
	}
	return out, err
}

// WriteInt64 write int64 number to buffer.
func (out *DataOutputX) WriteInt64(value int64) (*DataOutputX, error) {
	out.written += 8
	err := binary.Write(out.writer, binary.BigEndian, value)
	if err != nil {
		return out, err
	}
	return out, err
}

// WriteUInt64 write uint64 number to buffer.
func (out *DataOutputX) WriteUInt64(value uint64) (*DataOutputX, error) {
	out.written += 8
	err := binary.Write(out.writer, binary.BigEndian, value)
	return out, err
}

//WriteInt8 write int8 number to buffer
func (out *DataOutputX) WriteInt8(value int8) (*DataOutputX, error) {
	out.written++
	err := binary.Write(out.writer, binary.BigEndian, value)
	return out, err
}

//WriteUInt8 write int8 number to buffer
func (out *DataOutputX) WriteUInt8(value uint8) (*DataOutputX, error) {
	out.written++
	err := binary.Write(out.writer, binary.BigEndian, value)
	return out, err
}

// WriteFloat32 writes float32 value to buffer
func (out *DataOutputX) WriteFloat32(value float32) (*DataOutputX, error) {
	out.written += 4
	err := binary.Write(out.writer, binary.BigEndian, value)
	return out, err
}

// WriteFloat64 writes float64 value to buffer
func (out *DataOutputX) WriteFloat64(value float64) (*DataOutputX, error) {
	out.written += 8
	err := binary.Write(out.writer, binary.BigEndian, value)
	return out, err
}

// WriteDecimal writes number type value to buffer
func (out *DataOutputX) WriteDecimal(value int64) (*DataOutputX, error) {
	var err error
	if value == 0 {
		_, err = out.WriteInt8(0)
	} else if value >= math.MinInt8 && value <= math.MaxInt8 {
		_, err = out.WriteInt8(1)
		_, err = out.WriteInt8(int8(value))
	} else if value >= math.MinInt16 && value <= math.MaxInt16 {
		_, err = out.WriteInt8(2)
		_, err = out.WriteInt16(int16(value))
	} else if value >= math.MinInt32 && value <= math.MaxInt32 {
		_, err = out.WriteInt8(4)
		_, err = out.WriteInt32(int32(value))
	} else if value >= math.MinInt64 && value <= math.MaxInt64 {
		_, err = out.WriteInt8(8)
		_, err = out.WriteInt64(value)
	}
	return out, err
}

func (out *DataOutputX) WriteDecimal32(value int32) (*DataOutputX, error) {
	var err error
	if value == 0 {
		_, err = out.WriteInt8(0)
	} else if value >= math.MinInt8 && value <= math.MaxInt8 {
		_, err = out.WriteInt8(1)
		_, err = out.WriteInt8(int8(value))
	} else if value >= math.MinInt16 && value <= math.MaxInt16 {
		_, err = out.WriteInt8(2)
		_, err = out.WriteInt16(int16(value))
	} else if value >= math.MinInt32 && value <= math.MaxInt32 {
		_, err = out.WriteInt8(4)
		_, err = out.WriteInt32(int32(value))
	}
	return out, err
}

// WriteValue wtires value type to buffer
func (out *DataOutputX) WriteValue(value Value) (*DataOutputX, error) {
	if value == nil {
		value = NewNilValue()
	}
	_, err := out.WriteUInt8(uint8(value.GetValueType()))
	value.Write(out)
	return out, err
}

func (out *DataOutputX) WritePack(pack Pack) (*DataOutputX, error) {
	_, err := out.WriteUInt8(uint8(pack.GetPackType()))
	pack.Write(out)
	return out, err
}

func (out *DataOutputX) WriteStep(step Step) (*DataOutputX, error) {
	_, err := out.WriteUInt8(uint8(step.GetStepType()))
	err = step.Write(out)
	return out, err
}


// WriteString writes string value to buffer
func (out *DataOutputX) WriteString(value string) (*DataOutputX, error) {
	length := len(value)
	var err error
	if length == 0 {
		_, err = out.WriteInt8(0)
	} else if length > 100000 {
		err = out.WriteBlob([]byte("huge string omitted..."))
	} else {
		err = out.WriteBlob([]byte(value))
	}
	return out, err
}

func (out *DataOutputX) WriteIntBytes(value []byte) error {
	var err error
	_, err = out.WriteInt32(int32(len(value)))
	err = out.Write(value)
	return err
}

// WriteBlob writes byte array to buffer
func (out *DataOutputX) WriteBlob(value []byte) error {
	valueLen := len(value)
	var err error
	if valueLen == 0 {
		_, err = out.WriteUInt8(0)
	} else {
		if valueLen <= 253 {
			_, err = out.WriteUInt8(byte(valueLen))
			err = out.Write(value)
		} else if valueLen <= 65535 {
			_, err = out.WriteUInt8(255)
			_, err = out.WriteInt16(int16(valueLen))
			err = out.Write(value)
		} else {
			_, err = out.WriteUInt8(254)
			_, err = out.WriteInt32(int32(valueLen))
			err = out.Write(value)
		}
	}
	return err
}

func (out *DataOutputX) Write(value []byte) error {
	out.written += int32(len(value))
	_, err := out.writer.Write(value)
	return err
}

// Bytes returns buffer's bytes
func (out *DataOutputX) Bytes() []byte{
	if out == nil {
		return nil
	}
	switch v := out.writer.(type) {
	case *bytes.Buffer:
		return v.Bytes()
	default:
		return nil
	}

}

// WriteBoolean write bool valvue to buffer
func (out *DataOutputX) WriteBoolean(value bool) (*DataOutputX, error) {
	var err error
	if value == true {
		_, err = out.WriteInt8(1)
	} else {
		_, err = out.WriteInt8(0)
	}
	return out, err
}

// Size returns written size
func (out *DataOutputX) Size() int32 {
	return out.written
}

// GetWriteSize returns written size
func (out *DataOutputX) GetWriteSize() int32 {
	return out.written
}
