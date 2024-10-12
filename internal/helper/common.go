package helper

import "math"

type ErrorStruct struct {
	Err     error
	Code    int
	Message string
}

// Pembulatan dinamis
func RoundNumber(x float64, precision int, method ...string) float64 {
	if method == nil {
		method = append(method, "")
	}

	multiplier := math.Pow(10, float64(precision))
	switch method[0] {
	case "floor":
		return math.Floor(x*multiplier) / multiplier
	case "round":
		return math.Round(x*multiplier) / multiplier
	case "ceil":
		return math.Ceil(x*multiplier) / multiplier
	default:
		// Jika metode tidak diketahui, gunakan round sebagai default
		return math.Round(x*multiplier) / multiplier
	}
}
