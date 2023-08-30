package util

import "strconv"

func FloatString(float float64, precision int) string {
	return strconv.FormatFloat(float, 'f', precision, 64)
}
