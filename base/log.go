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
		logTime   time.Time
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

func (this *Log) Init(fileName string) bool {
	//this.logFile = nil
	this.fileName = fileName
	this.logSuffix = "log"
	this.errSuffix = "err"
	log.SetPrefix(fmt.Sprintf("[%s]", this.fileName))
	return true
}

func (this *Log) GetSuffix(nType LG_TYPE) string {
	if nType == LG_WARN {
		return this.logSuffix
	} else {
		return this.errSuffix
	}
}

func (this *Log) Write(nType LG_TYPE) {
	this.WriteFile(nType)
	tTime := time.Now()
	this.logger[nType].SetPrefix(fmt.Sprintf("[%s][%04d-%02d-%02d %02d:%02d:%02d]", this.fileName, tTime.Year(), tTime.Month(), tTime.Day(),
		tTime.Hour(), tTime.Minute(), tTime.Second()))
}

func (this *Log) Println(v1 ...interface{}) {
	this.Write(LG_WARN)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r"
	this.logger[LG_WARN].Output(2, fmt.Sprintln(params...))
	log.Println(params...)
}

func (this *Log) Print(v1 ...interface{}) {
	this.Write(LG_WARN)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	this.logger[LG_WARN].Output(2, fmt.Sprint(params...))
	log.Print(params...)
}

func (this *Log) Printf(format string, params ...interface{}) {
	this.Write(LG_WARN)
	format += "\r\n"
	this.logger[LG_WARN].Output(2, fmt.Sprintf(format, params...))
	log.Printf(format, params...)
}

func (this *Log) Fatalln(v1 ...interface{}) {
	this.Write(LG_ERROR)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r"
	this.logger[LG_ERROR].Output(2, fmt.Sprintln(params...))
	log.Println(params...)
}

func (this *Log) Fatal(v1 ...interface{}) {
	this.Write(LG_ERROR)
	params := make([]interface{}, len(v1)+1)
	for i, v := range v1 {
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	this.logger[LG_ERROR].Output(2, fmt.Sprint(params...))
	log.Print(params...)
}

func (this *Log) Fatalf(format string, params ...interface{}) {
	this.Write(LG_ERROR)
	format += "\r\n"
	this.logger[LG_ERROR].Output(2, fmt.Sprintf(format, params...))
	log.Printf(format, params...)
}

func (this *Log) WriteFile(nType LG_TYPE) {
	var err error
	tTime := time.Now()
	if this.logTime.Year() != tTime.Year() ||
		this.logTime.Month() != tTime.Month() || this.logTime.Day() != tTime.Day() {
		this.loceker.Lock()
		if this.logFile[nType] != nil {
			defer this.logFile[nType].Close()
		}

		if PathExists(PATH) == false {
			os.Mkdir(PATH, os.ModeDir)
		}

		sFileName := fmt.Sprintf("%s/%s_%d%02d%02d.%s", PATH, this.fileName, tTime.Year(), tTime.Month(), tTime.Day(),
			this.GetSuffix(nType))

		if PathExists(sFileName) == false {
			os.Create(sFileName)
		}

		this.logFile[nType], err = os.OpenFile(sFileName, os.O_RDWR|os.O_APPEND, 0)
		if err != nil {
			log.Fatalf("open logfile[%s] error", sFileName)
		}

		this.logger[nType].SetOutput(this.logFile[nType])
		this.logger[nType].SetPrefix(fmt.Sprintf("[%s][%04d-%02d-%02d %02d:%02d:%02d]", this.fileName, tTime.Year(), tTime.Month(), tTime.Day(),
			tTime.Hour(), tTime.Minute(), tTime.Second()))
		this.logger[nType].SetFlags(log.Llongfile)
		Stat, _ := this.logFile[nType].Stat()
		if Stat != nil {
			this.logTime = Stat.ModTime()
		} else {
			this.logTime = time.Now()
		}
		this.loceker.Unlock()
	}
}
