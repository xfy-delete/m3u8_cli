package parser

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/xfy520/m3u8_cli/package/decode"
	"github.com/xfy520/m3u8_cli/package/download"
	"github.com/xfy520/m3u8_cli/package/ffmpeg"
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
	DownName          string
	DownDir           string
	M3u8Url           string
	KeyBase64         string
	KeyIV             string
	KeyFile           string
	Headers           string
	BaseUrl           string
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
	BaseUrl      = ""
	m3u8SavePath = ""
	jsonSavePath = ""
)

var (
	extLists          = []string{}
	MEDIA_AUDIO_GROUP = make(map[string][]Audio)
	MEDIA_SUB_GROUP   = make(map[string][]Subtitle)
)

func New(DownName string, DownDir string, M3u8Url string, KeyBase64 string, KeyIV string, KeyFile string, Headers string) Parser {
	return &parser{DownName, DownDir, M3u8Url, KeyBase64, KeyIV, KeyFile,
		Headers, BaseUrl, m3u8SavePath, jsonSavePath, extLists, MEDIA_AUDIO_GROUP,
		MEDIA_SUB_GROUP,
	}
}

func (p *parser) SetBaseUrl(baseUrl string) {
	p.BaseUrl = baseUrl
}

func (p *parser) Parse() error {
	ffmpeg.REC_TIME = ""
	p.m3u8SavePath = path.Join(p.DownDir, "raw.m3u8")
	if !tool.Exists(p.DownDir) {
		if err := os.MkdirAll(p.DownDir, os.ModePerm); err != nil {
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
		// isEndlist      bool     = false
		// isAd           bool     = false
		// isM3u          bool     = false
	)
	if strings.Contains(p.M3u8Url, ".cntv.") {
		p.M3u8Url = strings.ReplaceAll(p.M3u8Url, "/h5e/", "/")
	}
	if strings.HasPrefix(p.M3u8Url, "http") {
		if strings.Contains(p.M3u8Url, "nfmovies.com/hls") {
			infbytes, err := download.HttpDownloadFileToBytes(p.M3u8Url, p.Headers, 6)
			if err != nil {
				return err
			}
			m3u8Content = decode.NfmoviesDecryptM3u8(infbytes)
		}
	}
	fmt.Println(m3u8Content)
	return nil
}
