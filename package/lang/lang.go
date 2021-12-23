package lang

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/xfy520/m3u8_cli/package/log"
)

type Contact struct {
	ParseExit                     string `json:"ParseExit"`
	StartReParsing                string `json:"StartReParsing"`
	SelectPlaylist                string `json:"SelectPlaylist"`
	WrtingMasterMeta              string `json:"WrtingMasterMeta"`
	MasterListFound               string `json:"MasterListFound"`
	WrtingMeta                    string `json:"WrtingMeta"`
	DownloadingExternalAudioTrack string `json:"DownloadingExternalAudioTrack"`
	InvalidM3u8Error              string `json:"InvalidM3u8Error"`
	NotSupportMethodError         string `json:"NotSupportMethodError"`
	DownloadingM3u8Key            string `json:"DownloadingM3u8Key"`
	ParseError                    string `json:"ParseError"`
	RedirectCountError            string `json:"RedirectCountError"`
	FilePathError                 string `json:"FilePathError"`
	StartParsing                  string `json:"StartParsing"`
	AragError                     string `json:"AragError"`
	SavePath                      string `json:"SavePath"`
	FileName                      string `json:"FileName"`
	DirLongError                  string `json:"DirLongError"`
	FileNameError                 string `json:"FileNameError"`
	InputRetryCount               string `json:"InputRetryCount"`
	UrlError                      string `json:"UrlError"`
	CommandError                  string `json:"CommandError"`
	Usage                         string `json:"Usage"`
	StartUp                       string `json:"StartUp"`
	ProjectPathError              string `json:"ProjectPathError"`
	AnyKey                        string `json:"AnyKey"`
	FfmpegPathUsage               string `json:"FfmpegPathUsage"`
	FfmpegLost                    string `json:"FfmpegLost"`
	FfmpegTip                     string `json:"FfmpegTip"`
	Url                           string `json:"Url"`
	WorkDir                       string `json:"WorkDir"`
	SaveName                      string `json:"SaveName"`
	BaseUrl                       string `json:"BaseUrl"`
	Headers                       string `json:"Headers"`
	MaxThreads                    string `json:"MaxThreads"`
	MinThreads                    string `json:"MinThreads"`
	RetryCount                    string `json:"RetryCount"`
	TimeOut                       string `json:"TimeOut"`
	MuxSetJson                    string `json:"MuxSetJson"`
	UseKeyFile                    string `json:"UseKeyFile"`
	UseKeyBase64                  string `json:"UseKeyBase64"`
	UseKeyIV                      string `json:"UseKeyIV"`
	DownloadRange                 string `json:"DownloadRange"`
	LiveRecDur                    string `json:"LiveRecDur"`
	StopSpeed                     string `json:"StopSpeed"`
	MaxSpeed                      string `json:"MaxSpeed"`
	ProxyAddress                  string `json:"ProxyAddress"`
	EnableDelAfterDone            string `json:"EnableDelAfterDone"`
	EnableMuxFastStart            string `json:"EnableMuxFastStart"`
	EnableBinaryMerge             string `json:"EnableBinaryMerge"`
	EnableParseOnly               string `json:"EnableParseOnly"`
	EnableAudioOnly               string `json:"EnableAudioOnly"`
	DisableDateInfo               string `json:"DisableDateInfo"`
	NoMerge                       string `json:"NoMerge"`
	NoProxy                       string `json:"NoProxy"`
	DisableIntegrityCheck         string `json:"DisableIntegrityCheck"`
}

var Lang Contact

func GetFile(dataFile string) ([]byte, error) {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		datapath := path.Join(path.Dir(filename), dataFile)
		f, err := os.Open(datapath)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	return nil, errors.New("system error")
}

func GetLocale() (string, error) {
	envlang, ok := os.LookupEnv("LANG")
	if ok {
		return strings.Split(envlang, ".")[0], nil
	}
	cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
	output, err := cmd.Output()
	if err == nil {
		return strings.Trim(string(output), "\r\n"), nil
	}
	return "", errors.New("cannot determine locale")
}

func init() {
	var lang = "en_US"
	locale, err := GetLocale()
	if err != nil {
		log.Error(err.Error())
	}
	if strings.Contains(locale, "TW") || strings.Contains(locale, "HK") || strings.Contains(locale, "MO") {
		lang = "zh_TW"
	} else if strings.Contains(locale, "CN") || strings.Contains(locale, "SG") {
		lang = "zh_CN"
	}
	content, err := GetFile(path.Join("data", lang+".json"))
	if err != nil {
		log.Error("open file error: " + err.Error())
		return
	}
	err = json.Unmarshal([]byte(content), &Lang)
	if err != nil {
		log.Error("error: " + err.Error())
		return
	}
}
