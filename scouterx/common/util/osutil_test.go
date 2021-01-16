package util

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestOSUtil(t *testing.T) {
	p, _ := GetAppPath()
	fmt.Printf("%s", p)
}

func TestScouterDumpPath(t *testing.T) {
	p := GetScouterPath()
	fmt.Printf("%s\n", p)
	fmt.Printf("%s\n", filepath.Join(p, "dump"))
}
