package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	errorColor = "\x1b[41m"
	warnColor  = "\x1b[43m"
	infoColor  = "\x1b[44m"
	colorReset = "\x1b[0m"
)

var (
	filename   string
	filepath   string
	permission os.FileMode
	nowTime    string
	mu         sync.Mutex
)

func getDate() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func checkPerm() {
	fileInfo, err := os.Lstat("./")

	if err != nil {
		printLog(err, "ERROR")
	}

	fileMode := fileInfo.Mode()
	permission = fileMode & os.ModePerm
}

func fileNotExists(name string, t string) bool {
	f, err := os.Stat(name)
	if t == "dir" {
		return os.IsNotExist(err) || !f.IsDir()
	} else if t == "file" {
		return os.IsNotExist(err)
	}

	return false
}

func SetupLogger(path string) {
	checkPerm()

	if fileNotExists(path, "dir") {
		err := os.Mkdir(path, permission)
		if err != nil {
			printLog(err, "ERROR")
		}
	}

	filepath = path
	filename = fmt.Sprintf("%s.log", getDate())
}

func Warn(msg any) {
	log(msg, "WARN")
}

func Error(msg any) {
	log(msg, "ERROR")
}

func Info(msg any) {
	log(msg, "INFO")
}

func Panic(msg any) {
	log(msg, "PANIC")
}

func log(msg any, logType string) {
	printLog(msg, logType)
	writeLog(fmt.Sprintf("%s [%s] %s\n", nowTime, logType, msg))
	if logType == "PANIC" {
		os.Exit(1)
	}
}

func printLog(msg any, logType string) {
	nowTime = time.Now().Format("15:04:05")
	cc := ""
	switch logType {
	case "WARN":
		cc = warnColor
		break
	case "ERROR":
		cc = errorColor
		break
	case "INFO":
		cc = infoColor
		break
	case "PANIC":
		cc = errorColor
		break
	}

	fmt.Printf("%s %s%s%s %s\n", nowTime, cc, logType, colorReset, msg)
}

func writeLog(msg any) {
	mu.Lock()
	defer mu.Unlock()

	date := strings.Replace(filename, ".log", "", -1)
	if date != getDate() {
		filename = fmt.Sprintf("%s.log", date)
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", filepath, filename), os.O_CREATE|os.O_RDWR|os.O_APPEND, permission)
	if err != nil {
		printLog(err, "ERROR")
		return
	}

	_, err = file.WriteString(fmt.Sprint(msg))
	if err != nil {
		printLog(err, "ERROR")
		return
	}
}
