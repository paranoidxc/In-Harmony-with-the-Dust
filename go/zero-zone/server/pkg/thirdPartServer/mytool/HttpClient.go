package mytool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var ErrTimeOut = errors.New("HttpClient 超时错误")
var ErrRequest = errors.New("HttpClient 请求失败")
var ErrReadBody = errors.New("HttpClient 读取响应体失败")
var ErrHttpCode = errors.New("HttpClient 响应状态码错误")
var ErrPanic = errors.New("HttpClient 捕获到[panic]异常")
var ErrConvert = errors.New("HttpClient 转换数据类型失败")
var ErrSlice = errors.New("HttpClient 非预期切片类型无法处理")
var ErrUnSupport = errors.New("HttpClient 不支持请求参数类型")

var errList []error

func init() {
	errList = []error{}
	errList = append(errList, ErrTimeOut)
	errList = append(errList, ErrRequest)
	errList = append(errList, ErrReadBody)
	errList = append(errList, ErrHttpCode)
	errList = append(errList, ErrPanic)
	errList = append(errList, ErrConvert)
	errList = append(errList, ErrSlice)
	errList = append(errList, ErrUnSupport)
}

func IsHttpClientInnerErr(err error) bool {
	for _, k := range errList {
		if k == err {
			return true
		}
	}
	return false
}

func HeaderWithJson(header map[string]string) map[string]string {
	if header == nil {
		header = map[string]string{}
	}
	header["Content-Type"] = "application/json"
	return header
}

func HeaderWithForm(header map[string]string) map[string]string {
	if header == nil {
		header = map[string]string{}
	}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	return header
}

func Header() map[string]string {
	return map[string]string{}
}

func sliceToFormBody(reqSlice []string, reqFormBody url.Values) url.Values {
	if reqFormBody == nil {
		reqFormBody = url.Values{}
	}
	for _, v := range reqSlice {
		tmp := strings.Split(v, "=")
		reqFormBody.Set(tmp[0], strings.Replace(tmp[1], "EQUAL", "=", -1))
	}

	return reqFormBody
}

func parseRequestParams(ctx context.Context, path *string, method string, header map[string]string, req interface{}) (buf *bytes.Buffer, err error) {
	defer func() {
		if xerr := recover(); xerr != nil {
			logc.Errorw(ctx, "--------- parseRequestParams 捕获到[panic]异常---------", logx.Field("err", xerr))
			//err = errors.New(fmt.Sprintf("捕获到[panic]异常:%+v", xerr))
			err = ErrPanic
		}
	}()
	reqFormBody := url.Values{}
	var reqJsonBody []byte
	if req != nil {
		kind := reflect.TypeOf(req).Kind()
		switch kind {
		case reflect.Map:
			//fmt.Printf("type:%+v\n", kind)
			reqSlice := []string{}
			reqVal := reflect.ValueOf(req)
			keys := reqVal.MapKeys()
			for _, key := range keys {
				value := reqVal.MapIndex(key)
				reqSlice = append(reqSlice, fmt.Sprintf("%v=%v", key.Interface(), value.Interface()))
			}
			if method == http.MethodGet {
				*path += "?" + url.PathEscape(strings.Replace(strings.Join(reqSlice, "&"), "EQUAL", "=", -1))
			}
			if method == http.MethodPost {
				reqFormBody = sliceToFormBody(reqSlice, reqFormBody)
			}
			//fmt.Printf("params:%s\n", strings.Join(reqSlice, "&"))
		case reflect.String:
			//fmt.Printf("type string:%s\n", req.(string))
			if method == http.MethodGet {
				*path += "?" + url.PathEscape(req.(string))
			}
		case reflect.Struct:
			//fmt.Printf("type:%+v\n", kind)
			var xerr error
			reqJsonBody, xerr = json.Marshal(req)
			if xerr != nil {
				//return nil, fmt.Errorf("json.Marshal(req) 数据失败: %v", err)
				logc.Errorw(ctx, "---------json.Marshal(req) 数据失败---------", logx.Field("req", req))
				return nil, ErrConvert
			}
			// reqJsonBodyRe := strings.Replace(string(reqJsonBody), "EQUAL", "=", -1)
			// reqJsonBody = []byte(reqJsonBodyRe)
			//fmt.Printf("reqBody:%s\n", string(reqJsonBody))
		case reflect.Slice:
			//fmt.Printf("type:%+v\n", kind)
			//fmt.Println("reqslice:", req)
			if reqSlice, ok := req.([]string); ok {
				if method == http.MethodGet {
					*path += "?" + url.PathEscape(strings.Join(reqSlice, "&"))
				} else {
					reqFormBody = sliceToFormBody(reqSlice, reqFormBody)
				}
			} else {
				logc.Errorw(ctx, "---------非预期切片类型无法处理---------", logx.Field("req", req))
				return nil, ErrSlice
			}
		default:
			logc.Errorw(ctx, "---------不支持请求参数类型---------", logx.Field("req", req))
			return nil, ErrUnSupport
		}
	}
	isJson := false
	if header != nil {
		for _, value := range header {
			if strings.Contains(value, "json") {
				isJson = true
			}
			//fmt.Println("key", key, "val", value)
		}
	}

	if isJson {
		buf = bytes.NewBuffer(reqJsonBody)
	} else {
		buf = bytes.NewBufferString(reqFormBody.Encode())
	}

	logc.Infow(ctx, "---------请求第三方信息---------",
		logx.Field("path", *path),
		logx.Field("method", method),
		logx.Field("reqJson", string(reqJsonBody)),
		logx.Field("reqFormBody", reqFormBody),
		logx.Field("isJson", isJson),
		logx.Field("header", header),
	)

	return
}

func DoHttpRequest(ctx context.Context, path string, method string, header map[string]string, req interface{}) (body []byte, err error) {
	defer func() {
		if xerr := recover(); xerr != nil {
			logc.Errorw(ctx, "--------- DoHttpRequest 捕获到[panic]异常---------", logx.Field("err", xerr))
			err = ErrPanic
		}
	}()

	buffer, err := parseRequestParams(ctx, &path, method, header, req)
	if err != nil {
		return
	}

	var request *http.Request
	request, err = http.NewRequest(method, path, buffer)

	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}

	// 发起 HTTP 请求
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			logc.Errorw(ctx, "---------发起 HTTP 超时---------", logx.Field("err", err))
			return nil, ErrTimeOut
		}
		logc.Errorw(ctx, "---------发起 HTTP 请求失败---------", logx.Field("err", err))
		return nil, ErrRequest
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logc.Errorw(ctx, "---------读取响应体失败---------", logx.Field("err", err))
		return nil, ErrReadBody
	}

	logc.Infow(ctx, "---------请求第三方响应---------", logx.Field("body", string(body)))
	// 检查响应状态码 只要不是正常的200 都是错误
	if resp.StatusCode != http.StatusOK {
		logc.Infow(ctx, "---------请求第三方响应 HTTP 状态码错误---------", logx.Field("CODE", resp.StatusCode))
		return nil, ErrHttpCode
	}

	// 解析 body 逻辑层自己处理
	return body, nil
}
