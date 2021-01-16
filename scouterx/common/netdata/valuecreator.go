package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

//CreateValue return Value instacne
func CreateValue(valueType byte) Value {
	switch valueType {
	case valueconstants.NULL:
		return NewNilValue()
	case valueconstants.FLOAT:
		return NewFloatEmptyValue()
	case valueconstants.DECIMAL:
		return NewDecimalEmptyValue()
	case valueconstants.LIST:
		return NewListValue()
	case valueconstants.TEXT:
		return NewTextEmptyValue()
	case valueconstants.BOOLEAN:
		return NewBooleanEmptyValue()
	case valueconstants.MAP:
		return NewMapValue()
	default:
		return nil
	}
}
