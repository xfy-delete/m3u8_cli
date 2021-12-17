package decode

import (
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/xfy520/m3u8_cli/package/tool"
)

func NfmoviesDecryptM3u8(byteArray []byte) string {
	return string(byteArray)
}

func DdyunDecryptM3u8(byteArray []byte) string {
	return ""
}

func GetVaildM3u8Url(url string) (string, error) {
	reg := regexp.MustCompile(`\w{20,}`)
	s := reg.FindAllString(url, -1)
	if len(s) == 0 {
		return "", errors.New("system error")
	}
	id := s[0]
	tm := tool.GetTimeStamp(false)
	t := strconv.FormatInt((tm/0x186a0)*0x64, 10)
	tmp := id + "duoduo" + "1" + t
	bs := tool.StrToBytes(tmp)
	has := md5.Sum(bs)
	key := fmt.Sprintf("%x", has)
	re, err := regexp.Compile(`1/\w{20,}`)
	return re.ReplaceAllString(url, key), err
}
