package scouterx

import (
	"context"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/strace"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestScouterAgent(T *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	Init()
	go loadTest()
	go loadTest()
	go loadTest()
	go loadTest()
	go loadTest()

	wg.Wait()
}

func loadTest() {
	for {
		randomSleeps()
	}
}


func randomSleeps() {
	ctx := context.Background()
	ctx = strace.StartService(ctx, "randomSleeps", "")
	defer strace.EndService(ctx)

	randomSleep(ctx, 1500)
	strace.GoWithTrace(ctx, "myGoFunc()", func (cascadeGoCtx context.Context) {
		randomSleep(cascadeGoCtx, 500)
	})
	randomSleep(ctx, 800)
}


func randomSleep(ctx context.Context, ms int) {
	step := strace.StartMethod(ctx)
	defer strace.EndMethod(ctx, step)

	rand := rand.Intn(ms)
	sleepFunc(ctx, rand)
}

func sleepFunc(ctx context.Context, ms int) {
	step := strace.StartMethod(ctx)
	defer strace.EndMethod(ctx, step)

	time.Sleep(time.Duration(ms) * time.Millisecond)
}
