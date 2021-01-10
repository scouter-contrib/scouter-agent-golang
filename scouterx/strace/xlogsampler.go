package strace

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-project/scouter-go-lib/common/netdata"
	"github.com/scouter-project/scouter-go-lib/common/util"
	"math/rand"
	"sync"
)

var xlogSamplerOnce sync.Once
var xlogSampler *XlogSampler

type XlogSampler struct {
	currentExcludeSamplingPattern   string
	currentDiscardServicePatterns   string
	currentSamplingServicePatterns  string
	currentSampling2ServicePatterns string
	currentSampling3ServicePatterns string

	excludeSamplingPatternMatcher util.CommaSeparatedChainedStrMatcher
	discardPatternMatcher         util.CommaSeparatedChainedStrMatcher
	samplingPatternMatcher        util.CommaSeparatedChainedStrMatcher
	sampling2PatternMatcher       util.CommaSeparatedChainedStrMatcher
	sampling3PatternMatcher       util.CommaSeparatedChainedStrMatcher
}

func GetXlogSampler() *XlogSampler {
	xlogSamplerOnce.Do(func() {
		xlogSampler = new(XlogSampler)
		xlogSampler.currentExcludeSamplingPattern = ac.XlogSamplingExcludePatterns
		xlogSampler.currentDiscardServicePatterns = ac.XlogDiscardServicePatterns
		xlogSampler.currentSamplingServicePatterns = ac.XlogPatternedSamplingServicePatterns
		xlogSampler.currentSampling2ServicePatterns = ac.XlogPatterned2SamplingServicePatterns
		xlogSampler.currentSampling3ServicePatterns = ac.XlogPatterned3SamplingServicePatterns

		xlogSampler.excludeSamplingPatternMatcher = util.NewCommaSeparatedChainedStrMatcher(xlogSampler.currentExcludeSamplingPattern)
		xlogSampler.discardPatternMatcher = util.NewCommaSeparatedChainedStrMatcher(xlogSampler.currentDiscardServicePatterns)
		xlogSampler.samplingPatternMatcher = util.NewCommaSeparatedChainedStrMatcher(xlogSampler.currentSamplingServicePatterns)
		xlogSampler.sampling2PatternMatcher = util.NewCommaSeparatedChainedStrMatcher(xlogSampler.currentSampling2ServicePatterns)
		xlogSampler.sampling3PatternMatcher = util.NewCommaSeparatedChainedStrMatcher(xlogSampler.currentSampling3ServicePatterns)

		conf.AddToConfObserver("XlogSampler", xlogSampler)
	})
	return xlogSampler
}

func (sampler *XlogSampler) Run() {
	ac := conf.GetInstance()

	if sampler.currentExcludeSamplingPattern != ac.XlogSamplingExcludePatterns {
		sampler.currentExcludeSamplingPattern = ac.XlogSamplingExcludePatterns
		sampler.excludeSamplingPatternMatcher = util.NewCommaSeparatedChainedStrMatcher(ac.XlogSamplingExcludePatterns)
	}
	if sampler.currentDiscardServicePatterns != ac.XlogDiscardServicePatterns {
		sampler.currentDiscardServicePatterns = ac.XlogDiscardServicePatterns
		sampler.discardPatternMatcher = util.NewCommaSeparatedChainedStrMatcher(ac.XlogDiscardServicePatterns)
	}
	if sampler.currentSamplingServicePatterns != ac.XlogPatternedSamplingServicePatterns {
		sampler.currentSamplingServicePatterns = ac.XlogPatternedSamplingServicePatterns
		sampler.samplingPatternMatcher = util.NewCommaSeparatedChainedStrMatcher(ac.XlogPatternedSamplingServicePatterns)
	}
	if sampler.currentSampling2ServicePatterns != ac.XlogPatterned2SamplingServicePatterns {
		sampler.currentSampling2ServicePatterns = ac.XlogPatterned2SamplingServicePatterns
		sampler.sampling2PatternMatcher = util.NewCommaSeparatedChainedStrMatcher(ac.XlogPatterned2SamplingServicePatterns)
	}
	if sampler.currentSampling3ServicePatterns != ac.XlogPatterned3SamplingServicePatterns {
		sampler.currentSampling3ServicePatterns = ac.XlogPatterned3SamplingServicePatterns
		sampler.sampling3PatternMatcher = util.NewCommaSeparatedChainedStrMatcher(ac.XlogPatterned3SamplingServicePatterns)
	}
}

func (sampler *XlogSampler) eEvaluateXLogDiscard(elapsed int32, serviceName string) netdata.XlogDiscardType {
	discardMode := netdata.XLOG_DISCARD_NONE

	if ac.XlogSamplingEnabled && sampler.isExcludeSamplingServicePattern(serviceName) {
		return netdata.XLOG_DISCARD_NONE
	}

	discarded := false
	if ac.XlogPatternedSamplingEnabled {
		discarded = sampler.isSamplingServicePattern(serviceName)
		if discarded {
			discardMode = sampler.samplingPatterned1(elapsed, discardMode)
		}
	}
	if !discarded && ac.XlogPatterned2SamplingEnabled {
		discarded = sampler.isSampling2ServicePattern(serviceName)
		if discarded {
			discardMode = sampler.samplingPatterned2(elapsed, discardMode)
		}
	}
	if !discarded && ac.XlogPatterned3SamplingEnabled {
		discarded = sampler.isSampling3ServicePattern(serviceName)
		if discarded {
			discardMode = sampler.samplingPatterned3(elapsed, discardMode)
		}
	}
	if !discarded && ac.XlogSamplingEnabled {
		discardMode = sampler.sampling4Elapsed(elapsed, discardMode)
	}

	return discardMode
}

func (sampler *XlogSampler) sampling4Elapsed(elapsed int32, discardMode netdata.XlogDiscardType) netdata.XlogDiscardType {
	if elapsed < ac.XlogSamplingStep1Ms {
		if rand.Intn(100) >= ac.XlogSamplingStep1RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogSamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogSamplingStep2Ms {
		if rand.Intn(100) >= ac.XlogSamplingStep2RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogSamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogSamplingStep3Ms {
		if rand.Intn(100) >= ac.XlogSamplingStep3RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogSamplingOnlyProfile)
		}
	} else {
		if rand.Intn(100) >= ac.XlogSamplingOverRatePct {
			discardMode = sampler.referToProfileMode(ac.XlogSamplingOnlyProfile)
		}
	}
	return discardMode
}

func (sampler *XlogSampler) samplingPatterned1(elapsed int32, discardMode netdata.XlogDiscardType) netdata.XlogDiscardType {
	if elapsed < ac.XlogPatternedSamplingStep1Ms {
		if rand.Intn(100) >= ac.XlogPatternedSamplingStep1RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatternedSamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogPatternedSamplingStep2Ms {
		if rand.Intn(100) >= ac.XlogPatternedSamplingStep2RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatternedSamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogPatternedSamplingStep3Ms {
		if rand.Intn(100) >= ac.XlogPatternedSamplingStep3RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatternedSamplingOnlyProfile)
		}
	} else {
		if rand.Intn(100) >= ac.XlogPatternedSamplingOverRatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatternedSamplingOnlyProfile)
		}
	}
	return discardMode
}

func (sampler *XlogSampler) samplingPatterned2(elapsed int32, discardMode netdata.XlogDiscardType) netdata.XlogDiscardType {
	if elapsed < ac.XlogPatterned2SamplingStep1Ms {
		if rand.Intn(100) >= ac.XlogPatterned2SamplingStep1RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned2SamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogPatterned2SamplingStep2Ms {
		if rand.Intn(100) >= ac.XlogPatterned2SamplingStep2RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned2SamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogPatterned2SamplingStep3Ms {
		if rand.Intn(100) >= ac.XlogPatterned2SamplingStep3RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned2SamplingOnlyProfile)
		}
	} else {
		if rand.Intn(100) >= ac.XlogPatterned2SamplingOverRatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned2SamplingOnlyProfile)
		}
	}
	return discardMode
}

func (sampler *XlogSampler) samplingPatterned3(elapsed int32, discardMode netdata.XlogDiscardType) netdata.XlogDiscardType {
	if elapsed < ac.XlogPatterned3SamplingStep1Ms {
		if rand.Intn(100) >= ac.XlogPatterned3SamplingStep1RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned3SamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogPatterned3SamplingStep2Ms {
		if rand.Intn(100) >= ac.XlogPatterned3SamplingStep2RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned3SamplingOnlyProfile)
		}
	} else if elapsed < ac.XlogPatterned3SamplingStep3Ms {
		if rand.Intn(100) >= ac.XlogPatterned3SamplingStep3RatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned3SamplingOnlyProfile)
		}
	} else {
		if rand.Intn(100) >= ac.XlogPatterned3SamplingOverRatePct {
			discardMode = sampler.referToProfileMode(ac.XlogPatterned3SamplingOnlyProfile)
		}
	}
	return discardMode
}

func (sampler *XlogSampler) referToProfileMode(onlyProfile bool) netdata.XlogDiscardType {
	if onlyProfile {
		return netdata.XLOG_DISCARD_PROFILE
	} else {
		return netdata.XLOG_DISCARD_ALL
	}
}

func (sampler *XlogSampler) isExcludeSamplingServicePattern(serviceName string) bool {
	if ac.XlogSamplingExcludePatterns == "" {
		return false
	}
	if sampler.excludeSamplingPatternMatcher.IsMatch(serviceName) {
		return true
	} else {
		return false
	}
}

func (sampler *XlogSampler) isDiscardServicePattern(serviceName string) bool {
	if ac.XlogDiscardServicePatterns == "" {
		return false
	}
	if sampler.discardPatternMatcher.IsMatch(serviceName) {
		return true
	} else {
		return false
	}
}

func (sampler *XlogSampler) isSamplingServicePattern(serviceName string) bool {
	if ac.XlogPatternedSamplingServicePatterns == "" {
		return false
	}
	if sampler.samplingPatternMatcher.IsMatch(serviceName) {
		return true
	} else {
		return false
	}
}

func (sampler *XlogSampler) isSampling2ServicePattern(serviceName string) bool {
	if ac.XlogPatterned2SamplingServicePatterns == "" {
		return false
	}
	if sampler.sampling2PatternMatcher.IsMatch(serviceName) {
		return true
	} else {
		return false
	}
}

func (sampler *XlogSampler) isSampling3ServicePattern(serviceName string) bool {
	if ac.XlogPatterned3SamplingServicePatterns == "" {
		return false
	}
	if sampler.sampling3PatternMatcher.IsMatch(serviceName) {
		return true
	} else {
		return false
	}
}
