package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetAppPath returns current parogram path
func GetAppPath() (string, error) {
	c := os.Getenv("SCOUTER_AGENT_HOME")
	if c == "" {
		var err error
		c, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	return c, nil
}

func GetScouterPath() string {
	scouterPath := os.Getenv("SCOUTER_AGENT_HOME")
	if scouterPath == "" {
		appPath, err := GetAppPath()
		if err != nil {
			appPath = os.TempDir()
		}
		scouterPath = filepath.Join(appPath, "scouter")
	}
	MakeDir(scouterPath)
	return scouterPath
}

// MakeDir makes dir given path.
func MakeDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("cannot create folder %s \n", path)
	}

}
