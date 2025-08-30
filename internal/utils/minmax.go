package utils

// Internal function that controls a value with max and min
func MinMax(min, v, max int) int {
	switch {
	case v < min:
		return min
	case max == 0:
		return v
	case v > max:
		return max
	default:
		return v
	}
}
