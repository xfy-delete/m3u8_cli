package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/xfy520/m3u8_cli/package/download/DownloadManager"
	"github.com/xfy520/m3u8_cli/package/ffmpeg"
	"github.com/xfy520/m3u8_cli/package/lang"
	"github.com/xfy520/m3u8_cli/package/log"
	"github.com/xfy520/m3u8_cli/package/parser"
	"github.com/xfy520/m3u8_cli/package/request"
	"github.com/xfy520/m3u8_cli/package/tool"
)

var (
	inputRetryCount int      = 30
	VERSION         string   = "1.0.0"
	BUILD_TIME      string   = "nil"
	GO_VERSION      string   = "1.17.1"
	maxThreads      int      = runtime.NumCPU()
	url             string   = ""
	minThreads      int      = 16
	retryCount      int      = 15
	timeOut         int      = 10
	baseUrl         string   = ""
	reqHeaders      string   = ""
	keyFile         string   = ""
	keyBase64       string   = ""
	keyIV           string   = ""
	muxSetJson      string   = "MUXSETS.json"
	muxFastStart    bool     = false
	delAfterDone    bool     = false
	parseOnly       bool     = false
	noMerge         bool     = false
	fileName        string   = ""
	workDir         string   = ""
	Args            []string = []string{}
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP:
				fmt.Println("终端控制进程结束(终端连接断开)", s)
				tool.Exit()
			case syscall.SIGINT:
				fmt.Println("用户发送INTR字符(Ctrl+C)触发", s)
				tool.Exit()
			case syscall.SIGTERM:
				fmt.Println("结束程序(可以被捕获、阻塞或忽略)", s)
				tool.Exit()
			case syscall.SIGQUIT:
				fmt.Println("用户发送QUIT字符(Ctrl+/)触发", s)
				tool.Exit()
			case syscall.SIGILL:
				fmt.Println("非法指令(程序错误、试图执行数据段、栈溢出等)", s)
				tool.Exit()
			default:
				fmt.Println("其他错误退出，不作处理，继续执行程序", s)
			}
		}
	}()
	runtime.GOMAXPROCS(maxThreads)
	log.Info("m3u8_cli version " + VERSION)
	log.Info("go version " + GO_VERSION)
	log.Info("build time " + BUILD_TIME)
	app := &cli.App{
		Name:    "m3u8_cli",
		Usage:   lang.Lang.Usage,
		Action:  run,
		Version: VERSION,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ffmpegPath",
				Aliases: []string{"fp"},
				Usage:   lang.Lang.FfmpegPathUsage,
				EnvVars: []string{"FFMPEG_PATH", "ffmpeg_path", "FFMPEGPATH", "FFMPEG-PATH", "FFMPEG"},
			},
			&cli.StringFlag{
				Name:    "workDir",
				Aliases: []string{"wd"},
				Usage:   lang.Lang.WorkDir,
			},
			&cli.StringFlag{
				Name:    "saveName",
				Aliases: []string{"sn"},
				Usage:   lang.Lang.SaveName,
			},
			&cli.StringFlag{
				Name:    "baseUrl",
				Aliases: []string{"bu"},
				Usage:   lang.Lang.BaseUrl,
			},
			&cli.StringFlag{
				Name:        "headers",
				Aliases:     []string{"hd"},
				Usage:       lang.Lang.Headers,
				DefaultText: "{}",
			},
			&cli.IntFlag{
				Name:        "maxThreads",
				Aliases:     []string{"maxT"},
				Usage:       lang.Lang.MaxThreads,
				DefaultText: "16",
			},
			&cli.IntFlag{
				Name:        "minThreads",
				Aliases:     []string{"minT"},
				Usage:       lang.Lang.MinThreads,
				DefaultText: "16",
			},
			&cli.IntFlag{
				Name:        "retryCount",
				Aliases:     []string{"rc"},
				Usage:       lang.Lang.RetryCount,
				DefaultText: "20",
			},
			&cli.IntFlag{
				Name:        "timeOut",
				Aliases:     []string{"to"},
				Usage:       lang.Lang.TimeOut,
				DefaultText: "10",
			},
			&cli.StringFlag{
				Name:        "muxSetJson",
				Aliases:     []string{"muxSJ"},
				Usage:       lang.Lang.MuxSetJson,
				DefaultText: "MUXSETS.json",
			},
			&cli.StringFlag{
				Name:    "useKeyFile",
				Aliases: []string{"ukf"},
				Usage:   lang.Lang.UseKeyFile,
			},
			&cli.StringFlag{
				Name:    "useKeyBase64",
				Aliases: []string{"ukb"},
				Usage:   lang.Lang.UseKeyBase64,
			},
			&cli.StringFlag{
				Name:    "useKeyIV",
				Aliases: []string{"uki"},
				Usage:   lang.Lang.UseKeyIV,
			},
			&cli.StringFlag{
				Name:    "downloadRange",
				Aliases: []string{"dr"},
				Usage:   lang.Lang.DownloadRange,
			},
			&cli.StringFlag{
				Name:    "liveRecDur",
				Aliases: []string{"ld"},
				Usage:   lang.Lang.LiveRecDur,
			},
			&cli.IntFlag{
				Name:        "stopSpeed",
				Aliases:     []string{"ss"},
				Usage:       lang.Lang.StopSpeed,
				DefaultText: "-999",
			},
			&cli.IntFlag{
				Name:        "maxSpeed",
				Aliases:     []string{"maxS"},
				Usage:       lang.Lang.MaxSpeed,
				DefaultText: "-999",
			},
			&cli.StringFlag{
				Name:    "proxyAddress",
				Aliases: []string{"pa"},
				Usage:   lang.Lang.ProxyAddress,
			},
			&cli.BoolFlag{
				Name:    "enableDelAfterDone",
				Aliases: []string{"eda"},
				Usage:   lang.Lang.EnableDelAfterDone,
			},
			&cli.BoolFlag{
				Name:    "enableMuxFastStart",
				Aliases: []string{"emfs"},
				Usage:   lang.Lang.EnableMuxFastStart,
			},
			&cli.BoolFlag{
				Name:    "enableBinaryMerge",
				Aliases: []string{"ebm"},
				Usage:   lang.Lang.EnableBinaryMerge,
			},
			&cli.BoolFlag{
				Name:    "enableParseOnly",
				Aliases: []string{"epo"},
				Usage:   lang.Lang.EnableParseOnly,
			},
			&cli.BoolFlag{
				Name:    "enableAudioOnly",
				Aliases: []string{"eao"},
				Usage:   lang.Lang.EnableAudioOnly,
			},
			&cli.BoolFlag{
				Name:    "disableDateInfo",
				Aliases: []string{"dd"},
				Usage:   lang.Lang.DisableDateInfo,
			},
			&cli.BoolFlag{
				Name:    "noMerge",
				Aliases: []string{"nm"},
				Usage:   lang.Lang.NoMerge,
			},
			&cli.BoolFlag{
				Name:    "noProxy",
				Aliases: []string{"np"},
				Usage:   lang.Lang.NoProxy,
			},
			&cli.BoolFlag{
				Name:    "disableIntegrityCheck",
				Aliases: []string{"dic"},
				Usage:   lang.Lang.DisableIntegrityCheck,
			},
		},
	}
	args, err := tool.GetArgs(os.Args, 1)
	if err != nil {
		log.Error(err.Error())
		tool.Pause()
	}
	if len(args) <= 1 {
		log.Error(lang.Lang.AragError)
		tool.Pause()
	} else {
		Args = args
		url = args[1]
		args = append(args[:1], args[2:]...)
		if err := app.Run(args); err != nil {
			log.Error(err.Error())
			tool.Pause()
		}
	}
}

func run(c *cli.Context) error {
	_, CurrentPath, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New(lang.Lang.ProjectPathError)
	}
	CurrentPath = path.Dir(CurrentPath)
	if ffmpeg.Init(c.String("ffmpegPath")) != nil {
		fmt.Printf("\033[1;31;40m%s\033[0m\n", lang.Lang.FfmpegLost)
		fmt.Printf("\033[1;31;40m%s\033[0m\n\n", lang.Lang.FfmpegTip)
		log.Info("http://ffmpeg.org/download.html")
		tool.Pause()
	}
	if url == "" {
		return errors.New(lang.Lang.UrlError)
	}

	delAfterDone = c.Bool("enableDelAfterDone")
	fmt.Println(delAfterDone)
	parseOnly = c.Bool("enableParseOnly")
	fmt.Println(parseOnly)
	BinaryMerge := c.Bool("enableBinaryMerge")
	fmt.Println(BinaryMerge)
	WriteDate := !c.Bool("disableDateInfo")
	fmt.Println(WriteDate)
	noMerge = c.Bool("noMerge")
	fmt.Println(noMerge)

	request.NoProxy = c.Bool("noProxy")

	if c.String("proxyAddress") != "" && strings.HasPrefix(c.String("proxyAddress"), "http://") {
		request.UseProxyAddress = c.String("proxyAddress")
	}

	if c.String("headers") != "" {
		reqHeaders = c.String("headers")
	}

	muxFastStart = c.Bool("enableMuxFastStart")
	fmt.Println(muxFastStart)
	DisableIntegrityCheck := c.Bool("disableIntegrityCheck")
	fmt.Println(DisableIntegrityCheck)
	if c.Bool("enableAudioOnly") {
		VIDEO_TYPE := "IGNORE"
		fmt.Println(VIDEO_TYPE)
	}
	muxSetJson = c.String("muxSetJson")
	fmt.Println(muxSetJson)

	if c.String("workDir") != "" {
		workDir = c.String("workDir")
	} else {
		workDir = path.Join(CurrentPath, "Downloads")
	}

	if c.String("saveName") != "" {
		fileName = tool.GetFileName(c.String("saveName"))
	} else {
		fileName = tool.GetUrlFileName(url) + "_" + time.Now().Format("2006-01-02.15-04-05")
	}

	if c.String("useKeyFile") != "" {
		keyFile = c.String("useKeyFile")
	}

	if c.String("useKeyBase64") != "" {
		keyBase64 = c.String("useKeyBase64")
	}

	if c.String("useKeyIV") != "" {
		keyIV = c.String("useKeyIV")
	}

	if c.Int("stopSpeed") != -999 {
		STOP_SPEED := c.Int("stopSpeed")
		fmt.Println(STOP_SPEED)
	}

	if c.Int("maxSpeed") != -999 {
		MAX_SPEED := c.Int("maxSpeed")
		fmt.Println(MAX_SPEED)
	}

	if c.String("baseUrl") != "" {
		baseUrl = c.String("baseUrl")
	}

	maxThreads = c.Int("maxThreads")
	fmt.Println(maxThreads)

	minThreads = c.Int("minThreads")
	fmt.Println(minThreads)

	retryCount = c.Int("retryCount")
	fmt.Println(retryCount)

	timeOut = c.Int("timeOut")
	fmt.Println(timeOut)

	if c.String("liveRecDur") != "" {
		reg := regexp.MustCompile(`(\d+):(\d+):(\d+)`)
		liveRecDur := c.String("liveRecDur")
		params := reg.FindStringSubmatch(liveRecDur)
		for _, param := range params {
			fmt.Println(param)
		}
		// int HH = Convert.ToInt32(reg2.Match(t).Groups[1].Value)
		// int MM = Convert.ToInt32(reg2.Match(t).Groups[2].Value)
		// int SS = Convert.ToInt32(reg2.Match(t).Groups[3].Value)
		// HLSLiveDownloader.REC_DUR_LIMIT = SS + MM * 60 + HH * 60 * 60
	}
	if c.String("downloadRange") != "" {
		downloadRange := c.String("downloadRange")
		if strings.Contains(downloadRange, ":") {
			reg := regexp.MustCompile(`((\d+):(\d+):(\d+))?-((\d+):(\d+):(\d+))?`)
			params := reg.FindStringSubmatch(downloadRange)
			for _, param := range params {
				fmt.Println(param)
			}
			// Parser.DurStart = reg2.Match(p).Groups[1].Value;
			// Parser.DurEnd = reg2.Match(p).Groups[5].Value;
			// if (Parser.DurEnd == "00:00:00") Parser.DurEnd = "";
			// Parser.DelAd = false;
		} else {
			reg := regexp.MustCompile(`(\d*)-(\d*)`)
			params := reg.FindStringSubmatch(downloadRange)
			fmt.Println(params)
			// if (!string.IsNullOrEmpty(reg.Match(p).Groups[1].Value))
			// {
			//     Parser.RangeStart = Convert.ToInt32(reg.Match(p).Groups[1].Value);
			//     Parser.DelAd = false;
			// }
			// if (!string.IsNullOrEmpty(reg.Match(p).Groups[2].Value))
			// {
			//     Parser.RangeEnd = Convert.ToInt32(reg.Match(p).Groups[2].Value);
			//     Parser.DelAd = false;
			// }
		}
	}
	return input(CurrentPath)
}

func input(CurrentPath string) error {
	if inputRetryCount == 0 {
		log.Error(lang.Lang.InputRetryCount)
		tool.Pause()
	}

	if strings.Contains(url, "twitcasting") && strings.Contains(url, "/fmp4/") {
		DownloadManager.BinaryMerge = true
	}

	// m3u8Content := ""
	// isVOD := true
	if len(workDir) >= 300 {
		return errors.New(lang.Lang.DirLongError)
	} else if len(workDir)+len(fileName) >= 300 {
		for {
			fileName = fileName[0 : len(fileName)-1]
			if len(workDir)+len(fileName) < 300 {
				break
			}
		}
	}

	log.Info(lang.Lang.FileName + fileName)
	log.Info(lang.Lang.SavePath + path.Join(workDir, fileName))
	parser := parser.New(fileName, path.Join(workDir, fileName), url, keyBase64, keyIV, keyFile, reqHeaders)
	if baseUrl != "" {
		parser.SetBaseUrl(baseUrl)
	}
	log.LogFile = path.Join(CurrentPath, "Logs", time.Now().Format("2006-01-02_15-04-05.000")+".log")
	if err := log.InitLog(url + " " + strings.Join(append(Args[:0], Args[1:]...), " ")); err != nil {
		return err
	}
	tool.Check(log.WriteInfo(lang.Lang.StartParsing + url))
	log.Warn(lang.Lang.StartParsing + url)
	if strings.HasSuffix(url, ".json") && tool.Exists(url) {
		if !tool.Exists(path.Join(workDir, fileName)) {
			if err := os.MkdirAll(path.Join(workDir, fileName), os.ModePerm); err != nil {
				return err
			}
		}
		if err := tool.CopyFile(url, path.Join(workDir, fileName, "meta.json")); err != nil {
			return err
		}
	} else {
		parser.Parse()
	}
	return nil
}
