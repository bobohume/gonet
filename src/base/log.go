package base

import (
	"os"
	"log"
	"time"
	"fmt"
	"sync"
)

const (
	LG_WARN = iota
	LG_ERROR = iota
	PATH = "Log"
)

type (
	CLog struct {
		log.Logger
		m_pFile *os.File
		m_Time time.Time
		m_FileName string
		m_LogSuffix string
		m_ErrSuffix string
		m_Loceker sync.Mutex
	}

	ILog interface {
		Init(string) bool
		Write(int)
		WreiteFile(int)
	}
)

var(
	G_Log CLog
)

func (this *CLog) Init(fileName string) bool {
	this.m_pFile = nil
	this.m_FileName = fileName
	this.m_LogSuffix = "log"
	this.m_ErrSuffix = "err"
	return true
}

func (this *CLog)Write(nType int){
	this.WriteFile(nType)
	tTime := time.Now()
	this.SetPrefix(fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",tTime.Year(), tTime.Month(), tTime.Day(),
		tTime.Hour(), tTime.Minute(), tTime.Second()))
}

func (this *CLog) Println(v1 ...interface{}) {
	this.Write(LG_WARN)
	params := make([]interface{}, len(v1) + 1)
	for i,v := range v1{
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	this.Logger.Println(params...)
	log.Println(params...)
}

func (this *CLog) Print(v1 ...interface{}) {
	this.Write(LG_WARN)
	params := make([]interface{}, len(v1) + 1)
	for i,v := range v1{
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	this.Logger.Print(params...)
	log.Print(params...)
}

func (this *CLog) Printf(format string, params ...interface{}) {
	this.Write(LG_WARN)
	format += "\r\n";
	this.Logger.Printf(format, params...)
	log.Printf(format,params...)
}

func (this *CLog) Fatalln(v1 ...interface{}) {
	this.Write(LG_ERROR)
	params := make([]interface{}, len(v1) + 1)
	for i,v := range v1{
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	this.Logger.Println(params...)
}

func (this *CLog) Fatal(v1 ...interface{}) {
	this.Write(LG_ERROR)
	params := make([]interface{}, len(v1) + 1)
	for i,v := range v1{
		params[i] = v
	}
	params[len(v1)] = "\r\n"
	this.Logger.Print(params...)
}

func (this *CLog) Fatalf(format string, params ...interface{}) {
	this.Write(LG_ERROR)
	format += "\r\n"
	this.Logger.Printf(format, params...)
}

func (this *CLog) WriteFile(nType int){
	var err error
	tTime := time.Now()
	if (this.m_pFile == nil ||  this.m_Time.Year() != tTime.Year() ||
	this.m_Time.Month() != tTime.Month() || this.m_Time.Day() != tTime.Day()){
		this.m_Loceker.Lock()
		if this.m_pFile != nil {
			defer this.m_pFile.Close()
		}

		if PathExists(PATH) == false {
			os.Mkdir(PATH, os.ModeDir)
		}

		GetSuffix := func(nType int) string{
			if nType == LG_WARN{
				return this.m_LogSuffix
			}else{
				return this.m_ErrSuffix
			}
		}

		sFileName := fmt.Sprintf("%s/%s_%d%02d%02d.%s", PATH, this.m_FileName,tTime.Year(), tTime.Month(), tTime.Day(),
			GetSuffix(nType))

		if PathExists(sFileName) == false{
			os.Create(sFileName);
		}

		this.m_pFile,err  = os.OpenFile(sFileName, os.O_RDWR, 0)
		if err != nil {
			log.Fatalln("open logfile[%s] error", sFileName)
		}

		this.SetOutput(this.m_pFile)
		this.SetPrefix(fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d]",tTime.Year(), tTime.Month(), tTime.Day(),
			tTime.Hour(), tTime.Minute(), tTime.Second()))
		this.SetFlags(log.Llongfile)
		Stat,_ := this.m_pFile.Stat()
		if Stat != nil {
			this.m_Time = Stat.ModTime()
		}else{
			this.m_Time = time.Now()
		}
		this.m_Loceker.Unlock()
	}
}
