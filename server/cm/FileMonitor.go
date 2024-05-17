package cm

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/rpc"
	"os"
	"time"
	"unsafe"
)

type (
	FileRead func() //reload
	FileInfo struct {
		Info os.FileInfo
		Call FileRead //call reload
	}

	FileMonitor struct {
		actor.Actor
		filesMap map[string]*FileInfo
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
	f.RegisterTimer(3*time.Second, f.update)
	actor.MGR.RegisterActor(f)
	f.Actor.Start()
}

func (f *FileMonitor) AddFile(fileName string, pFunc FileRead) {
	ponit := unsafe.Pointer(&pFunc)
	f.SendMsg(rpc.RpcHead{}, "Addfile", fileName, (*int64)(ponit))
}

func (f *FileMonitor) addFile(fileName string, pFunc FileRead) {
	file, err := os.Open(fileName)
	if err == nil {
		defer file.Close()
		fileInfo, err := file.Stat()
		if err == nil {
			f.filesMap[fileName] = &FileInfo{fileInfo, pFunc}
		}
	}
}

func (f *FileMonitor) delFile(fileName string) {
	delete(f.filesMap, fileName)
}

func (f *FileMonitor) update() {
	for i, v := range f.filesMap {
		file, err := os.Open(i)
		if err == nil {
			defer file.Close()
			fileInfo, err := file.Stat()
			if err == nil && v.Info.ModTime() != fileInfo.ModTime() {
				v.Call()
				v.Info = fileInfo
				fmt.Println(fmt.Sprintf("file [%s] reload", v.Info.Name()))
			}
		}
	}
}

func (f *FileMonitor) Addfile(ctx context.Context, fileName string, p *int64) {
	pFunc := (*FileRead)(unsafe.Pointer(p))
	f.addFile(fileName, *pFunc)
}

func (f *FileMonitor) Delfile(ctx context.Context, fileName string) {
	f.delFile(fileName)
}
