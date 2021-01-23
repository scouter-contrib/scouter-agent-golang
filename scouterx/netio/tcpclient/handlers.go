package tcpclient

import (
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/conf"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/dump"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/netio"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/strace/tctxmanager"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/logger"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"strconv"
)

const OBJECT_RESET_CACHE string = "OBJECT_RESET_CACHE";
const TRIGGER_THREAD_DUMP string = "TRIGGER_THREAD_DUMP"
const OBJECT_DUMP_FILE_LIST string = "OBJECT_DUMP_FILE_LIST"
const OBJECT_DUMP_FILE_DETAIL string = "OBJECT_DUMP_FILE_DETAIL"
const TRIGGER_BLOCK_PROFILE string = "TRIGGER_BLOCK_PROFILE"
const TRIGGER_MUTEX_PROFILE string = "TRIGGER_MUTEX_PROFILE"

const OBJECT_LIST_HEAP_DUMP string = "OBJECT_LIST_HEAP_DUMP"
const OBJECT_DOWNLOAD_HEAP_DUMP string = "OBJECT_DOWNLOAD_HEAP_DUMP"
const OBJECT_DELETE_HEAP_DUMP string = "OBJECT_DELETE_HEAP_DUMP"

const OBJECT_CALL_CPU_PROFILE string = "OBJECT_CALL_CPU_PROFILE"
const OBJECT_CALL_BLOCK_PROFILE string = "OBJECT_CALL_BLOCK_PROFILE"
const OBJECT_CALL_MUTEX_PROFILE string = "OBJECT_CALL_MUTEX_PROFILE"


const GET_CONFIGURE_WAS string = "GET_CONFIGURE_WAS"
const SET_CONFIGURE_WAS string = "SET_CONFIGURE_WAS"
const LIST_CONFIGURE_WAS string = "LIST_CONFIGURE_WAS"
const CONFIGURE_VALUE_TYPE string = "CONFIGURE_VALUE_TYPE"
const CONFIGURE_DESC string = "CONFIGURE_DESC"
const OBJECT_ACTIVE_SERVICE_LIST string = "OBJECT_ACTIVE_SERVICE_LIST"
const OBJECT_THREAD_DETAIL string = "OBJECT_THREAD_DETAIL"

func handle(cmd string, pack netdata.Pack, in *netdata.DataInputX, out *netdata.DataOutputX) netdata.Pack {
	defer common.ReportScouterPanic()

	switch cmd {
	case "KEEP_ALIVE":
		if ac.IsTrace() {
			logger.Trace.Println("KEEP_ALIVE")
		}
		return nil
	case OBJECT_RESET_CACHE:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_RESET_CACHE")
		}
		return ResetCache(pack)
	case TRIGGER_THREAD_DUMP:
		if ac.IsTrace() {
			logger.Trace.Println("TRIGGER_THREAD_DUMP")
		}
		return dump.HeavyAllStackTrace()
	case OBJECT_DUMP_FILE_LIST:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_DUMP_FILE_LIST")
		}
		return dump.ListDumpFiles()
	case OBJECT_DUMP_FILE_DETAIL:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_DUMP_FILE_DETAIL")
		}
		dump.StreamDumpFileContents(pack, out)
		return nil
	case TRIGGER_BLOCK_PROFILE:
		if ac.IsTrace() {
			logger.Trace.Println("TRIGGER_BLOCK_PROFILE")
		}
		triggerBlockProfile(pack, out)
		return nil
	case TRIGGER_MUTEX_PROFILE:
		if ac.IsTrace() {
			logger.Trace.Println("TRIGGER_MUTEX_PROFILE")
		}
		triggerMutexProfile(pack, out)
		return nil
	case OBJECT_LIST_HEAP_DUMP:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_LIST_HEAP_DUMP")
		}
		return listBinaryDump(pack, out)
	case OBJECT_DOWNLOAD_HEAP_DUMP:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_DOWNLOAD_HEAP_DUMP")
		}
		return downloadBinaryDump(pack, out)
	case OBJECT_DELETE_HEAP_DUMP:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_DELETE_HEAP_DUMP")
		}
		return deleteBinaryDump(pack, out)
	case OBJECT_CALL_CPU_PROFILE:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_CALL_CPU_PROFILE")
		}
		return triggerBinaryCpuProfile(pack, out)
	case OBJECT_CALL_BLOCK_PROFILE:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_CALL_BLOCK_PROFILE")
		}
		return triggerBinaryBlockProfile(pack, out)
	case OBJECT_CALL_MUTEX_PROFILE:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_CALL_MUTEX_PROFILE")
		}
		return triggerBinaryMutexProfile(pack, out)
	case GET_CONFIGURE_WAS:
		if ac.IsTrace() {
			logger.Trace.Println("GET_CONFIGURE_WAS")
		}
		return loadConfig()
	case SET_CONFIGURE_WAS:
		if ac.IsTrace() {
			logger.Trace.Println("SET_CONFIGURE_WAS")
		}
		return saveConfig(pack)
	case LIST_CONFIGURE_WAS:
		if ac.IsTrace() {
			logger.Trace.Println("LIST_CONFIGURE_WAS")
		}
		//TODO
		return netdata.NewMapPack()
	case CONFIGURE_VALUE_TYPE:
		if ac.IsTrace() {
			logger.Trace.Println("CONFIGURE_VALUE_TYPE")
		}
		//TODO
		return netdata.NewMapPack()
	case CONFIGURE_DESC:
		if ac.IsTrace() {
			logger.Trace.Println("CONFIGURE_DESC")
		}
		//TODO
		return netdata.NewMapPack()
	case OBJECT_ACTIVE_SERVICE_LIST:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_ACTIVE_SERVICE_LIST")
		}
		return GetActiveList()
	case OBJECT_THREAD_DETAIL:
		if ac.IsTrace() {
			logger.Trace.Println("OBJECT_THREAD_DETAIL")
		}
		return GetGoroutineDetail(pack)
	default:
		if ac.IsTrace() {
			logger.Trace.Println("UNKNWON-HANDLER: " + cmd)
		}
	}

	return nil
}

func deleteBinaryDump(pack netdata.Pack, out *netdata.DataOutputX) netdata.Pack {
	paramPack, ok := pack.(*netdata.MapPack);
	if !ok {
		return nil
	}
	return dump.DeleteBinaryDumpFiles(paramPack.GetString("delfileName"))
}

func downloadBinaryDump(pack netdata.Pack, out *netdata.DataOutputX) netdata.Pack {
	paramPack, ok := pack.(*netdata.MapPack);
	if !ok {
		return nil
	}
	dump.DownloadBinaryDumpFiles(out, paramPack.GetString("fileName"))
	return nil
}

func listBinaryDump(pack netdata.Pack, out *netdata.DataOutputX) netdata.Pack {
	return dump.ListBinaryDumpFiles()
}

func triggerBinaryCpuProfile(pack netdata.Pack, out *netdata.DataOutputX) netdata.Pack {
	dump.ProfileBinaryCpu(30)
	p := netdata.NewMapPack()
	p.Put("success", netdata.NewBooleanValue(true));
	p.Put("msg", "Success. it take about 30 seconds.");
	return p
}

func triggerBlockProfile(pack netdata.Pack, out *netdata.DataOutputX) *netdata.MapPack {
	dump.ProfileBlock(30, 1000 * 1000, 1)
	p := netdata.NewMapPack()
	p.Put("success", netdata.NewBooleanValue(true));
	p.Put("msg", "Success. it take about 30 seconds.");
	return p
}

func triggerMutexProfile(pack netdata.Pack, out *netdata.DataOutputX) *netdata.MapPack {
	dump.ProfileMutex(30, 10, 1)
	p := netdata.NewMapPack()
	p.Put("success", netdata.NewBooleanValue(true));
	p.Put("msg", "Success. it take about 30 seconds.");
	return p
}

func triggerBinaryBlockProfile(pack netdata.Pack, out *netdata.DataOutputX) *netdata.MapPack {
	dump.ProfileBlockBinaryDump(30, 1000 * 1000)
	p := netdata.NewMapPack()
	p.Put("success", netdata.NewBooleanValue(true));
	p.Put("msg", "Success. it take about 30 seconds.");
	return p
}

func triggerBinaryMutexProfile(pack netdata.Pack, out *netdata.DataOutputX) *netdata.MapPack {
	dump.ProfileMutexBinaryDump(30, 5)
	p := netdata.NewMapPack()
	p.Put("success", netdata.NewBooleanValue(true));
	p.Put("msg", "Success. it take about 30 seconds.");
	return p
}

func GetGoroutineDetail(param netdata.Pack) *netdata.MapPack {
	p := netdata.NewMapPack()
	p.Put("Thread Name", "[No Goroutine] End");
	p.Put("State", "end");

	paramPack, ok := param.(*netdata.MapPack);
	if !ok {
		return p
	}
	tctx := tctxmanager.GetTraceContextByTxid(paramPack.GetInt64("txid"))
	if tctx == nil {
		return p
	}
	goid := tctx.Goid
	p.Put("Service Txid", util.IntToXlogString32(tctx.Txid));
	p.Put("Service Name", tctx.ServiceName);
	p.Put("Service Elapsed", util.MillisToNow(tctx.StartTime));
	p.Put("Thread Id", -1);
	p.Put("State", "n/a");
	p.Put("Thread Name", "[goroutine] " + strconv.Itoa(goid));
	p.Put("Stack Trace", "stacktrace for goroutine is not yet supported");
	p.Put("Last trace method", tctx.LastMethod);

	return p
}

func GetActiveList() *netdata.MapPack {
	mpack := netdata.NewMapPack()

	id := mpack.NewList("id")
	elapsed := mpack.NewList("elapsed")
	service := mpack.NewList("service")
	name := mpack.NewList("name")
	txid := mpack.NewList("txid")
	ip := mpack.NewList("ip")
	state := mpack.NewList("stat")
	cpu := mpack.NewList("cpu")
	sql := mpack.NewList("sql")
	subcall := mpack.NewList("subcall")

	ifTctx := tctxmanager.GetAllTctx()

	for _, tctx0 := range ifTctx {
		if tctx0 == nil {
			return mpack
		}
		tctx := tctx0.(*netio.TraceContext)
		goid := tctx.Goid
		id.AddInt64(int64(goid))
		elapsed.AddInt32(util.MillisToNow(tctx.StartTime))
		service.AddString(tctx.ServiceName)
		name.AddString("[goroutine] " + strconv.Itoa(goid))
		txid.AddString(util.IntToXlogString32(tctx.Txid))
		ip.AddString(tctx.RemoteIp)
		state.AddString("n/a")
		cpu.AddInt64(-1)
		sql.AddString("n/a")
		subcall.AddString("n/a")

	}
	mpack.Put("complete", netdata.NewBooleanValue(true))

	return mpack
}

func loadConfig() *netdata.MapPack {
	configText := conf.LoadConfigText()
	configKeys := conf.LoadConfigKeys()

	mp := netdata.NewMapPack()
	mp.Put("configKey", configKeys)
	mp.Put("agentConfig", configText)
	return mp
}

func saveConfig(pack netdata.Pack) *netdata.MapPack {
	result := netdata.NewMapPack()
	result.Put("result", false)

	if p, ok := pack.(*netdata.MapPack); ok {
		text := p.GetString("setConfig")
		success := conf.SaveConfigText(text)
		if success {
			conf.GetInstance().Refresh()
			result.Put("result", "true")
		}
	}
	return result
}

func ResetCache(pack netdata.Pack) netdata.Pack {
	netio.ResetTextSent()
	return pack
}

