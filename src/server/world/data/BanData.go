package data

import (
	"server/common"
	"base"
	"log"
	"strings"
)

type(
	CBanData struct{
		BanName string
	}

	CBanDataRes struct{
		common.CBaseDataRes
	}
)

var(
	BANDREPOSITORY	CBanDataRes
)

func (this *CBanDataRes) Read() bool {
	this.Init()
	var file base.CDataFile
	lineData := &base.RData{}

	if (!file.ReadDataFile("data/BanWord.dat")) {
		log.Fatalf("read BanWord.dat error")
		return false
	}

	for i := 0; i < file.RecordNum; i++{
		pData := CBanData{}

		file.GetData(lineData)
		base.IFAssert(lineData.Type == base.DType_String, "read BanWord.dat BanName error")
		pData.BanName = lineData.String
		if pData.BanName == ""{
			continue
		}

		this.AddData(pData.BanName, pData)
	}

	return true
}

func ReplaceBanWord(msg string, replace string) string{
	for i,_ := range BANDREPOSITORY.DataMap{
		msg = strings.Replace(msg, i.(string), replace, -1 )
	}
	return msg
}