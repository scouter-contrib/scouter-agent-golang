package util

import (
	"strconv"
	"strings"
)

var emptyIp = []byte{0, 0, 0, 0}

func IpToBytes(ip string) []byte {
	result := make([]byte, 4)
	split := strings.Split(ip, ".")
	if len(split) != 4 {
		return emptyIp
	}
	for i, part := range split {
		v, err := strconv.ParseUint(part, 10, 8)
		if (err != nil || v < 0) {
			return emptyIp
		}
		result[i] = byte(v)
	}
	return result
}
