package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"math/rand"
	"net/http"
	"reflect"

	"encoding/json"
	"fmt"
	"time"
	"zero-zone/pkg/globalkey"

	"github.com/zeromicro/go-zero/core/logx"
)

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

func Sha256(str string) string {
	sum256 := sha256.Sum256([]byte(str))
	shaStr := fmt.Sprintf("%x", sum256) // 将[]byte转成16进制
	return shaStr
}

type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func GetUserId(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(globalkey.SysJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}

	return uid
}

func ArrayUniqueValue[T any](arr []T) []T {
	size := len(arr)
	result := make([]T, 0, size)
	temp := map[any]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}

	return result
}

func ArrayContainValue(arr []int64, search int64) bool {
	for _, v := range arr {
		if v == search {
			return true
		}
	}

	return false
}

func Intersect(slice1 []int64, slice2 []int64) []int64 {
	m := make(map[int64]int64)
	n := make([]int64, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}

	return n
}

func Difference(slice1 []int64, slice2 []int64) []int64 {
	m := make(map[int64]int)
	n := make([]int64, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, v := range slice1 {
		times, _ := m[v]
		if times == 0 {
			n = append(n, v)
		}
	}

	return n
}

func Time2Str(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func DateTimeStrCompact() string {
	return time.Now().Format("20060102150405")
}

func DateTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Unix2TimeStr(unix int64) string {
	if unix > 0 {
		tm := time.Unix(unix, 0)
		return tm.Format("2006-01-02 15:04:05")
	}
	return ""
}

func RandomNum(n int) string {
	var letters = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[RandomInt(0, len(letters))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[RandomInt(0, len(letters))]
	}
	return string(b)
}

func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func Host(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	return scheme + r.Host
}
