package tool

import (
	"os"

	"github.com/xfy520/m3u8_cli/package/lang"
	"github.com/xfy520/m3u8_cli/package/log"
)

func Pause() {
	log.Info(lang.Lang.AnyKey)
	b := make([]byte, 1)
	os.Stdin.Read(b)
	ExitFunc()
}

func ExitFunc() {
	log.Info("开始退出...")
	log.Info("执行清理...")
	log.Info("结束退出...")
	os.Exit(0)
}

func IfString(condition bool, trueVal string, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}
