package cm

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/rpc"
	"os"
	"sync"
	"time"
)

type (
	FileRead func() //reload
	FileInfo struct {
		Name    string
		ModTime time.Time
		Call    FileRead //call reload
	}

	FileMonitor struct {
		actor.Actor
		filesMap  map[string]*FileInfo
		callsMap  map[string]FileRead
		callsLock sync.Locker
	}

	IFileMonitor interface {
		actor.IActor
		addFile(string, FileRead)
		delFile(string)
		update()
		AddFile(string, FileRead)
	}
)

func (f *FileMonitor) Init() {
	f.Actor.Init()
	f.filesMap = map[string]*FileInfo{}
	f.callsMap = map[string]FileRead{}
	f.callsLock = &sync.Mutex{}
	f.RegisterTimer(1*time.Second, f.update)
	actor.MGR.RegisterActor(f)
	f.Actor.Start()
}

func (f *FileMonitor) AddFile(fileName string, pFunc FileRead) {
	f.callsLock.Lock()
	defer f.callsLock.Unlock()
	f.callsMap[fileName] = pFunc
	f.SendMsg(rpc.RpcHead{SendType: rpc.SEND_LOCAL}, "Addfile", fileName)
}

func (f *FileMonitor) addFile(fileName string, pFunc FileRead) {
	fileInfo, err := os.Stat(fileName)
	if err == nil {
		f.filesMap[fileName] = &FileInfo{fileInfo.Name(), fileInfo.ModTime(), pFunc}
	}
}

func (f *FileMonitor) delFile(fileName string) {
	delete(f.filesMap, fileName)
}

func (f *FileMonitor) update() {
	for i, v := range f.filesMap {
		fileInfo, err := os.Stat(i)
		if err == nil && v.ModTime != fileInfo.ModTime() {
			v.Call()
			v.ModTime = fileInfo.ModTime()
			base.LOG.Println(fmt.Sprintf("file [%s] reload", v.Name))
		}
	}
}

func (f *FileMonitor) Addfile(ctx context.Context, fileName string) {
	f.callsLock.Lock()
	defer f.callsLock.Unlock()
	if pFunc, isOk := f.callsMap[fileName]; isOk {
		f.addFile(fileName, pFunc)
	}
}

func (f *FileMonitor) Delfile(ctx context.Context, fileName string) {
	f.delFile(fileName)
}

var FILEMONITOR FileMonitor
