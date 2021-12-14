package log

import (
	"fmt"
	"time"
)

type Level uint32

const (
	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
)

var LogFile = ""

func PrintLine(msg string, Level Level) {
	switch Level {
	case InfoLevel:
		fmt.Printf("\033[1;34;40m%s\033[0m", time.Now().Format("2006-01-02 15:04:05.000")+" ")
		fmt.Printf("\033[1;32;40m%s\033[0m\n", msg)
	case WarnLevel:
		fmt.Printf("\033[1;34;40m%s\033[0m", time.Now().Format("2006-01-02 15:04:05.000")+" ")
		fmt.Printf("\033[1;33;40m%s\033[0m\n", msg)
	default:
		fmt.Printf("\033[1;34;40m%s\033[0m", time.Now().Format("2006-01-02 15:04:05.000")+" ")
		fmt.Printf("\033[1;31;40m%s\033[0m\n", msg)
	}
}

func Info(msg string) {
	PrintLine(msg, InfoLevel)
}

func Warn(msg string) {
	PrintLine(msg, WarnLevel)
}

func Error(msg string) {
	PrintLine(msg, ErrorLevel)
}

func WriteLine(msg string, Level Level) {

}

func WriteLineInfo(msg string) {
	WriteLine(msg, InfoLevel)
}

func WriteLineWarn(msg string) {
	WriteLine(msg, WarnLevel)
}

func WriteLineError(msg string) {
	WriteLine(msg, ErrorLevel)
}

func Show(msg string) {

}
