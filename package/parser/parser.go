package parser

type Parser interface {
	SetBaseUrl(baseUrl string)
}

type parser struct {
	DownName  string
	DownDir   string
	M3u8Url   string
	KeyBase64 string
	KeyIV     string
	KeyFile   string
	Headers   string
	BaseUrl   string
}

func New(DownName string, DownDir string, M3u8Url string, KeyBase64 string, KeyIV string, KeyFile string, Headers string) Parser {
	BaseUrl := ""
	return &parser{DownName, DownDir, M3u8Url, KeyBase64, KeyIV, KeyFile, Headers, BaseUrl}
}

func (p *parser) SetBaseUrl(baseUrl string) {
	p.BaseUrl = baseUrl
}
