package download

import (
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/xfy520/m3u8_cli/package/log"
	"github.com/xfy520/m3u8_cli/package/request"
	"github.com/xfy520/m3u8_cli/package/tool"
)

// 下载文件字节流
func HttpDownloadFileToBytes(uri string, headers string, timeOut time.Duration) ([]byte, error) {
	if strings.HasPrefix(uri, "file:") {
		u, err := url.Parse(uri)
		if err != nil {
			log.Error(err.Error())
		}
		uri = u.Path
		sysType := runtime.GOOS
		uri = tool.IfString(sysType == "windows", uri[:1], uri)
		infbytes, err := tool.ReadFile(uri)
		if err != nil {
			return nil, err
		}
		return infbytes, nil
	}
	req, err := request.New(uri, http.MethodGet, timeOut, headers)
	if err != nil {
		return nil, err
	}
	return req.Send()
}
