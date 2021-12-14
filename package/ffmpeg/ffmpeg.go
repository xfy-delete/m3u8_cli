package ffmpeg

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/xfy520/m3u8_cli/package/tool"
)

var ffmpeg_path = ""

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Init(ffmpegPath string) error {
	cmd := exec.Command("ffmpeg")
	if cmd.Path != "ffmpeg" {
		ffmpeg_path = cmd.Path
	}
	if Exists(ffmpeg_path) {
		return nil
	}
	sysType := runtime.GOOS
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		ffmpeg_path = path.Join(path.Dir(filename), tool.IfString(sysType == "windows", "ffmpeg.exe", "ffmpeg"))
	}
	if !Exists(ffmpeg_path) {
		if ffmpegPath != "" {
			if !IsFile(ffmpegPath) {
				ffmpeg_path = path.Join(ffmpegPath, tool.IfString(sysType == "windows", "ffmpeg.exe", "ffmpeg"))
			} else {
				ffmpeg_path = ffmpegPath
			}
		}
		if !Exists(ffmpeg_path) {
			return errors.New("file path error")
		}
	}
	return nil
}
