package strace

import (
	"context"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"github.com/scouter-project/scouter-go-lib/common/util"
	"testing"
	"time"
)

func TestXlog(t *testing.T) {
	RegisterObj()
	service := "/test-service/witprofile/0"

	func() {
		ctx := StartService(nil, service, "10.10.10.10")
		defer EndService(ctx)

		time.Sleep(time.Duration(100) * time.Millisecond)
		testMethod1(ctx)
		testMethod1(ctx)
		testMethod1(ctx)
		testMethod1(ctx)
		testMethod1(ctx)
		time.Sleep(time.Duration(200) * time.Millisecond)
	}()
	time.Sleep(time.Duration(100) * time.Millisecond)
}

func TestXlogWithGo(t *testing.T) {
	RegisterObj()
	service := "/test-service/withgo/0"

	func() {
		ctx := StartService(nil, service, "10.10.10.10")
		defer EndService(ctx)
		testMethod1(ctx)

		GoWithTrace(ctx, "testMethod4Go()", func(cascadeGoCtx context.Context) {
			testMethod4Go(cascadeGoCtx)
		})
		time.Sleep(time.Duration(30) * time.Millisecond)

		testMethod1(ctx)
		time.Sleep(time.Duration(200) * time.Millisecond)
	}()
	time.Sleep(time.Duration(100) * time.Millisecond)
}

func RegisterObj() *netdata.ObjectPack {
	objPack := netdata.NewObjectPack()
	objPack.ObjName = "node-testcase0"
	objPack.ObjHash = util.HashString(objPack.ObjName)
	objPack.ObjType = "java"
	netio.SendPackDirect(objPack)
	conf.GetInstance().ObjHash = objPack.ObjHash

	return objPack
}

func testMethod1(ctx context.Context) {
	methodStep := StartMethod(ctx)
	defer EndMethod(ctx, methodStep)
	time.Sleep(time.Duration(150) * time.Millisecond)
}

func testMethod4Go(ctx context.Context) {
	methodStep := StartMethod(ctx)
	defer EndMethod(ctx, methodStep)
	time.Sleep(time.Duration(50) * time.Millisecond)
}
