package conf

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/magiconair/properties"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/logger"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var once sync.Once
var configureOnce sync.Once
var descSync sync.RWMutex

type ValueType int

const vtValue = 1
const vtNum = 2
const vtBool = 3
const vtCommaSeparatedValue = 4
const vtCommaColonSeparatedValue = 5

//var descMap = make(map[string]ConfigureDesc)
var descMap = linkedhashmap.New()

type ConfigureDesc struct {
	Key          string
	Value        string
	DefaultValue string
	Desc         string
	ValueType    ValueType
}

type Configure struct {
	sync.Mutex

	lastModified time.Time

	_trace bool
	TraceObjSend bool

	SendQueueSize int
	ObjHash int32
	ObjName string
	ObjType string
	ObjNameSimple string
	NetCollectorIP string
	NetCollectorUDPPort int
	NetCollectorTCPPort int
	NetCollectorTcpSoTimeoutMs int
	NetCollectorTcpConnectionTimeoutMs int
	UDPMaxBytes int

	StuckServiceBaseTimeMs int
	StuckServiceRemoveEnabled bool

	TraceActiveserviceYellowTime int
	TraceActiveserviceRedTime int

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
		logger.Init()
		configure = new(Configure)
		configure.NetCollectorIP = "127.0.0.1"
		configure.NetCollectorUDPPort = 6100
		configure.NetCollectorTCPPort = 6100
		configure.NetCollectorTcpSoTimeoutMs = 60000
		configure.NetCollectorTcpConnectionTimeoutMs = 3000
		configure.UDPMaxBytes = 60000

		configure.StuckServiceBaseTimeMs = 300000
		configure.StuckServiceRemoveEnabled = true

		configure.TraceActiveserviceYellowTime = 3000
		configure.TraceActiveserviceRedTime = 8000
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
	conf.Refresh()
}

func (conf *Configure) run() {
	for true {
		time.Sleep(5 * time.Second)
		conf.Refresh()
	}
}

func getConfFilePath() string {
	scouterConfFile := os.Getenv("SCOUTER_CONFIG")
	if scouterConfFile == "" {
		scouterConfFile = os.Getenv("scouter.config")
	}
	if scouterConfFile == "" {
		path := util.GetScouterPath()
		confPath := filepath.Join(path, "conf")
		util.MakeDir(confPath)
		scouterConfFile = filepath.Join(path, "conf", "scouter.conf")
	}

	return scouterConfFile
}

func (conf *Configure) SetTrace(mode bool) {
	conf._trace = mode
}

func (conf *Configure) IsTrace() bool {
	return conf._trace
}

func (conf *Configure) resetObjNameAndType(props *properties.Properties) {
	defaultName := "go1"
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	oldObjName := conf.ObjName
	newObjSimpleName := stringOf(props,"obj_name", defaultName, "object name")
	conf.ObjNameSimple = newObjSimpleName
	conf.ObjName = "/" + hostname + "/" + newObjSimpleName
	conf.ObjHash = util.HashString(conf.ObjName)

	conf.ObjType = stringOf(props,"obj_type", "golang", "object type(monitoring group)")
	if oldObjName != conf.ObjName {
		objChangeNotify()
	}
}

func (conf *Configure) Refresh() {
	conf.Lock()
	defer conf.Unlock()

	filePath := getConfFilePath()
	info, err := os.Stat(filePath)
	if err != nil {
		//TODO logging
		switch e := err.(type) {
		case *os.PathError:
			if e.Err != syscall.ENOENT {
				return
			}
		default:
			return
		}
	}
	if info == nil || info.ModTime().After(conf.lastModified) {
		props := properties.NewProperties()
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			props.Set(pair[0], pair[1])
		}

		fileProps, err := properties.LoadFile(filePath, properties.UTF8)
		if err != nil {
			switch e := err.(type) {
			case *os.PathError:
				if e.Err != syscall.ENOENT {
					return
				}
			default:
				return
			}
		}
		if fileProps != nil {
			props.Merge(fileProps)
		}
		props = props.FilterFunc(func (k, v string) bool {
			return v != ""
		})

		conf.addToConf(props)
		if info != nil {
			conf.lastModified = info.ModTime()
		}

		//last
		confChangeNotify()
	}
}

func GetConfigureDescMap() *linkedhashmap.Map {
	descSync.RLock()
	defer descSync.RUnlock()

	copy := linkedhashmap.New()
	it := descMap.Iterator()
	for it.Next() {
		copy.Put(it.Key(), it.Value())
	}
	return copy
}

func boolOf(props *properties.Properties, key string, defaultValue bool, desc string) bool {
	descSync.Lock()
	defer descSync.Unlock()

	value := props.GetBool(key, defaultValue)
	descMap.Put(key, ConfigureDesc{
		Key:          key,
		Value:        strconv.FormatBool(value),
		DefaultValue: strconv.FormatBool(defaultValue),
		Desc:         desc,
		ValueType:    vtBool,
	})
	return value
}

func intOf(props *properties.Properties, key string, defaultValue int, desc string) int {
	descSync.Lock()
	defer descSync.Unlock()

	value := props.GetInt(key, defaultValue)
	descMap.Put(key, ConfigureDesc{
		Key:          key,
		Value:        strconv.Itoa(value),
		DefaultValue: strconv.Itoa(defaultValue),
		Desc:         desc,
		ValueType:    vtNum,
	})
	return value
}

func stringOf(props *properties.Properties, key string, defaultValue string, desc string) string {
	descSync.Lock()
	defer descSync.Unlock()

	value := props.GetString(key, defaultValue)
	descMap.Put(key, ConfigureDesc{
		Key:          key,
		Value:        value,
		DefaultValue: defaultValue,
		Desc:         desc,
		ValueType:    vtValue,
	})
	return value
}

func stringOfType(props *properties.Properties, key string, defaultValue string, valueType ValueType, desc string) string {
	descSync.Lock()
	defer descSync.Unlock()

	value := props.GetString(key, defaultValue)
	descMap.Put(key, ConfigureDesc{
		Key:          key,
		Value:        value,
		DefaultValue: defaultValue,
		Desc:         desc,
		ValueType:    valueType,
	})
	return value
}

func (conf *Configure) addToConf(props *properties.Properties) {

	conf.resetObjNameAndType(props)

	conf._trace = props.GetBool("_trace", false)

	conf.TraceObjSend = boolOf(props, "trace_obj_send", false, "")
	conf.SendQueueSize = intOf(props, "send_queue_size", 3000, "")

	conf.NetCollectorIP = stringOf(props, "net_collector_ip", "127.0.0.1", "")
	conf.NetCollectorUDPPort = intOf(props, "net_collector_udp_port", 6100, "")
	conf.NetCollectorTCPPort = intOf(props, "net_collector_tcp_port", 6100, "")
	conf.NetCollectorTcpSoTimeoutMs = intOf(props, "net_collector_tcp_so_timeout_ms", 60000, "")
	conf.NetCollectorTcpConnectionTimeoutMs = intOf(props, "net_collector_tcp_connection_timeout_ms", 3000, "")
	conf.UDPMaxBytes = intOf(props, "udp_max_bytes", 60000, "")

	conf.StuckServiceBaseTimeMs = intOf(props, "stuck_service_base_time_ms", 300000, "")
	conf.StuckServiceRemoveEnabled = boolOf(props, "stuck_service_remove_enabled", true, "")

	conf.TraceActiveserviceYellowTime = intOf(props, "trace_activeservice_yellow_time", 3000, "")
	conf.TraceActiveserviceRedTime = intOf(props, "trace_activeservice_red_time", 8000, "")

	conf.ProfileStepMaxKeepInMemoryCount = intOf(props, "profile_step_max_keep_in_memory_count", 2048, "")
	conf.ProfileStepMaxCount = intOf(props, "profile_step_max_count", 1024, "")

	conf.ProfileHttpQuerystringEnabled = boolOf(props, "profile_http_querystring_enabled", false, "")
	conf.ProfileHttpHeaderEnabled = boolOf(props, "profile_http_header_enabled", false, "")
	conf.ProfileHttpHeaderKeys = stringOfType(props, "profile_http_header_keys", "", vtCommaSeparatedValue, "")

	conf.TraceHttpClientIpHeaderKey = stringOf(props, "trace_http_client_ip_header_key", "", "")

	conf.XlogDiscardServicePatterns = stringOfType(props, "xlog_discard_service_patterns", "", vtCommaSeparatedValue, "")
	conf.XlogDiscardServiceShowError = boolOf(props, "xlog_discard_service_show_error", true, "")

	conf.XlogSamplingExcludePatterns = stringOfType(props, "xlog_sampling_exclude_patterns", "", vtCommaSeparatedValue, "")

	conf.XlogSamplingEnabled = boolOf(props, "xlog_sampling_enabled", false, "")
	conf.XlogSamplingOnlyProfile = boolOf(props, "xlog_sampling_only_profile", false, "")
	conf.XlogSamplingStep1Ms = int32(intOf(props, "xlog_sampling_step1_ms", 100, ""))
	conf.XlogSamplingStep1RatePct = intOf(props, "xlog_sampling_step1_rate_pct", 2, "")
	conf.XlogSamplingStep2Ms = int32(intOf(props, "xlog_sampling_step2_ms", 500, ""))
	conf.XlogSamplingStep2RatePct = intOf(props, "xlog_sampling_step2_rate_pct", 7, "")
	conf.XlogSamplingStep3Ms = int32(intOf(props, "xlog_sampling_step3_ms", 1000, ""))
	conf.XlogSamplingStep3RatePct = intOf(props, "xlog_sampling_step3_rate_pct", 15, "")
	conf.XlogSamplingOverRatePct = intOf(props, "xlog_sampling_over_rate_pct", 3000, "")

	conf.XlogPatternedSamplingEnabled = boolOf(props, "xlog_patterned_sampling_enabled", false, "")
	conf.XlogPatternedSamplingServicePatterns = stringOfType(props, "xlog_patterned_sampling_service_patterns", "", vtCommaSeparatedValue, "")
	conf.XlogPatternedSamplingOnlyProfile = boolOf(props, "xlog_patterned_sampling_only_profile", false, "")
	conf.XlogPatternedSamplingStep1Ms = int32(intOf(props, "xlog_patterned_sampling_step1_ms", 100, ""))
	conf.XlogPatternedSamplingStep1RatePct = intOf(props, "xlog_patterned_sampling_step1_rate_pct", 2, "")
	conf.XlogPatternedSamplingStep2Ms = int32(intOf(props, "xlog_patterned_sampling_step2_ms", 500, ""))
	conf.XlogPatternedSamplingStep2RatePct = intOf(props, "xlog_patterned_sampling_step2_rate_pct", 7, "")
	conf.XlogPatternedSamplingStep3Ms = int32(intOf(props, "xlog_patterned_sampling_step3_ms", 1000, ""))
	conf.XlogPatternedSamplingStep3RatePct = intOf(props, "xlog_patterned_sampling_step3_rate_pct", 15, "")
	conf.XlogPatternedSamplingOverRatePct = intOf(props, "xlog_patterned_sampling_over_rate_pct", 3000, "")

	conf.XlogPatterned2SamplingEnabled = boolOf(props, "xlog_patterned2_sampling_enabled", false, "")
	conf.XlogPatterned2SamplingServicePatterns = stringOfType(props, "xlog_patterned2_sampling_service_patterns", "", vtCommaSeparatedValue, "")
	conf.XlogPatterned2SamplingOnlyProfile = boolOf(props, "xlog_patterned2_sampling_only_profile", false, "")
	conf.XlogPatterned2SamplingStep1Ms = int32(intOf(props, "xlog_patterned2_sampling_step1_ms", 100, ""))
	conf.XlogPatterned2SamplingStep1RatePct = intOf(props, "xlog_patterned2_sampling_step1_rate_pct", 2, "")
	conf.XlogPatterned2SamplingStep2Ms = int32(intOf(props, "xlog_patterned2_sampling_step2_ms", 500, ""))
	conf.XlogPatterned2SamplingStep2RatePct = intOf(props, "xlog_patterned2_sampling_step2_rate_pct", 7, "")
	conf.XlogPatterned2SamplingStep3Ms = int32(intOf(props, "xlog_patterned2_sampling_step3_ms", 1000, ""))
	conf.XlogPatterned2SamplingStep3RatePct = intOf(props, "xlog_patterned2_sampling_step3_rate_pct", 15, "")
	conf.XlogPatterned2SamplingOverRatePct = intOf(props, "xlog_patterned2_sampling_over_rate_pct", 3000, "")

	conf.XlogPatterned3SamplingEnabled = boolOf(props, "xlog_patterned3_sampling_enabled", false, "")
	conf.XlogPatterned3SamplingServicePatterns = stringOfType(props, "xlog_patterned3_sampling_service_patterns", "", vtCommaSeparatedValue, "")
	conf.XlogPatterned3SamplingOnlyProfile = boolOf(props, "xlog_patterned3_sampling_only_profile", false, "")
	conf.XlogPatterned3SamplingStep1Ms = int32(intOf(props, "xlog_patterned3_sampling_step1_ms", 100, ""))
	conf.XlogPatterned3SamplingStep1RatePct = intOf(props, "xlog_patterned3_sampling_step1_rate_pct", 2, "")
	conf.XlogPatterned3SamplingStep2Ms = int32(intOf(props, "xlog_patterned3_sampling_step2_ms", 500, ""))
	conf.XlogPatterned3SamplingStep2RatePct = intOf(props, "xlog_patterned3_sampling_step2_rate_pct", 7, "")
	conf.XlogPatterned3SamplingStep3Ms = int32(intOf(props, "xlog_patterned3_sampling_step3_ms", 1000, ""))
	conf.XlogPatterned3SamplingStep3RatePct = intOf(props, "xlog_patterned3_sampling_step3_rate_pct", 15, "")
	conf.XlogPatterned3SamplingOverRatePct = intOf(props, "xlog_patterned3_sampling_over_rate_pct", 3000, "")
}



