package data

import (
	"gonet/base"
	"gonet/server/cm"
	"log"
	"strings"
)

type (
	BanData struct {
		BanName string
	}

	BanDataRes struct {
		cm.BaseDataRes
	}
)

var (
	BANDATA BanDataRes
)

func (b *BanDataRes) Read() bool {
	b.Init()
	var file base.DataFile
	lineData := &base.Data{}

	if !file.ReadDataFile("data/BanWord.dat") {
		log.Fatalf("read BanWord.dat error")
		return false
	}

	for i := 0; i < file.RecordNum; i++ {
		pData := BanData{}

		file.GetData(lineData)
		pData.BanName = lineData.String("BanName")
		if pData.BanName == "" {
			continue
		}

		//b.AddData(pData.BanName, pData)
	}

	return true
}

func (b *BanDataRes) GetData(id int) *BanData {
	pData := b.BaseDataRes.GetBaseData(id)
	if pData != nil {
		return pData.(*BanData)
	}

	return nil
}

func ReplaceBanWord(msg string, replace string) string {
	for i, _ := range BANDATA.DataMap {
		msg = strings.Replace(msg, i.(string), replace, -1)
	}
	return msg
}
