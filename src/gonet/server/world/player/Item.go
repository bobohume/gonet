package player

import (
	"gonet/server/world/data"
	"gonet/base"
	"math"
	"gonet/db"
	"gonet/server/world"
	"database/sql"
)

type(
	Item struct {
		Id int64 `sql:"primary"`						//物品唯一Id
		PlayerId int `sql:"primary;name:player_id"`		//玩家Id
		ItemId int	`sql:"name:item_id"`				//模板Id
		Quantity int `sql:name:quantity`				//数量(对于装备，只能为1)
	}

	Equip struct {
		Id int64 `sql:"primary"`						//物品唯一Id
		PlayerId int `sql:"primary;name:player_id"`		//玩家Id
		ItemId int `sql:"name:item_id"`					//模板Id
		Level int	`sql:"name:level"`					//等级
		StrengthenLv int `sql:"name:strengthen_lv"`		//强化等级
	}

	ItemEquipPair struct {
		Item *Item
		Equip *Equip
	}

	ItemMgr struct {
		m_ItemMap map[int64] *Item
		m_EquipMap map[int64] *Equip
		m_Player *Player
		m_db *sql.DB
		m_Log *base.CLog
	}

	IItemMgr interface {
		Init(*Player)
		CreateItem(int, int, int) (*Item, *Equip)		//创建物品
		AddItem(int, int, int)	bool					//物品操作
		//SortItem(int) bool							//排序物品
		CanReduceItem(int, int, int) bool				//能否扣除
		addItem(int, int, int)	bool					//添加物品
		reduceItem(int, int, int) bool					//删除物品
		DelEquipById(uint64, int) bool					//删除装备
		DelEquip(*Equip) bool							//删除装备
	}
)

func (this *ItemMgr) Init(pPlayer *Player){
	this.m_Player = pPlayer
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
}

func (this *ItemMgr) CreateItem(ItemId, PlayerId, Quantity int) (*Item, *Equip) {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return nil, nil
	}

	pItem := &Item{}
	pItem.Id = base.UUID.UUID()
	pItem.ItemId = ItemId
	pItem.Quantity = Quantity
	pItem.PlayerId = PlayerId

	var pEquip *Equip
	if pItemData.IsEquip(){
		pEquip = &Equip{}
		pEquip.Id = pItem.Id
		pEquip.ItemId = ItemId
		pEquip.PlayerId = PlayerId
	}
	return pItem, pEquip
}

func (this *ItemMgr) AddItem(ItemId, PlayerId, Quantity int) bool{
	if Quantity > 0 {
		return this.addItem(ItemId, PlayerId, Quantity)
	}

	return  this.reduceItem(ItemId, PlayerId, Quantity)
}

func (this *ItemMgr) CanReduceItem(ItemId, PlayerId, Quantity int) bool{
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return false
	}
	iLeftQuantity := int(math.Abs(float64(Quantity)))
	bEnough := false
	for _,pItem := range this.m_ItemMap{
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

	return !bEnough
}

func (this *ItemMgr) DelEquip(pEquip *Equip) bool{
	if pEquip != nil{
		pItem, exist := this.m_ItemMap[pEquip.Id]
		if exist{
			this.m_db.Exec(db.InsertSql(pItem, "tbl_item"))
		}
		this.m_db.Exec(db.InsertSql(pEquip, "tbl_equip"))
		delete(this.m_ItemMap, pEquip.Id)
		delete(this.m_EquipMap, pEquip.Id)
		return true
	}
	return false
}

func (this *ItemMgr) DelEquipById(Id int64, PlayerId int) bool{
	pEquip, exist := this.m_EquipMap[Id]
	if exist{
		return this.DelEquip(pEquip)
	}
	return false
}

func (this *ItemMgr) addItem(ItemId, PlayerId, Quantity int) bool{
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return false
	}

	iLeftQuantity, iNeedQuantity:= Quantity, 0
	bEnough := false
	BatMap := make(map[int64] int)
	CreateMap := make(map[int64] (*ItemEquipPair))
	for _, pItem := range this.m_ItemMap{
		if pItem != nil && pItem.ItemId == ItemId && pItem.Quantity < pItemData.MaxDie{
			iNeedQuantity = iLeftQuantity
			iLeftQuantity -= pItemData.MaxDie - pItem.Quantity
		}

		if iLeftQuantity > 0 {
			BatMap[pItem.Id] = pItemData.MaxDie - pItem.Quantity
		}else{
			BatMap[pItem.Id] = iNeedQuantity
		}
	}

	for iLeftQuantity > 0{
		iNeedQuantity = iLeftQuantity
		iLeftQuantity -= pItemData.MaxDie

		if iLeftQuantity > 0 {
			pItem, pEquip := this.CreateItem(ItemId, PlayerId, pItemData.MaxDie)
			if pItem != nil && pEquip != nil{
				CreateMap[pItem.Id] = &ItemEquipPair{pItem, pEquip}
			} else {
				bEnough = true
				break
			}
		} else{
			pItem, pEquip := this.CreateItem(ItemId, PlayerId, iNeedQuantity)
			if pItem != nil && pEquip != nil{
				CreateMap[pItem.Id] = &ItemEquipPair{pItem, pEquip}
			} else {
				bEnough = true
			}
			break
		}
	}

	if !bEnough{
		for i, v := range BatMap{
			pItem, exist := this.m_ItemMap[i]
			if exist{
				pItem.Quantity += v
				this.m_db.Exec(db.UpdateSqlEx(pItem, "tbl_item", "quantity"))
			}
		}

		for _, v := range CreateMap{
			if v.Item != nil{
				this.m_ItemMap[v.Item.Id] = v.Item
				this.m_db.Exec(db.InsertSql(v.Item, "tbl_item"))
			}

			if v.Equip != nil{
				this.m_EquipMap[v.Equip.Id] = v.Equip
				this.m_db.Exec(db.InsertSql(v.Equip, "tbl_equip"))
			}
		}
	}

	return !bEnough
}

func (this *ItemMgr) reduceItem(ItemId, PlayerId, Quantity int) bool{
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return false
	}

	iLeftQuantity, iNeedQuantity := int(math.Abs(float64(Quantity))), 0
	bEnough := false
	bEquip := pItemData.IsEquip()
	BatMap := make(map[int64] int)
	for _, pItem := range this.m_ItemMap{
		if pItem != nil && pItem.ItemId == ItemId{
			iNeedQuantity = iLeftQuantity
			iLeftQuantity -= pItem.Quantity

			if iLeftQuantity > 0 {
				BatMap[pItem.Id] = pItem.Quantity
			}else{
				BatMap[pItem.Id] = iNeedQuantity
			}
		}
	}

	if iLeftQuantity > 0{
		bEnough = true
	}

	if !bEnough{
		for i, v := range BatMap{
			pItem, exist := this.m_ItemMap[i]
			if exist{
				pItem.Quantity -= v
				if pItem.Quantity == 0{
					delete(this.m_ItemMap, i)
					this.m_db.Exec(db.DeleteSql(pItem, "tbl_item"))
				}else{
					this.m_db.Exec(db.UpdateSqlEx(pItem, "tbl_item", "quantity"))
				}
			}

			if bEquip{
				pEquip, exist := this.m_EquipMap[i]
				if exist{
					delete(this.m_EquipMap, i)
					this.m_db.Exec(db.DeleteSql(pEquip, "tbl_equip"))
				}
			}
		}
	}

	return !bEnough
}
