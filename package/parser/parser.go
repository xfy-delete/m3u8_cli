package parser

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/xfy520/m3u8_cli/package/decode"
	"github.com/xfy520/m3u8_cli/package/download"
	"github.com/xfy520/m3u8_cli/package/ffmpeg"
	"github.com/xfy520/m3u8_cli/package/lang"
	"github.com/xfy520/m3u8_cli/package/tool"
)

type Audio struct {
	Name     string
	Language string
	Uri      string
	Channels string
	ToString func() string
}

type Subtitle struct {
	Name     string
	Language string
	Uri      string
	ToString func() string
}

type Parser interface {
	SetBaseUrl(baseUrl string)
	Parse() error
}

type parser struct {
	downName          string
	downDir           string
	m3u8Url           string
	keyBase64         string
	keyIV             string
	keyFile           string
	headers           string
	baseUrl           string
	m3u8SavePath      string
	jsonSavePath      string
	extLists          []string
	MEDIA_AUDIO_GROUP map[string][]Audio
	MEDIA_SUB_GROUP   map[string][]Subtitle
}

func AudioNew(Name string, Language string, Uri string, Channels string) *Audio {
	return &Audio{Name: Language, Uri: Uri, Channels: Channels, ToString: func() string {
		return strings.ReplaceAll("["+Name+"] ["+Language+"] ["+tool.IfString(Channels == "", "", Channels+"ch")+"]", "[]", "")
	}}
}

func SubtitleNew(Name string, Language string, Uri string) *Subtitle {
	return &Subtitle{Name: Language, Uri: Uri, ToString: func() string {
		return "[" + Name + "] [" + Language + "]"
	}}
}

const (
	baseUrl      = ""
	m3u8SavePath = ""
	jsonSavePath = ""
)

var (
	extLists          = []string{}
	MEDIA_AUDIO_GROUP = make(map[string][]Audio)
	MEDIA_SUB_GROUP   = make(map[string][]Subtitle)
)

func New(downName string, downDir string, m3u8Url string, keyBase64 string, keyIV string, keyFile string, headers string) Parser {
	return &parser{downName, downDir, m3u8Url, keyBase64, keyIV, keyFile,
		headers, baseUrl, m3u8SavePath, jsonSavePath, extLists, MEDIA_AUDIO_GROUP,
		MEDIA_SUB_GROUP,
	}
}

func (p *parser) SetBaseUrl(baseUrl string) {
	p.baseUrl = baseUrl
}

func (p *parser) Parse() error {
	ffmpeg.REC_TIME = ""
	p.m3u8SavePath = path.Join(p.downDir, "raw.m3u8")
	if !tool.Exists(p.downDir) {
		if err := os.MkdirAll(p.downDir, os.ModePerm); err != nil {
			return err
		}
	}
	//存放分部的所有信息(#EXT-X-DISCONTINUITY)
	// parts := []string{}
	//存放分片的所有信息
	// segments := []string{}
	// segInfo := make(map[string]string)
	p.extLists = []string{}
	p.MEDIA_AUDIO_GROUP = make(map[string][]Audio)
	p.MEDIA_SUB_GROUP = make(map[string][]Subtitle)

	var (
		m3u8Content string = ""
		// m3u8Method     string   = ""
		// extMAP         []string = []string{"", ""}
		// extList        []string = [10]string{}
		// segIndex       int32    = 0
		// startIndex     int32    = 0
		// targetDuration int      = 0
		// totalDuration  float32  = 0
		// expectSegment  bool     = false
		// expectPlaylist bool     = false
		isEndlist bool = false
		// isAd           bool     = false
		// isM3u          bool     = false
	)
	if strings.Contains(p.m3u8Url, ".cntv.") {
		p.m3u8Url = strings.ReplaceAll(p.m3u8Url, "/h5e/", "/")
	}
	if strings.HasPrefix(p.m3u8Url, "http") {
		if strings.Contains(p.m3u8Url, "nfmovies.com/hls") {
			infbytes, err := download.HttpDownloadFileToBytes(p.m3u8Url, p.headers, 60)
			if err != nil {
				return err
			}
			m3u8Content = decode.NfmoviesDecryptM3u8(infbytes)
		} else if strings.Contains(p.m3u8Url, "hls.ddyunp.com/ddyun") || strings.Contains(p.m3u8Url, "hls.90mm.me/ddyun") {
			m3u8Url, err := decode.GetVaildM3u8Url(p.m3u8Url)
			if err != nil {
				return err
			}
			infbytes, err := download.HttpDownloadFileToBytes(m3u8Url, p.headers, 60)
			if err != nil {
				return err
			}
			m3u8Content = decode.DdyunDecryptM3u8(infbytes)
		} else {
			infbytes, err := download.GetWebSource(p.m3u8Url, p.headers, 60)
			if err != nil {
				return err
			}
			m3u8Content = tool.BytesToStr(infbytes)
		}
	} else if strings.HasPrefix(p.m3u8Url, "file:") {
		u, err := url.Parse(p.m3u8Url)
		if err != nil {
			return err
		}
		uri := u.Path
		sysType := runtime.GOOS
		uri = tool.IfString(sysType == "windows", uri[:1], uri)
		infbytes, err := tool.ReadFile(uri)
		if err != nil {
			return err
		}
		m3u8Content = tool.BytesToStr(infbytes)
	} else {
		infbytes, err := tool.ReadFile(p.m3u8Url)
		if err != nil {
			return err
		}
		m3u8Content = tool.BytesToStr(infbytes)
	}
	if m3u8Content == "" {
		return errors.New(lang.Lang.ParseError)
	}
	fmt.Println(m3u8Content)
	if strings.Contains(p.m3u8Url, "tlivecloud-playback-cdn.ysp.cctv.cn") && strings.Contains(p.m3u8Url, "endtime") {
		isEndlist = true
	}
	return nil
}
