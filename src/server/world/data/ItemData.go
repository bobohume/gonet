package data

import (
	"gonet/server/common"
	"gonet/base"
	"log"
)

type(
	ItemData struct{
		ItemId int
		Type int	//类型
		MaxDie int	//最大叠加数
	}

	ItemDataRes struct{
		common.BaseDataRes
	}

	IItemDataRes interface {
		common.IBaseDataRes
	}
)

var(
	ITEMDATA	ItemDataRes
)

func (this *ItemData) IsEquip() bool{
	return false
}

func (this *ItemDataRes) Read() bool{
	this.Init()
	var file base.CDataFile
	//lineData := &base.RData{}

	if (!file.ReadDataFile("data/BanWord.dat")) {
		log.Fatalf("read BanWord.dat error")
		return false
	}

	for i := 0; i < file.RecordNum; i++{
		pData := ItemData{}
		this.AddData(pData.ItemId, pData)
	}

	return true
}

func (this *ItemDataRes) GetData(id int) *ItemData {
	pData := this.BaseDataRes.GetBaseData(id)
	if pData != nil{
		return pData.(*ItemData)
	}

	return nil
}