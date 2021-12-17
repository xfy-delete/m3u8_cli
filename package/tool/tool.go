package tool

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/flytam/filenamify"

	"github.com/xfy520/m3u8_cli/package/lang"
	"github.com/xfy520/m3u8_cli/package/log"
)

// 判断是否是目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断是否是文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 判断目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 按任意键退出
func Pause() {
	fmt.Printf("\033[1;34;40m%s\033[0m", time.Now().Format("2006-01-02 15:04:05.000")+" ")
	fmt.Printf("\033[1;32;40m%s\033[0m", lang.Lang.AnyKey)
	b := make([]byte, 1)
	os.Stdin.Read(b)
	Exit()
}

// 错误判断
func Check(err error) {
	if err != nil {
		log.Error(err.Error())
		Pause()
	}
}

// 退出处理
func Exit() {
	log.Info("开始退出...")
	log.Info("执行清理...")
	log.Info("结束退出...")
	os.Exit(0)
}

// 字符串判断三目
func IfString(condition bool, trueVal string, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

// 解析字符串命令行
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

// 打开配置文件
func openArgsFile(file_path string) (string, error) {
	if Exists(file_path) && IsFile(file_path) {
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
	} else {
		return "", nil
	}
}

// 获取命令行配置参数
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

// 获取合法文件名
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

// 通过url获取文件名
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

// 复制文件内容工具函数
func CopyFile(srcFile string, destFile string) error {
	file1, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	var file2 *os.File
	if Exists(destFile) {
		file2, err = os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	} else {
		file2, err = os.Create(destFile)
	}
	if err != nil {
		return err
	}
	defer file1.Close()
	defer file2.Close()
	_, err = io.Copy(file2, file1)
	return err
}

// 读取文件工具函数
func ReadFile(file_path string) ([]byte, error) {
	if Exists(file_path) && IsFile(file_path) {
		f, err := os.Open(file_path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fd, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return fd, err
	}
	return nil, errors.New(lang.Lang.FilePathError + file_path)
}

// 字符串转字节数组
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// 字节数组转字符串
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 获取本地时间与世界标准时间差秒数或者毫秒数
func GetTimeStamp(bflag bool) int64 {
	now := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)
	utc := time.Now().UTC()
	ts := utc.Sub(now)
	if bflag {
		return int64(ts.Seconds())
	} else {
		return ts.Milliseconds()
	}
}
