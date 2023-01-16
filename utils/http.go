package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"io"

	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/thinkeridea/go-extend/exbytes"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
)

var (
	Httpclient *http.Client
)

func init() {
	Httpclient = &http.Client{Transport: &http.Transport{
		TLSHandshakeTimeout: 3 * time.Second,
	},
		Timeout: 6 * time.Second,
	}
}

func HttpRequest(method string, url string, data interface{}) ([]byte, int) {
	var req *http.Request
	if method == "GET" {
		dataMap, err := Struct2Map(data)
		if err != nil {
			return nil, 500
		}
		requestParamStr := getRequestParamStr(dataMap)
		req, _ = http.NewRequest(method, fmt.Sprintf("%s?%s", url, requestParamStr), nil)
	} else {
		md, _ := json.Marshal(data)
		req, _ = http.NewRequest(method, url, bytes.NewReader(md))
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := Httpclient.Do(req)
	if err != nil {
		logrus.Errorf("call appInternal api err: %s, %s", url, err.Error())
		return nil, 500
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Read response body err: %s", err)
		return nil, 500
	}

	return content, resp.StatusCode
}

func getRequestParamStr(data map[string]interface{}) string {
	resultString := make([]string, 0)
	keys := make([]string, 0)
	for key, _ := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value, _ := data[key]
		switch x := value.(type) {
		case string, float32, float64, int, int8, int16, int32, int64:
			resultString = append(resultString, fmt.Sprintf("%s=%v", key, x))
		case []interface{}:
			for _, v := range x {
				resultString = append(resultString, fmt.Sprintf("%s[]=%v", key, v))
			}
		}
	}
	return strings.Join(resultString, "&")
}

type HttpPool struct {
	pool sync.Pool
}

func NewHttpPool() *HttpPool {
	return &HttpPool{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 4096))
			},
		},
	}
}
func (api *HttpPool) HttpRequest(method string, reqUrl string, body interface{}) ([]byte, error) {
	var err error
	buffer := api.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			api.pool.Put(buffer)
			buffer = nil
		}
	}()

	e := jsoniter.NewEncoder(buffer)
	err = e.Encode(body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"request": body}).Errorf("jsoniter.Marshal fail: %v", err)
		return nil, err
	}
	data := buffer.Bytes()
	req, err := http.NewRequest(method, reqUrl, buffer)
	req.Header.Add("Content-Type", "application/json")
	//logrus.WithFields(logrus.Fields{"data": exbytes.ToString(data),}).Errorf("http.NewRequest failed: %v", err)

	if err != nil {
		logrus.WithFields(logrus.Fields{"data": exbytes.ToString(data)}).Errorf("http.NewRequest failed: %v", err)
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": reqUrl,
		}).Errorf("query service failed %v", err)
		return nil, fmt.Errorf("query service failed %v", err)
	}
	buffer.Reset()
	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, err
	}

	respData := buffer.Bytes()
	logrus.WithFields(logrus.Fields{"response_json": exbytes.ToString(respData)}).Debug("response json")
	logrus.WithFields(logrus.Fields{"url": reqUrl}).Debug("url")

	if resp.StatusCode != 200 {
		var errorOut struct {
			Errors []string `json:"errors"`
		}
		err = jsoniter.Unmarshal(respData, &errorOut)
		logrus.WithFields(logrus.Fields{
			"url":         reqUrl,
			"status":      resp.Status,
			"status_code": resp.StatusCode,
			"errors":      errorOut.Errors,
		}).Errorf("invalid http status code")
		errorStr := ""
		for k, v := range errorOut.Errors {
			if 0 == k {
				errorStr = v
				continue
			}
			errorStr += "," + v
		}
		return respData, formutil.NewWithStatusError(resp.StatusCode, errorStr)
		//return respData, err
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"data": exbytes.ToString(respData),
			"url":  reqUrl,
		}).Errorf("jsoniter.Unmarshal failed, error:%v", err)

		return respData, err
	}
	return respData, nil
}
