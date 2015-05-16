package kmgStrconv

import "strconv"

func AtoIDefault0(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
