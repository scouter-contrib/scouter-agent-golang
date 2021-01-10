package conf

import (
	"sync"
	"time"
)

var once sync.Once
var configureOnce sync.Once

type Configure struct {
	SendQueueSize int
	ObjHash int32
	ObjName string
	NetCollectorIP string
	NetCollectorUDPPort int
	NetCollectorTCPPort int
	UDPMaxBytes int

	ProfileStepMaxKeepInMemoryCount int
	ProfileStepMaxCount int

	ProfileHttpQuerystringEnabled bool
	ProfileHttpHeaderEnabled bool
	ProfileHttpHeaderKeys string

	TraceHttpClientIpHeaderKey string

	XlogDiscardServicePatterns string
	XlogDiscardServiceShowError bool

	XlogSamplingExcludePatterns string

	XlogSamplingEnabled bool
	XlogSamplingOnlyProfile bool
	XlogSamplingStep1Ms int32
	XlogSamplingStep1RatePct int
	XlogSamplingStep2Ms int32
	XlogSamplingStep2RatePct int
	XlogSamplingStep3Ms int32
	XlogSamplingStep3RatePct int
	XlogSamplingOverRatePct int

	XlogPatternedSamplingEnabled bool
	XlogPatternedSamplingServicePatterns string
	XlogPatternedSamplingOnlyProfile bool
	XlogPatternedSamplingStep1Ms int32
	XlogPatternedSamplingStep1RatePct int
	XlogPatternedSamplingStep2Ms int32
	XlogPatternedSamplingStep2RatePct int
	XlogPatternedSamplingStep3Ms int32
	XlogPatternedSamplingStep3RatePct int
	XlogPatternedSamplingOverRatePct int

	XlogPatterned2SamplingEnabled bool
	XlogPatterned2SamplingServicePatterns string
	XlogPatterned2SamplingOnlyProfile bool
	XlogPatterned2SamplingStep1Ms int32
	XlogPatterned2SamplingStep1RatePct int
	XlogPatterned2SamplingStep2Ms int32
	XlogPatterned2SamplingStep2RatePct int
	XlogPatterned2SamplingStep3Ms int32
	XlogPatterned2SamplingStep3RatePct int
	XlogPatterned2SamplingOverRatePct int

	XlogPatterned3SamplingEnabled bool
	XlogPatterned3SamplingServicePatterns string
	XlogPatterned3SamplingOnlyProfile bool
	XlogPatterned3SamplingStep1Ms int32
	XlogPatterned3SamplingStep1RatePct int
	XlogPatterned3SamplingStep2Ms int32
	XlogPatterned3SamplingStep2RatePct int
	XlogPatterned3SamplingStep3Ms int32
	XlogPatterned3SamplingStep3RatePct int
	XlogPatterned3SamplingOverRatePct int

}

var configure *Configure

func GetInstance() *Configure {
	once.Do(func() {
		configure = new(Configure)
		configure.NetCollectorIP = "127.0.0.1"
		configure.NetCollectorUDPPort = 6100
		configure.NetCollectorTCPPort = 6100
		configure.UDPMaxBytes = 60000
		configure.SendQueueSize = 3000
		configure.ProfileStepMaxKeepInMemoryCount = 2048
		configure.ProfileStepMaxCount = 1024
		configure.XlogDiscardServiceShowError = true

		configure.init()
		go configure.run()
	})
	return configure
}

func (conf *Configure) init() {
	//TODO read system prop
	conf.refresh()
}

func (conf *Configure) refresh() {
	//TODO


	//last
	confChangeNotify()
}

func (conf *Configure) run() {
	for true {
		time.Sleep(5 * time.Second)
		conf.refresh()
	}
}


