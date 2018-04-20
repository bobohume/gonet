package main

import (
	"server/common"
	"base"
	"log"
	"net/http"
	"fmt"
	"time"
)

type(
	CApplyData struct{
		Id int
		Des string
		Type [2]int
	}

	CApplyDataRes struct{
		common.CBaseDataRes
	}
)

func (this *CApplyDataRes) Read() bool {
	this.Init()
	var file base.CDataFile
	lineData := &base.RData{}

	if (!file.ReadDataFile("data/Apply.dat")) {
		log.Fatalf("read apply.date error")
		return false
	}

	for i := 0; i < file.RecordNum; i++{
		pData := CApplyData{}

		file.GetData(lineData)
		base.IFAssert(lineData.Type == base.DType_U16, "read applydata ID error")
		pData.Id = int(lineData.U16)

		file.GetData(lineData)
		base.IFAssert(lineData.Type == base.DType_String, "read applydata 奏章描述 error")
		pData.Des = lineData.String

		for i := 0; i < 2; i++{
			file.GetData(lineData)
			base.IFAssert(lineData.Type == base.DType_U8, "read applydata 执行效果 error")
			pData.Type[i] = int(lineData.U8)
		}
		this.AddData(pData.Id, pData)
	}

	return true
}

func main() {
	//这里实现了远程获取pprof数据的接口
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	file := &CApplyDataRes{}
	file.Read()
	fmt.Println(file)

	for {
		time.Sleep(1000)
	}
}