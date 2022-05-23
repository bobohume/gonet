package data

import (
	"gonet/base"
	"gonet/common"
	"log"
)

type (
	ItemData struct {
		ItemId int
		Type   int //类型
		MaxDie int //最大叠加数
	}

	ItemDataRes struct {
		common.BaseDataRes
	}

	IItemDataRes interface {
		common.IBaseDataRes
	}
)

var (
	ITEMDATA ItemDataRes
)

func (it *ItemData) IsEquip() bool {
	return false
}

func (it *ItemDataRes) Read() bool {
	it.Init()
	var file base.DataFile
	//lineData := &base.RData{}

	if !file.ReadDataFile("data/BanWord.dat") {
		log.Fatalf("read BanWord.dat error")
		return false
	}

	for i := 0; i < file.RecordNum; i++ {
		pData := ItemData{}
		it.AddData(pData.ItemId, pData)
	}

	return true
}

func (it *ItemDataRes) GetData(id int) *ItemData {
	pData := it.BaseDataRes.GetBaseData(id)
	if pData != nil {
		return pData.(*ItemData)
	}

	return nil
}
