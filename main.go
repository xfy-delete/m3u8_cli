package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/xfy520/m3u8_cli/package/ffmpeg"
	"github.com/xfy520/m3u8_cli/package/lang"
	"github.com/xfy520/m3u8_cli/package/log"
	"github.com/xfy520/m3u8_cli/package/tool"
)

var (
	VERSION    string = "1.0.0"
	BUILD_TIME string = "nil"
	GO_VERSION string = "1.17.1"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP:
				fmt.Println("终端控制进程结束(终端连接断开)", s)
				tool.ExitFunc()
			case syscall.SIGINT:
				fmt.Println("用户发送INTR字符(Ctrl+C)触发", s)
				tool.ExitFunc()
			case syscall.SIGTERM:
				fmt.Println("结束程序(可以被捕获、阻塞或忽略)", s)
				tool.ExitFunc()
			case syscall.SIGQUIT:
				fmt.Println("用户发送QUIT字符(Ctrl+/)触发", s)
				tool.ExitFunc()
			case syscall.SIGILL:
				fmt.Println("非法指令(程序错误、试图执行数据段、栈溢出等)", s)
				tool.ExitFunc()
			default:
				fmt.Println("其他错误退出，不作处理，继续执行程序", s)
			}
		}
	}()
	CurrentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error(lang.Lang.ProjectPathError)
		tool.Pause()
	}
	app := cli.NewApp()
	app.Name = "m3u8_cli"
	app.Usage = lang.Lang.Usage
	app.Action = run
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "ffmpegPath",
			Aliases: []string{"fp"},
			Usage:   lang.Lang.FfmpegPathUsage,
			EnvVars: []string{"FFMPEG_PATH", "ffmpeg_path", "FFMPEGPATH", "FFMPEG-PATH", "FFMPEG"},
		},
		&cli.StringFlag{
			Name:        "workDir",
			Aliases:     []string{"wd"},
			Usage:       lang.Lang.WorkDir,
			DefaultText: path.Join(CurrentPath, "Downloads"),
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
			Name:    "stopSpeed",
			Aliases: []string{"ss"},
			Usage:   lang.Lang.StopSpeed,
		},
		&cli.IntFlag{
			Name:    "maxSpeed",
			Aliases: []string{"maxS"},
			Usage:   lang.Lang.MaxSpeed,
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

func run(c *cli.Context) error {
	if ffmpeg.Init(c.String("ffmpeg.path")) != nil {
		fmt.Printf("\033[1;31;40m%s\033[0m\n", lang.Lang.FfmpegLost)
		fmt.Printf("\033[1;31;40m%s\033[0m\n\n", lang.Lang.FfmpegTip)
		log.Info("http://ffmpeg.org/download.html")
		tool.Pause()
	}
	log.Info("m3u8_cli version " + VERSION)
	log.Info("go version " + GO_VERSION)
	log.Info("build time " + BUILD_TIME)
	log.Info(lang.Lang.StartUp)
	return nil
}
