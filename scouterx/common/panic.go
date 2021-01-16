package common

import (
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/dump"
)

var fSendPanic func(string)

func RegisterSendPanicF(f func(string)) {
	fSendPanic = f
}

func ReportScouterPanic() {
	if r := recover(); r != nil {
		stackTrace := dump.StackTrace(100)
		fmt.Printf("[scouter][panic]%+v\n%s", r, stackTrace)
		fSendPanic(fmt.Sprintf("%+v - %s", r, stackTrace))
	}
}

