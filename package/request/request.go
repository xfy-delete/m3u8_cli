package request

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"

	"github.com/xfy520/m3u8_cli/package/agent"
	"github.com/xfy520/m3u8_cli/package/tool"
)

var (
	UseProxyAddress string = ""
	NoProxy         bool   = false
)

type Request interface {
	Send(redirectCount int) ([]byte, error)
	Set(key string, value string)
	InitHeader()
	SetHeaders(headers string)
	Get302() (string, error)
}

type request struct {
	client *http.Client
	req    *http.Request
}

func Strval(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value := value.(type) {
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 64)
	case int:
		return strconv.Itoa(value)
	case uint:
		return strconv.Itoa(int(value))
	case int8:
		return strconv.Itoa(int(value))
	case uint8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case uint16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case uint32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case string:
		return value
	case []byte:
		return string(value)
	default:
		newValue, _ := json.Marshal(value)
		return string(newValue)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func getHeaderStr(headers string) string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return headers
	}
	if tool.Exists(headers) && tool.IsFile(headers) {
		headersByte, err := tool.ReadFile(headers)
		tool.Check(err)
		return string(headersByte)
	}
	headers = path.Join(filename, headers)
	if tool.Exists(headers) && tool.IsFile(headers) {
		headersByte, err := tool.ReadFile(headers)
		tool.Check(err)
		return string(headersByte)
	}
	return headers
}

func getHeaderMap(headers string) map[string]interface{} {
	headersBytes := tool.StrToBytes(getHeaderStr(headers))
	if json.Valid(headersBytes) {
		jsonMap := make(map[string]interface{})
		err := json.Unmarshal(headersBytes, &jsonMap)
		if err != nil {
			return make(map[string]interface{})
		}
		return jsonMap
	} else {
		headersArray := strings.Split(headers, "|")
		jsonMap := make(map[string]interface{})
		for _, value := range headersArray {
			values := strings.SplitN(value, ":", 1)
			if len(values) == 2 {
				jsonMap[values[0]] = values[1]
			}
		}
		return jsonMap
	}
}

func (r *request) Set(key string, value string) {
	r.req.Header.Set(key, value)
}

func (r *request) InitHeader() {
	s := rand.NewSource(time.Now().Unix())
	r.req.Header = http.Header{}
	userAgent := agent.UserAgent[rand.New(s).Intn(len(agent.UserAgent))]
	r.req.Header.Set("accept-encoding", "gzip, deflate")
	r.req.Header.Set("accept", "*/*")
	r.req.Header.Set("user-agent", userAgent)
}

func (r *request) SetHeaders(headers string) {
	if headers != "" {
		jsonMap := getHeaderMap(headers)
		for k, v := range jsonMap {
			r.req.Header.Set(k, Strval(v))
		}
	}
}

func New(uri string, method string, timeOut time.Duration, banRedirect bool) (Request, error) {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	var proxy *url.URL = nil
	if !NoProxy && UseProxyAddress != "" {
		proxy, _ = url.Parse(UseProxyAddress)
	}
	return &request{
		client: &http.Client{
			Timeout: time.Second * timeOut,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if banRedirect {
					return http.ErrUseLastResponse
				}
				if len(via) >= 30 {
					return errors.New("redirect too times")
				}
				return nil
			},
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*timeOut)
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * timeOut))
					return conn, nil
				},
			},
		},
		req: req,
	}, nil
}

func (r *request) Send(redirectCount int) ([]byte, error) {
	if redirectCount == 0 && redirectCount != -1 {
		return nil, errors.New("")
	}
	redirectCount -= 1
	res, err := r.client.Do(r.req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 302 {
		loc, err := res.Location()
		if err != nil {
			return nil, err
		}
		r.req.URL = loc
		return r.Send(redirectCount)
	}
	body := res.Body
	if res.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	}
	if res.Header.Get("Content-Encoding") == "br" {
		b, err := ioutil.ReadAll(brotli.NewReader(res.Body))
		return b, err
	}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *request) Get302() (string, error) {
	res, err := r.client.Do(r.req)
	if err != nil {
		return "", err
	}
	return res.Request.URL.String(), nil
}
