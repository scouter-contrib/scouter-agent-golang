package util

import (
	"fmt"
	"time"
)

var secondsPerMin = int64(60)
var secondsPerHour = int64(60 * secondsPerMin)
var secondsPerDay = int64(24 * secondsPerHour)

//GetDuration returns long to string format
func GetDuration(seconds int64) string {
	day := seconds / secondsPerDay
	hour := (seconds - (day * secondsPerDay)) / secondsPerHour
	min := (seconds - (day * secondsPerDay) - (hour * secondsPerHour)) / secondsPerMin
	secs := (seconds - (day * secondsPerDay) - (hour * secondsPerHour) - (min * secondsPerMin))
	return fmt.Sprintf("%dD %dH %dM %dS", day, hour, min, secs)
}

func MillisBetween(from time.Time, to time.Time) int32 {
	return int32(to.Sub(from).Milliseconds())
}

func MillisToNow(from time.Time) int32 {
	return int32(time.Now().Sub(from).Milliseconds())
}

func TimeToMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
