package player

import (
	"gonet/base"
	"gonet/server/game/data"
	"gonet/server/model"
	"math"
)

type (
	ItemEquipPair struct {
		Item  *model.Item
		Equip *model.Equip
	}
)

//创建物品
func (p *Player) CreateItem(ItemId int, Quantity int) (*model.Item, *model.Equip) {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil {
		return nil, nil
	}

	pItem := &model.Item{}
	pItem.Id = base.UUID.UUID()
	pItem.ItemId = ItemId
	pItem.Quantity = Quantity
	pItem.PlayerId = p.PlayerId

	var pEquip *model.Equip
	if pItemData.IsEquip() {
		pEquip = &model.Equip{}
		pEquip.Id = pItem.Id
		pEquip.ItemId = ItemId
		pEquip.PlayerId = p.GetPlayerId()
	}
	return pItem, pEquip
}

//物品操作
func (p *Player) AddItem(ItemId int, Quantity int) bool {
	if Quantity > 0 {
		return p.addItem(ItemId, Quantity)
	}

	return p.reduceItem(ItemId, Quantity)
}

//能否扣除
func (p *Player) CanReduceItem(ItemId int, Quantity int) bool {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil {
		return false
	}
	iLeftQuantity := int(math.Abs(float64(Quantity)))
	bEnough := false
	for _, pItem := range p.Bag.DataMap {
		if pItem != nil && pItem.ItemId == ItemId {
			iLeftQuantity -= pItem.Quantity

			if iLeftQuantity > 0 {
			} else {
				break
			}
		}
	}

	if iLeftQuantity > 0 {
		bEnough = true
	}

	return bEnough
}

//删除装备
func (p *Player) DelEquip(Id int64) bool {
	_, exist := p.EquipData.DataMap[Id]
	if exist {
		delete(p.Bag.DataMap, Id)
		delete(p.Bag.DataMap, Id)
		p.SaveItemData()
		p.SaveEquipData()
		return true
	}
	return false
}

func (p *Player) addItem(ItemId int, Quantity int) bool {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil {
		return false
	}

	iLeftQuantity, iNeedQuantity := Quantity, 0
	bEnough := false
	BatMap := make(map[int64]int)
	CreateMap := make(map[int64](*ItemEquipPair))
	for _, pItem := range p.Bag.DataMap {
		if pItem != nil && pItem.ItemId == ItemId && pItem.Quantity < pItemData.MaxDie {
			iNeedQuantity = iLeftQuantity
			iLeftQuantity -= pItemData.MaxDie - pItem.Quantity

			if iLeftQuantity > 0 {
				BatMap[pItem.Id] = pItemData.MaxDie - pItem.Quantity
			} else {
				BatMap[pItem.Id] = iNeedQuantity
				break
			}
		}
	}

	for iLeftQuantity > 0 {
		iNeedQuantity = iLeftQuantity
		iLeftQuantity -= pItemData.MaxDie

		if iLeftQuantity > 0 {
			pItem, pEquip := p.CreateItem(ItemId, pItemData.MaxDie)
			if pItem != nil {
				CreateMap[pItem.Id] = &ItemEquipPair{pItem, pEquip}
			} else {
				bEnough = true
				break
			}
		} else {
			pItem, pEquip := p.CreateItem(ItemId, iNeedQuantity)
			if pItem != nil {
				CreateMap[pItem.Id] = &ItemEquipPair{pItem, pEquip}
			} else {
				bEnough = true
			}
			break
		}
	}

	if !bEnough {
		for i, v := range BatMap {
			pItem, exist := p.Bag.DataMap[i]
			if exist {
				pItem.Quantity += v
			}
		}

		for _, v := range CreateMap {
			if v.Item != nil {
				p.Bag.DataMap[v.Item.Id] = v.Item
			}

			if v.Equip != nil {
				p.EquipData.DataMap[v.Equip.Id] = v.Equip
			}
		}
	}

	p.SaveItemData()
	p.SaveEquipData()
	return bEnough
}

func (p *Player) reduceItem(ItemId int, Quantity int) bool {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil {
		return false
	}

	iLeftQuantity, iNeedQuantity := int(math.Abs(float64(Quantity))), 0
	bEnough := false
	bEquip := pItemData.IsEquip()
	BatMap := make(map[int64]int)
	for _, pItem := range p.Bag.DataMap {
		if pItem != nil && pItem.ItemId == ItemId {
			iNeedQuantity = iLeftQuantity
			iLeftQuantity -= pItem.Quantity

			if iLeftQuantity > 0 {
				BatMap[pItem.Id] = pItem.Quantity
			} else {
				BatMap[pItem.Id] = iNeedQuantity
				break
			}
		}
	}

	if iLeftQuantity > 0 {
		bEnough = true
	}

	if !bEnough {
		for i, v := range BatMap {
			pItem, exist := p.Bag.DataMap[i]
			if exist {
				pItem.Quantity -= v
				if pItem.Quantity == 0 {
					delete(p.Bag.DataMap, i)
				}
			}

			if bEquip {
				_, exist := p.EquipData.DataMap[i]
				if exist {
					delete(p.EquipData.DataMap, i)
				}
			}
		}
	}

	p.SaveItemData()
	p.SaveEquipData()
	return bEnough
}
