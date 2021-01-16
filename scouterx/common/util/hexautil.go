package util

import "strconv"

const xplus = "x"
const xminus = "z"

func IntToXlogString32(num int64) string {
	minus := num < 0
	if minus {
		return xminus + strconv.FormatInt(-num, 32)
	} else {
		return xplus + strconv.FormatInt(num, 32)
	}
}
