package dump

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/tcpflag"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/logger"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/netdata"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/util"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

const suffix = "pprof"

func StackTrace(count int) string {
	// fine, line
	pcLen := count
	if pcLen == 0 {
		pcLen = 1
	}
	var sb strings.Builder
	pc := make([]uintptr, pcLen)
	n := runtime.Callers(1, pc)
	if n < 1 {
		return ""
	}
	for _, e := range pc {
		f := runtime.FuncForPC(e)
		if f != nil && f.Name() != "" {
			filePath, line := f.FileLine(e)
			_, file := path.Split(filePath)
			var str string
			if sb.Len() == 0 {
				str = fmt.Sprintf("%v %v:%v", f.Name(), file, line)
			} else {
				str = fmt.Sprintf("\n%v %v:%v", f.Name(), file, line)
			}
			sb.WriteString(str)
		}
	}

	return sb.String()
}


func HeavyAllStackTrace() netdata.Pack {
	fileName, err := heavyAllStackTrace()
	pack := netdata.NewMapPack()
	if err != nil {
		return pack
	}
	pack.Put("name", fileName)
	return pack
}

func heavyAllStackTrace() (string, error) {
	fileName := filepath.Join(getDumpPath(), "go_dump_" + time.Now().Format("20060102_150405") + ".log")

	var w io.Writer = os.Stdout
	if fileName != "" {
		f, err := os.Create(fileName)
		if err == nil {
			w = bufio.NewWriter(f)
			defer f.Close()
		} else {
			return "", err
		}
	}
	profile := pprof.Lookup("goroutine")
	if profile == nil {
		logger.Error.Printf("could not find goroutine profile\n")
		return "", errors.New("could not find goroutine profile")
	}
	return fileName, profile.WriteTo(w, 2)
}

func ProfileBinaryCpu(sec int) netdata.Pack {
	fileName := filepath.Join(getBinaryDumpPath(), "go_cpu_" + time.Now().Format("20060102_150405") + "." + suffix)

	go profileCpu(fileName, sec)
	pack := netdata.NewMapPack()
	return pack
}

func profileCpu(fileName string, sec int) {
	f, err := os.Create(fileName)
	if err == nil {
		defer f.Close()
	} else {
		logger.Error.Printf("could not make file to start CPU profile: %s\n", err)
		return
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		logger.Error.Printf("could not start CPU profile: %s\n", err)
		return
	}
	defer pprof.StopCPUProfile()
	time.Sleep(time.Duration(sec) * time.Second)
}

func ProfileBlock(sec, rate, level int) netdata.Pack {
	fileName := filepath.Join(getDumpPath(), "go_block_" + time.Now().Format("20060102_150405") + ".log")

	go profileBlock(fileName, sec, rate, level)
	pack := netdata.NewMapPack()
	return pack
}

func ProfileBlockBinaryDump(sec, rate int) netdata.Pack {
	fileName := filepath.Join(getBinaryDumpPath(), "go_block_" + time.Now().Format("20060102_150405") + "." + suffix)

	go profileBlock(fileName, sec, rate, 0)
	pack := netdata.NewMapPack()
	return pack
}

func profileBlock(fileName string, sec, rate, level int) {
	runtime.SetBlockProfileRate(rate)
	defer runtime.SetBlockProfileRate(-1)
	time.Sleep(time.Duration(sec) * time.Second)

	var w *bufio.Writer
	if fileName != "" {
		f, err := os.Create(fileName)
		if err == nil {
			w = bufio.NewWriter(f)
			defer f.Close()
		} else {
			return
		}
	}
	profile := pprof.Lookup("block")
	if profile == nil {
		logger.Error.Printf("could not find block profile\n")
	}
	if err := profile.WriteTo(w, level); err != nil {
		logger.Error.Printf("could not run block profile: %s\n", err)
	}
	w.Flush()
}


func ProfileMutex(sec, rate, level int) {
	fileName := filepath.Join(getDumpPath(), "go_mutex_" + time.Now().Format("20060102_150405") + ".log")

	go profileMutex(fileName, sec, rate, level)
	return
}

func ProfileMutexBinaryDump(sec, rate int) {
	fileName := filepath.Join(getBinaryDumpPath(), "go_mutex_" + time.Now().Format("20060102_150405") + "." + suffix)

	go profileMutex(fileName, sec, rate, 0)
	return
}


func profileMutex(fileName string, sec, rate, level int) {
	old := runtime.SetMutexProfileFraction(rate)
	defer runtime.SetMutexProfileFraction(old)
	time.Sleep(time.Duration(sec) * time.Second)

	var w *bufio.Writer
	if fileName != "" {
		f, err := os.Create(fileName)
		if err == nil {
			w = bufio.NewWriter(f)
			defer f.Close()
		} else {
			return
		}
	}
	profile := pprof.Lookup("mutex")
	if profile == nil {
		logger.Error.Printf("could not find mutex profile\n")
	}
	if err := profile.WriteTo(w, level); err != nil {
		logger.Error.Printf("could not run mutex profile: %s\n", err)
	}
	w.Flush()
}

func getDumpPath() string {
	path := util.GetScouterPath()
	dumpPath := filepath.Join(path, "dump")
	util.MakeDir(dumpPath)
	return dumpPath
}

func getBinaryDumpPath() string {
	path := util.GetScouterPath()
	dumpPath := filepath.Join(path, "binary_dump")
	util.MakeDir(dumpPath)
	return dumpPath
}

func ListDumpFiles() netdata.Pack {
	pack := netdata.NewMapPack()
	nameLv := netdata.NewListValue()
	pack.Put("name", nameLv)
	sizeLv := netdata.NewListValue()
	pack.Put("size", sizeLv)
	modifiedLv := netdata.NewListValue()
	pack.Put("last_modified", modifiedLv)

	dumpPath := getDumpPath()

	filepath.Walk(dumpPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			nameLv.AddString(info.Name())
			sizeLv.AddInt64(info.Size())
			modifiedLv.AddInt64(util.TimeToMillis(info.ModTime()))
		}
		return nil
	})

	return pack
}

func ListBinaryDumpFiles() netdata.Pack {
	pack := netdata.NewMapPack()
	nameLv := netdata.NewListValue()
	pack.Put("name", nameLv)
	sizeLv := netdata.NewListValue()
	pack.Put("size", sizeLv)

	dumpPath := getBinaryDumpPath()

	filepath.Walk(dumpPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			nameLv.AddString(info.Name())
			sizeLv.AddInt64(info.Size())
		}
		return nil
	})

	return pack
}

func DownloadBinaryDumpFiles(out *netdata.DataOutputX, fileName string) {
	if fileName == "" || !strings.HasSuffix(fileName, suffix){
		return
	}
	dumpPath := getBinaryDumpPath()
	fullName := filepath.Join(dumpPath, fileName)

	f, err := os.Open(fullName)
	if err != nil {
		logger.Error.Printf("could not open file: %s\n", fullName)
		return
	}
	defer f.Close()

	//var b = make([]byte, 2 * 1024 * 1024, 2 * 1024 * 1024)
	var b = make([]byte, 300, 300)
	var offset int64 = 0

	for {
		len, err := f.ReadAt(b, offset)
		if err != nil {
			if err != io.EOF {
				logger.Error.Printf("file read error: %s\n", err.Error())
			} else {
				out.WriteUInt8(tcpflag.HasNEXT);
				out.WriteBlob(b[0:len]);
			}
			break
		}
		out.WriteUInt8(tcpflag.HasNEXT);
		out.WriteBlob(b[0:len]);
		offset = offset + int64(len)
	}
	return
}

func DeleteBinaryDumpFiles(fileName string) netdata.Pack {
	pack := netdata.NewMapPack()
	if fileName == "" || !strings.HasSuffix(fileName, suffix){
		pack.Put("success", netdata.NewBooleanValue(false))
		pack.Put("msg", "no fileName")
		return pack
	}
	dumpPath := getBinaryDumpPath()
	fullName := filepath.Join(dumpPath, fileName)

	err := os.Remove(fullName)
	if err != nil {
		pack.Put("success", netdata.NewBooleanValue(false))
		pack.Put("msg", err.Error())
	} else {
		pack.Put("success", netdata.NewBooleanValue(true))
		pack.Put("msg", "success")
	}
	return pack
}

func StreamDumpFileContents(param netdata.Pack, out *netdata.DataOutputX) {
	p, ok := param.(*netdata.MapPack)
	if !ok {
		return
	}
	name := p.GetString("name")
	dumpFile := filepath.Join(getDumpPath(), name)

	data, err := os.Open(dumpFile)
	if err != nil {
		return
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	part := make([]byte, 4 * 1024)
	var count int

	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		_, err := out.WriteUInt8(tcpflag.HasNEXT)
		err = out.WriteBlob(part[:count])
		if err != nil {
			return
		}
	}
	return

}
