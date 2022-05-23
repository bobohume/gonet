package Excel__test

import (
	"gonet/base"
	"gonet/common"
	"log"
	"testing"
)

const (
	BUFF_DATA = "buff.dat"
)

type (
	BuffData struct {
		Id       int
		Series   int
		Lv       int
		Effect   int64
		Flag     int
		Target   int
		BuffIds  []int
		BuffNums []int
	}

	BuffDataRes struct {
		common.BaseDataRes
	}

	IBuffDataRes interface {
		common.IBaseDataRes
	}
)

var (
	BUFFDATA IBuffDataRes
)

func (this *BuffDataRes) Read() bool {
	this.Init()
	var file base.DataFile
	lineData := &base.Data{}

	if !file.ReadDataFile(BUFF_DATA) {
		log.Fatalf("read buff.dat error")
		return false
	}

	for i := 0; i < file.RecordNum; i++ {
		pData := BuffData{}
		file.GetData(lineData)
		pData.Id = lineData.Int("id")

		file.GetData(lineData)
		pData.Series = lineData.Int("Series")

		file.GetData(lineData)
		pData.Lv = lineData.Int("Lv")

		file.GetData(lineData)
		file.GetData(lineData)
		file.GetData(lineData)

		file.GetData(lineData)
		pData.Effect = lineData.Int64("Effect")

		file.GetData(lineData)
		pData.Flag = lineData.Int("Flag")

		file.GetData(lineData)
		pData.Target = lineData.Enum("Target")

		file.GetData(lineData)
		pData.BuffIds = lineData.IntArray("BuffIds")

		file.GetData(lineData)
		pData.BuffNums = lineData.IntArray("BuffNums")

		this.AddData(pData.Id, pData)
	}

	return true
}

func (this *BuffDataRes) GetData(id int) *BuffData {
	pData := this.BaseDataRes.GetBaseData(id)
	if pData != nil {
		return pData.(*BuffData)
	}

	return nil
}

func TestBuffData(t *testing.T) {
	data := BuffDataRes{}
	data.Read()
}
