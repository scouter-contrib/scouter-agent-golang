package common

import (
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/dump"
	"runtime"
	"strconv"
	"strings"
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

func Goid() int {
	defer ReportScouterPanic()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

