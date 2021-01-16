package dump

import (
	"bufio"
	"fmt"
	"github.com/scouter-contrib/scouter-agent-golang/scouterx/common/constants/tcpflag"
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
	path := util.GetScouterPath()
	dumpPath := filepath.Join(path, "dump")
	util.MakeDir(dumpPath)

	fileName := filepath.Join(dumpPath, "go_dump_" + time.Now().Format("20060102_150405") + ".log")

	var w io.Writer = os.Stdout
	if fileName != "" {
		f, err := os.Create(fileName)
		if err == nil {
			w = bufio.NewWriter(f)
		} else {
			return "", err
		}
	}
	profile := pprof.Lookup("goroutine")
	return fileName, profile.WriteTo(w, 2)
}

func getDumpPath() string {
	path := util.GetScouterPath()
	dumpPath := filepath.Join(path, "dump")
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
