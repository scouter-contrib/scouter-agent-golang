package util
import "math"

//RoundUp returns round float
func RoundUp(input float64, places int) float64 {
	var round float64
	var newVal float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return newVal
}