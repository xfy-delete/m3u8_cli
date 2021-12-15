package tool

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/flytam/filenamify"

	"github.com/xfy520/m3u8_cli/package/lang"
	"github.com/xfy520/m3u8_cli/package/log"
)

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

func Pause() {
	fmt.Printf("\033[1;34;40m%s\033[0m", time.Now().Format("2006-01-02 15:04:05.000")+" ")
	fmt.Printf("\033[1;32;40m%s\033[0m", lang.Lang.AnyKey)
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

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func ParseCommandLine(command string) ([]string, error) {
	var args []string
	args = append(args, os.Args[0])
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	rune_command := []rune(command)
	for i := 0; i < len(rune_command); i++ {
		c := rune_command[i]
		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}
		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}
		if c == '\\' {
			escapeNext = true
			continue
		}
		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}
		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}
		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}
	if state == "quotes" {
		return nil, errors.New(lang.Lang.CommandError + command)
	}
	if current != "" {
		args = append(args, current)
	}
	return args, nil
}

func openArgsFile(file_path string) (string, error) {
	file, err := os.Open(file_path)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	args := []string{}
	input := bufio.NewScanner(file)
	for input.Scan() {
		arg := input.Text()
		if arg != "" {
			if strings.HasPrefix(arg, "--") {
				args = append(args, arg)
			} else {
				args = append(args, "--"+arg)
			}
		}
	}
	file.Close()
	return strings.Join(args, " "), nil
}

func GetArgs(args []string, skip int) ([]string, error) {
	argsLen := len(args)
	if argsLen == 1 {
		fmt.Printf("\n\033[1;36;40m%s\033[0m", "m3u8_cli")
		fmt.Print(">")
		inputReader := bufio.NewReader(os.Stdin)
		input, _, err := inputReader.ReadLine()
		if err == nil {
			args, err := ParseCommandLine(string(input))
			if err != nil {
				return nil, err
			}
			argsLen = len(args)
			if argsLen == 2 || argsLen == 3 || argsLen == 4 || argsLen == 5 {
				return GetArgs(args, 2)
			} else {
				return args, nil
			}
		} else {
			return nil, err
		}
	}
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("system error")
	}
	args_file := path.Join(path.Dir(filename), "m3ub_cli")
	if argsLen == 2 {
		args_str, err := openArgsFile(args_file)
		if err != nil {
			return nil, err
		}
		args, err = ParseCommandLine(args[1] + " " + args_str)
		return args, err
	}
	if argsLen == 4 && (args[2] == "--saveName" || args[2] == "--sn") {
		args_str, err := openArgsFile(args_file)
		if err != nil {
			return nil, err
		}
		return ParseCommandLine(args[1] + " --sn " + args[3] + " " + args_str)
	}
	return args, nil
}

func GetFileName(file_name string) string {
	output, err := filenamify.Filenamify(file_name, filenamify.Options{
		Replacement: ".",
	})
	if err != nil {
		log.Error(lang.Lang.FileNameError)
		return ""
	}
	return output
}

func GetUrlFileName(url string) string {
	if Exists(url) && IsFile(url) {
		return path.Base(url)
	}
	urls := strings.Split(url, "/")
	if len(urls) > 0 {
		return strings.Replace(strings.Split(urls[len(urls)-1], "?")[0], ".m3u8", "", 1)
	} else {
		return time.Now().Format("2006-01-02.15-04-05")
	}
}
