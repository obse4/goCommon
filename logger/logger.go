package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger

	logLevel = 0 //默认的LogLevel为0，即所有级别的日志都打印

	//
	logOut        io.Writer
	day           int
	dayChangeLock sync.RWMutex
	logFile       string
)

const (
	DebugLevel = iota //iota=0
	InfoLevel         //InfoLevel=iota, iota=1
	WarnLevel         //WarnLevel=iota, iota=2
	ErrorLevel        //ErrorLevel=iota, iota=3
)

// 日志配置
type LogConfig struct {
	// 默认的LogLevel为0，即所有级别的日志都打印
	LogLevel int
	// 是否输出日志文件
	LogOut bool
	// 文件所在文件夹位置, 默认执行文件所在目录的log文件夹下
	LogFile string
}

func InitLogger(conf *LogConfig) {
	conf.setLogLevel()
	conf.setLogFile()
}

func (l *LogConfig) setLogLevel() {
	logLevel = l.LogLevel
}

func (l *LogConfig) setLogFile() {

	now := time.Now()

	if l.LogOut {
		var (
			err  error
			file *os.File
		)
		if l.LogFile != "" {
			logFile = l.LogFile
		} else {
			logFile = filepath.Join(getExeDir(), "log")
		}
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)

		if err != nil {
			panic(err)
		}
		logOut = io.MultiWriter(os.Stdout, file)
	} else {
		logOut = os.Stdout
	}

	// 文件输出或控制台输出日志
	debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
	infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
	warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
	day = now.YearDay()
	dayChangeLock = sync.RWMutex{}

}

// 检查是否需要切换日志文件，如果需要则切换
func checkAndChangeLogfile() {
	if logOut != os.Stdout {
		dayChangeLock.Lock()
		defer dayChangeLock.Unlock()
		now := time.Now()
		if now.YearDay() == day {
			return
		}

		logOut.(*os.File).Close()
		postFix := now.Add(-24 * time.Hour).Format("20060102") //昨天的日期
		if err := os.Rename(logFile, logFile+"."+postFix); err != nil {
			fmt.Printf("append date postfix %s to log file %s failed: %v\n", postFix, logFile, err)
			return
		}
		var err error
		if logOut, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664); err != nil {
			fmt.Printf("create log file %s failed %v\n", logFile, err)
			return
		} else {
			debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
			infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
			warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
			errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
			day = now.YearDay()
		}
	}

}

func addPrefix() string {
	file, _, line := getLineNo()
	arr := strings.Split(file, "/")
	if len(arr) > 3 {
		arr = arr[len(arr)-3:]
	}
	return strings.Join(arr, "/") + ":" + strconv.Itoa(line)
}

func getLineNo() (string, string, int) {
	funcName, file, line, ok := runtime.Caller(3)
	if ok {
		return file, runtime.FuncForPC(funcName).Name(), line
	} else {
		return "", "", 0
	}
}

func Debug(format string, v ...interface{}) {
	if logLevel <= DebugLevel {
		checkAndChangeLogfile()
		debugLogger.Printf(addPrefix()+" "+format, v...)
	} else {
		fmt.Printf("[DEBUG] "+addPrefix()+" "+format+"\n", v...)
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		checkAndChangeLogfile()
		infoLogger.Printf(addPrefix()+" "+format, v...) //format末尾如果没有换行符会自动加上
	} else {
		fmt.Printf("[INFO] "+addPrefix()+" "+format+"\n", v...)
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= WarnLevel {
		checkAndChangeLogfile()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	} else {
		fmt.Printf("[WARN] "+addPrefix()+" "+format+"\n", v...)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ErrorLevel {
		checkAndChangeLogfile()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	} else {
		fmt.Printf("[ERROR] "+addPrefix()+" "+format+"\n", v...)
	}
}

func getExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	fmt.Println("执行文件所在目录:", exeDir)
	return exeDir
}
