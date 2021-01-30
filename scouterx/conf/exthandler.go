package conf

import (
	"io/ioutil"
)

func LoadConfigText() string {
	filePath := getConfFilePath()
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		//TODO logging
		return ""
	}
	return string(bytes)
}

func SaveConfigText(text string) bool {
	err := ioutil.WriteFile(getConfFilePath(), []byte(text), 0644)
	if err != nil {
		//TODO logging
		return false
	}
	return true
}

