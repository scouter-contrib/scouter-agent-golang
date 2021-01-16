package util

// CopyArray returns copied array
func CopyArray(data []byte, pos int, length int) []byte {
	return data[pos : pos+length]
}
