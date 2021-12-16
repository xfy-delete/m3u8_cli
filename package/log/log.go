package log

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Level uint32

const (
	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
)

var countGuard sync.RWMutex
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

func InitLog(command string) error {
	logDir := path.Dir(LogFile)
	if !Exists(logDir) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	num := 1
	filenameall := path.Base(LogFile)
	filesuffix := path.Ext(LogFile)
	fileName := filenameall[0 : len(filenameall)-len(filesuffix)]
	for {
		if !Exists(LogFile) {
			break
		}
		LogFile = path.Join(logDir, fileName+"-"+string(rune(num)))
		num += 1
	}
	filePath := LogFile
	logs := []string{
		"Log " + time.Now().Format("2006-01-02") + "\r\n",
		"Save Path: " + logDir + "\r\n",
		"Task Start: " + time.Now().Format("2006-01-02 15:04:05") + "\r\n",
		"Task CommandLine: " + command + "\r\n",
	}
	var err error
	var file *os.File
	if Exists(filePath) {
		file, err = os.OpenFile(filePath, os.O_APPEND, os.ModePerm)
	} else {
		file, err = os.Create(filePath)
	}
	if err != nil {
		return err
	}
	_, err = io.WriteString(file, strings.Join(logs, ""))
	defer file.Close()
	return err
}

func WriteLine(log string, msg string) error {
	if !Exists(LogFile) {
		return nil
	}
	filePath := LogFile
	countGuard.Lock()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer countGuard.Unlock()
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(time.Now().Format("") + " / (" + msg + ") " + log)
	write.Flush()
	return nil
}

func WriteError(log string) error {
	return WriteLine(log, "ERROR")
}

func WriteInfo(log string) error {
	return WriteLine(log, "INFO")
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
