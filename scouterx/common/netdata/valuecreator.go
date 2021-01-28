package netdata

import (
	valueconstants "github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/valueconstant"
)

//CreateValue return Value instance
func CreateValue(valueType byte) Value {
	switch valueType {
	case valueconstants.NULL:
		return NewNilValue()
	case valueconstants.BOOLEAN:
		return NewBooleanEmptyValue()
	case valueconstants.DECIMAL:
		return NewDecimalEmptyValue()
	case valueconstants.FLOAT:
		return NewFloatEmptyValue()
	case valueconstants.DOUBLE:
		return NewDoubleEmptyValue()
	case valueconstants.TEXT:
		return NewTextEmptyValue()
	case valueconstants.TEXT_HASH:
		return NewTextHashEmptyValue()
	case valueconstants.BLOB:
		return NewBlobEmptyValue()
	case valueconstants.IP4ADDR:
		return NewBlobEmptyValue()
	case valueconstants.LIST:
		return NewListValue()
	case valueconstants.MAP:
		return NewMapValue()
	default:
		return nil
	}
}
