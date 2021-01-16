package netdata

// Value is interface for all value type
type Value interface {
	Write(out *DataOutputX) error
	Read(in *DataInputX) (Value, error)
	GetValueType() byte
	ToString() string
}
