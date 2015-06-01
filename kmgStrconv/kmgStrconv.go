package kmgStrconv

import "strconv"

func AtoIDefault0(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ParseFloat64(f string) (float64, error) {
	return strconv.ParseFloat(f, 64)
}
