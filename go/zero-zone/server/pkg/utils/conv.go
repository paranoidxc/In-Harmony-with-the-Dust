package utils

import "strconv"

func Int642Str(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Float642Str(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

func Str2Float64(str string) (float64, error) {
	num, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return num, err
	}

	return num, nil
}

func Str2Int64(str string) (int64, error) {
	num, err := strconv.ParseInt(str, 10, strconv.IntSize)
	if err != nil {
		return num, err
	}

	return num, nil
}

func Str2Uint(str string) (uint, error) {
	var uintNum uint
	num, err := strconv.ParseUint(str, 10, strconv.IntSize)
	if err != nil {
		return uintNum, err
	}
	uintNum = uint(num)

	return uintNum, nil
}
