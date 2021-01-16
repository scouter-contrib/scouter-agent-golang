package netdata

import (
	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
)

//PerfCounterPack has metric counter which has key/value type
type PerfCounterK8SPack struct {
	Time     int64
	SiteID	string
	Cluster string
	NameSpace string
	NodeName string
	DaemonSetName string
	DeploymentName string
	PodName string
	ContainerName string
	ObjName string
	Timetype int8
	MetricLevel int8
	Data     *MapValue
}

// NewPerfCounterPack returns PerfCounterPack instance
func NewPerfCounterK8SPack() *PerfCounterK8SPack {
	pack := new(PerfCounterK8SPack)
	pack.Data = NewMapValue()
	return pack
}

// Write will write PerfCounterPack to datoutputx
func (pack *PerfCounterK8SPack) Write(out *DataOutputX) error {
	_, err := out.WriteInt64(pack.Time)
	_, err = out.WriteString(pack.SiteID)
	_, err = out.WriteString(pack.Cluster)
	_, err = out.WriteString(pack.NameSpace)
	_, err = out.WriteString(pack.NodeName)
	_, err = out.WriteString(pack.DaemonSetName)
	_, err = out.WriteString(pack.DeploymentName)
	_, err = out.WriteString(pack.PodName)
	_, err = out.WriteString(pack.ContainerName)
	_, err = out.WriteString(pack.ObjName)
	_, err = out.WriteInt8(pack.Timetype)
	_, err = out.WriteInt8(pack.MetricLevel)
	_, err = out.WriteValue(pack.Data)
	return err
}

// Read will read PerfCounterPack from datainputx
func (pack *PerfCounterK8SPack) Read(in *DataInputX) (Pack, error) {
	var err error
	pack.Time, err = in.ReadInt64()
	pack.SiteID, err = in.ReadString()
	pack.Cluster, err = in.ReadString()
	pack.NameSpace, err = in.ReadString()
	pack.NodeName, err = in.ReadString()
	pack.DaemonSetName, err = in.ReadString()
	pack.DeploymentName, err = in.ReadString()
	pack.PodName, err = in.ReadString()
	pack.ContainerName, err = in.ReadString()
	pack.ObjName, err = in.ReadString()
	pack.Timetype, err = in.ReadInt8()
	pack.MetricLevel, err = in.ReadInt8()
	value, err := in.ReadValue()
	pack.Data = value.(*MapValue)
	return pack, err
}

// Put will put key/value to PerfCounterPack
func (pack *PerfCounterK8SPack) Put(key string, any interface{}) {
	switch v := any.(type) {
	case int32:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case int64:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case int:
		pack.Data.Put(key, NewDecimalValue(int64(v)))
	case float32:
		pack.Data.Put(key, NewFloatValue(v))
	case float64:
		pack.Data.Put(key, NewFloatValue(float32(v)))
	case string:
		pack.Data.Put(key, NewTextValue(v))
	case bool:
		pack.Data.Put(key, NewBooleanValue(v))

	}

}

// ToString returns converted perfcounterpack value
func (pack *PerfCounterK8SPack) ToString() string {
	str := pack.Data.ToString()
	return str
}

//GetPackType returns pack type
func (pack *PerfCounterK8SPack) GetPackType() byte {
	return packconstants.PERFCOUNTER_K8S
}
