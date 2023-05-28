package base

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LG_TYPE int

const (
	LG_WARN  LG_TYPE = iota
	LG_ERROR LG_TYPE = iota
	LG_MAX   LG_TYPE = iota
)
const (
	PATH = "log"
)

type (
	Log struct {
		logger    [LG_MAX]log.Logger
		logFile   [LG_MAX]*os.File
		logTime   [LG_MAX]time.Time
		fileName  string
		logSuffix string
		errSuffix string
		loceker   sync.Mutex
	}

	ILog interface {
		Init(string) bool
		Write(LG_TYPE)
		WriteFile(LG_TYPE)
		Println(...interface{})
		Print(...interface{})
		Printf(string, ...interface{})
		Fatalln(...interface{})
		Fatal(...interface{})
		Fatalf(string, ...interface{})
	}
)

var (
	LOG Log
)

func (l *Log) Init(fileName string) bool {
	//l.logFile = nil
	l.fileName = fileName
	l.logSuffix = "log"
	l.errSuffix = "err"
	log.SetPrefix(fmt.Sprintf("[%s]", l.fileName))
	return true
}

func (l *Log) GetSuffix(nType LG_TYPE) string {
	if nType == LG_WARN {
		return l.logSuffix
	} else {
		return l.errSuffix
	}
}

func (l *Log) Write(nType LG_TYPE) {
	l.WriteFile(nType)
	tTime := time.Now()
	l.logger[nType].SetPrefix(fmt.Sprintf("[%s][%04d-%02d-%02d %02d:%02d:%02d]", l.fileName, tTime.Year(), tTime.Month(), tTime.Day(),
		tTime.Hour(), tTime.Minute(), tTime.Second()))
}

func (l *Log) Println(v1 ...interface{}) {
	l.Write(LG_WARN)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r"
	l.logger[LG_WARN].Output(2, fmt.Sprintln(params...))
	log.Println(params...)
}

func (l *Log) Print(v1 ...interface{}) {
	l.Write(LG_WARN)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	l.logger[LG_WARN].Output(2, fmt.Sprint(params...))
	log.Print(params...)
}

func (l *Log) Printf(format string, params ...interface{}) {
	l.Write(LG_WARN)
	format += "\r\n"
	l.logger[LG_WARN].Output(2, fmt.Sprintf(format, params...))
	log.Printf(format, params...)
}

func (l *Log) Fatalln(v1 ...interface{}) {
	l.Write(LG_ERROR)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r"
	l.logger[LG_ERROR].Output(2, fmt.Sprintln(params...))
	log.Println(params...)
}

func (l *Log) Fatal(v1 ...interface{}) {
	l.Write(LG_ERROR)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	l.logger[LG_ERROR].Output(2, fmt.Sprint(params...))
	log.Print(params...)
}

func (l *Log) Fatalf(format string, params ...interface{}) {
	l.Write(LG_ERROR)
	format += "\r\n"
	l.logger[LG_ERROR].Output(2, fmt.Sprintf(format, params...))
	log.Printf(format, params...)
}

func (l *Log) WriteFile(nType LG_TYPE) {
	var err error
	tTime := time.Now()
	if l.logTime[nType].Year() != tTime.Year() ||
		l.logTime[nType].Month() != tTime.Month() || l.logTime[nType].Day() != tTime.Day() {
		l.loceker.Lock()
		if l.logFile[nType] != nil {
			defer l.logFile[nType].Close()
		}

		if PathExists(PATH) == false {
			os.Mkdir(PATH, os.ModeDir)
		}

		sFileName := fmt.Sprintf("%s/%s_%d%02d%02d.%s", PATH, l.fileName, tTime.Year(), tTime.Month(), tTime.Day(),
			l.GetSuffix(nType))

		if PathExists(sFileName) == false {
			os.Create(sFileName)
		}

		l.logFile[nType], err = os.OpenFile(sFileName, os.O_RDWR|os.O_APPEND, 0)
		if err != nil {
			log.Fatalf("open logfile[%s] error", sFileName)
		}

		l.logger[nType].SetOutput(l.logFile[nType])
		l.logger[nType].SetPrefix(fmt.Sprintf("[%s][%04d-%02d-%02d %02d:%02d:%02d]", l.fileName, tTime.Year(), tTime.Month(), tTime.Day(),
			tTime.Hour(), tTime.Minute(), tTime.Second()))
		l.logger[nType].SetFlags(log.Llongfile)
		Stat, _ := l.logFile[nType].Stat()
		if Stat != nil {
			l.logTime[nType] = Stat.ModTime()
		} else {
			l.logTime[nType] = time.Now()
		}
		l.loceker.Unlock()
	}
}
