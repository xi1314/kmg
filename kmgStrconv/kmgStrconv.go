package kmgStrconv

import "strconv"

func AtoIDefault0(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func FormatInt(i int) string {
	return strconv.Itoa(i)
}
func MustParseInt(f string) int {
	i, err := strconv.Atoi(f)
	if err != nil {
		panic(err)
	}
	return i
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func FormatFloatPrec0(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64)
}

//以两位精度把浮点转成字符串
func FormatFloatPrec2(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func FormatFloatPrec4(f float64) string {
	return strconv.FormatFloat(f, 'f', 4, 64)
}

func ParseFloat64(f string) (float64, error) {
	return strconv.ParseFloat(f, 64)
}

func MustParseFloat64(f string) float64 {
	out, err := strconv.ParseFloat(f, 64)
	if err != nil {
		panic(err)
	}
	return out
}

func ParseFloat64Default0(s string) float64 {
	out, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return out
}

func MustParseBool(f string) bool {
	out, err := strconv.ParseBool(f)
	if err != nil {
		panic(err)
	}
	return out
}
