package netdata

import (
	packconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/packconstant"
)

//CreateValue return Value instacne
func CreatePack(packType byte) Pack {
	switch packType {
	case packconstants.MAP:
		return NewMapPack()
	case packconstants.TEXT:
		return NewTextPack()
	case packconstants.PERFCOUNTER:
		return NewPerfCounterPack()
	case packconstants.OBJECT:
		return NewObjectPack2()
	case packconstants.ALERT:
		return NewAlertPack()
	default:
		return nil
	}
}
