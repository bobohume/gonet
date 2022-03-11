package player

import (
	"gonet/base"
	"gonet/server/model"
	"gonet/server/game/data"
	"math"
)

type(
	ItemEquipPair struct {
		Item *model.Item
		Equip *model.Equip
	}
)

//创建物品
func (this *Player) CreateItem(ItemId int, Quantity int) (*model.Item, *model.Equip) {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return nil, nil
	}

	pItem := &model.Item{}
	pItem.Id = base.UUID.UUID()
	pItem.ItemId = ItemId
	pItem.Quantity = Quantity
	pItem.PlayerId = this.PlayerId

	var pEquip *model.Equip
	if pItemData.IsEquip(){
		pEquip = &model.Equip{}
		pEquip.Id = pItem.Id
		pEquip.ItemId = ItemId
		pEquip.PlayerId = this.GetPlayerId()
	}
	return pItem, pEquip
}

//物品操作
func (this *Player) AddItem(ItemId int, Quantity int) bool{
	if Quantity > 0 {
		return this.addItem(ItemId, Quantity)
	}

	return  this.reduceItem(ItemId, Quantity)
}

//能否扣除
func (this *Player) CanReduceItem(ItemId int, Quantity int) bool{
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return false
	}
	iLeftQuantity := int(math.Abs(float64(Quantity)))
	bEnough := false
	for _,pItem := range this.Bag.DataMap{
		if pItem != nil && pItem.ItemId == ItemId {
			iLeftQuantity -= pItem.Quantity

			if iLeftQuantity > 0 {
			} else {
				break
			}
		}
	}

	if iLeftQuantity > 0{
		bEnough = true
	}

	return bEnough
}

//删除装备
func (this *Player) DelEquip(Id int64) bool{
	_, exist := this.EquipData.DataMap[Id]
	if exist{
		delete(this.Bag.DataMap, Id)
		delete(this.Bag.DataMap, Id)
		this.SaveItemData()
		this.SaveEquipData()
		return true
	}
	return false
}

func (this *Player) addItem(ItemId int, Quantity int) bool{
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return false
	}

	iLeftQuantity, iNeedQuantity:= Quantity, 0
	bEnough := false
	BatMap := make(map[int64] int)
	CreateMap := make(map[int64] (*ItemEquipPair))
	for _, pItem := range this.Bag.DataMap{
		if pItem != nil && pItem.ItemId == ItemId && pItem.Quantity < pItemData.MaxDie{
			iNeedQuantity = iLeftQuantity
			iLeftQuantity -= pItemData.MaxDie - pItem.Quantity

			if iLeftQuantity > 0 {
				BatMap[pItem.Id] = pItemData.MaxDie - pItem.Quantity
			}else{
				BatMap[pItem.Id] = iNeedQuantity
				break
			}
		}
	}

	for iLeftQuantity > 0{
		iNeedQuantity = iLeftQuantity
		iLeftQuantity -= pItemData.MaxDie

		if iLeftQuantity > 0 {
			pItem, pEquip := this.CreateItem(ItemId, pItemData.MaxDie)
			if pItem != nil{
				CreateMap[pItem.Id] = &ItemEquipPair{pItem, pEquip}
			} else {
				bEnough = true
				break
			}
		} else{
			pItem, pEquip := this.CreateItem(ItemId, iNeedQuantity)
			if pItem != nil{
				CreateMap[pItem.Id] = &ItemEquipPair{pItem, pEquip}
			} else {
				bEnough = true
			}
			break
		}
	}

	if !bEnough{
		for i, v := range BatMap{
			pItem, exist := this.Bag.DataMap[i]
			if exist{
				pItem.Quantity += v
			}
		}

		for _, v := range CreateMap{
			if v.Item != nil{
				this.Bag.DataMap[v.Item.Id] = v.Item
			}

			if v.Equip != nil{
				this.EquipData.DataMap[v.Equip.Id] = v.Equip
			}
		}
	}

	this.SaveItemData()
	this.SaveEquipData()
	return bEnough
}

func (this *Player) reduceItem(ItemId int, Quantity int) bool{
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return false
	}

	iLeftQuantity, iNeedQuantity := int(math.Abs(float64(Quantity))), 0
	bEnough := false
	bEquip := pItemData.IsEquip()
	BatMap := make(map[int64] int)
	for _, pItem := range this.Bag.DataMap{
		if pItem != nil && pItem.ItemId == ItemId{
			iNeedQuantity = iLeftQuantity
			iLeftQuantity -= pItem.Quantity

			if iLeftQuantity > 0 {
				BatMap[pItem.Id] = pItem.Quantity
			}else{
				BatMap[pItem.Id] = iNeedQuantity
				break
			}
		}
	}

	if iLeftQuantity > 0{
		bEnough = true
	}

	if !bEnough{
		for i, v := range BatMap{
			pItem, exist := this.Bag.DataMap[i]
			if exist{
				pItem.Quantity -= v
				if pItem.Quantity == 0{
					delete(this.Bag.DataMap, i)
				}
			}

			if bEquip{
				_, exist := this.EquipData.DataMap[i]
				if exist{
					delete(this.EquipData.DataMap, i)
				}
			}
		}
	}

	this.SaveItemData()
	this.SaveEquipData()
	return bEnough
}
