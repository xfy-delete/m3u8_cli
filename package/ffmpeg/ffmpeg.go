package ffmpeg

import (
	"errors"
	"os/exec"
	"path"
	"runtime"

	"github.com/xfy520/m3u8_cli/package/tool"
)

var ffmpeg_path = ""

func Init(ffmpegPath string) error {
	cmd := exec.Command("ffmpeg")
	if cmd.Path != "ffmpeg" {
		ffmpeg_path = cmd.Path
		return nil
	}
	sysType := runtime.GOOS
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		ffmpeg_path = path.Join(path.Dir(filename), tool.IfString(sysType == "windows", "ffmpeg.exe", "ffmpeg"))
	}
	if !tool.Exists(ffmpeg_path) {
		if ffmpegPath != "" {
			if !tool.IsFile(ffmpegPath) {
				ffmpeg_path = path.Join(ffmpegPath, tool.IfString(sysType == "windows", "ffmpeg.exe", "ffmpeg"))
			} else {
				ffmpeg_path = ffmpegPath
			}
		}
		if !tool.Exists(ffmpeg_path) {
			return errors.New("file path error")
		}
	}
	return nil
}
