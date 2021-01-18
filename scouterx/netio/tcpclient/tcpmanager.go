package tcpclient

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/logger"
	"time"
)

func StartTcp() {
	conf.GetInstance()
	go startTcp()
}

func startTcp() {
	for (true) {
		minSleep := 3000
		sleep := minSleep
		maxSleep := 60000

		func() {
			defer common.ReportScouterPanic()
			time.Sleep(time.Duration(1000) * time.Millisecond)
			client := GetTcpClient()
			if client.Prepare() {
				err := client.Process()
				if err != nil {
					logger.Error.Printf("[scouter][err]connection to collector: %v\n", err)
					time.Sleep(time.Duration(min(sleep, maxSleep)) * time.Millisecond)
					sleep = sleep * 2
					if sleep > maxSleep {
						sleep = maxSleep
					}
				}
			} else {
				time.Sleep(time.Duration(min(sleep, maxSleep)) * time.Millisecond)
				sleep = sleep * 2
				if sleep > maxSleep {
					sleep = maxSleep
				}
			}
		}()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
