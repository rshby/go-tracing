package helper

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

func Dump(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return string(marshal)
}

func TimeToStringIndonesia(v time.Time) string {
	v = v.Add(7 * time.Hour)
	return v.Format("2006-01-02 15:04:05")
}

func ExpectNumber[T int | int64 | uint | float32 | float64](v any) T {
	var result T

	var inputString string
	switch v.(type) {
	case string:
		inputString = v.(string)
	default:
		inputString = Dump(v)
	}

	switch reflect.TypeOf(result).Kind() {
	case reflect.Uint, reflect.Int, reflect.Int64:
		i, _ := strconv.ParseInt(inputString, 10, 64)
		return T(i)
	case reflect.Float64, reflect.Float32:
		float, _ := strconv.ParseFloat(inputString, 64)
		return T(float)
	default:
		return result
	}
}

func MyCaller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		return details.Name()
	}

	return "not know"
}
