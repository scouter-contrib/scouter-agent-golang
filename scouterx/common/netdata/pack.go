package netdata

// Pack is a interface
type Pack interface {
	Write(out *DataOutputX) error
	Read(in *DataInputX) (Pack, error)
	ToString() string
	GetPackType() byte
}
