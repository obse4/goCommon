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

// 定义日志等级类型
type LogLeveL uint32

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger

	//默认的LogLevel为0，即所有级别的日志都打印
	logLevel LogLeveL = 0

	// 文件输出模式writer
	logOut io.Writer
	// 当年天数
	day int
	// 日志文件锁
	dayChangeLock sync.RWMutex
	// 日志文件指定路径
	logFile string
	// 文件保留天数
	stayDay int
)

const (
	DebugLevel LogLeveL = iota //iota=0
	InfoLevel                  //InfoLevel=iota, iota=1
	WarnLevel                  //WarnLevel=iota, iota=2
	ErrorLevel                 //ErrorLevel=iota, iota=3
)

// 日志配置
type LogConfig struct {
	// 默认的LogLevel为0，即所有级别的日志都打印
	LogLevel LogLeveL
	// 是否输出日志文件
	LogOut bool
	// 文件所在文件夹位置, 默认执行文件所在目录的log文件夹下
	// 使用绝对路径
	LogFile string
	// 输出日志文件保存时间, 默认0，即全部保存，1则只保留当日
	StayDay int
}

// 初始化日志
func InitLogger(conf *LogConfig) {
	conf.setLogLevel()
	conf.setLogFile()
	conf.setStayDay()
}

func (l *LogConfig) setLogLevel() {
	logLevel = l.LogLevel
}

func (l *LogConfig) setStayDay() {
	stayDay = l.StayDay
}

func (l *LogConfig) setLogFile() {

	now := time.Now()

	if l.LogOut {
		// 输出日志文件
		var (
			err  error
			file *os.File
		)
		if l.LogFile != "" {
			// 用户指定文件路径
			logFile = l.LogFile
		} else {
			// 默认文件路径 执行文件所在目录下log文件
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
func checkAndChangeLogFile() {
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
			Warn("append date postfix %s to log file %s failed: %v\n", postFix, logFile, err)
			return
		}

		if stayDay > 0 {
			removeHistoryLogFile()
		}

		var err error
		if logOut, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664); err != nil {
			Warn("create log file %s failed %v\n", logFile, err)
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
		checkAndChangeLogFile()
		debugLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= InfoLevel {
		checkAndChangeLogFile()
		infoLogger.Printf(addPrefix()+" "+format, v...) //format末尾如果没有换行符会自动加上
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= WarnLevel {
		checkAndChangeLogFile()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ErrorLevel {
		checkAndChangeLogFile()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	}
}

// 获取当前执行文件所在文件夹位置
func getExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	return exeDir
}

// 获取历史日志文件路径列表
func getLogFileHistoryList() []string {
	targetFilepath := logFile

	// 获取文件所在的目录
	dir := filepath.Dir(targetFilepath)

	// 列出目录下的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var fileList []string

	// 遍历并打印文件名
	for _, file := range files {
		if strings.Contains(filepath.Join(dir, file.Name()), fmt.Sprintf("%s.", targetFilepath)) {
			fileList = append(fileList, filepath.Join(dir, file.Name()))
			fmt.Println(filepath.Join(dir, file.Name()))
		}

	}
	return fileList
}

// 删除文件
func removeFile(path string) {
	// 删除文件
	err := os.Remove(path)
	if err != nil {
		return
	}
}

// 删除保留日之前日志
func removeHistoryLogFile() {
	historyPathList := getLogFileHistoryList()

	for _, path := range historyPathList {
		postFix := strings.ReplaceAll(path, fmt.Sprintf("%s.", logFile), "")
		filePostFix, err := strconv.Atoi(postFix)
		if err != nil {
			continue
		}

		stayStartPostFix, err := strconv.Atoi(time.Now().AddDate(0, 0, -stayDay).Format("20060102"))
		if err != nil {
			continue
		}
		if filePostFix <= stayStartPostFix {
			removeFile(path)
		}
	}
}
