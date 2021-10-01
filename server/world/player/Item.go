package player

import (
	"gonet/base"
	"gonet/db"
	"gonet/server/model"
	"gonet/server/world/data"
	"math"
)

type(
	ItemEquipPair struct {
		Item *model.Item
		Equip *model.Equip
	}

	ItemMgr struct {
		IPlayer
		m_ItemMap map[int64] *model.Item
		m_EquipMap map[int64] *model.Equip
	}

	IItemMgr interface {
		Init(IPlayer)
		CreateItem(int, int) (*model.Item, *model.Equip)		//创建物品
		AddItem(int, int)	bool					//物品操作
		//SortItem(int) bool						//排序物品
		CanReduceItem(int, int) bool				//能否扣除
		//addItem(int, int)	bool					//添加物品
		//reduceItem(int, int) bool					//删除物品
		DelEquipById(int64) bool					//删除装备
		DelEquip(*model.Equip) bool						//删除装备
	}
)

func (this *Player) GetItemMgr() IItemMgr{
	return this.m_ItemMgr
}

func (this *ItemMgr) Init(pPlayer IPlayer){
	this.IPlayer = pPlayer
	//test
	/*this.RegisterCall("C_W_AddEquipAttrRequest", func(ctx context.Context, packet *rpc.C_W_ChatMessage) {
		world.SendToClient(this.GetGateClusterId(), &rpc.W_C_ChatMessage{
			PacketHead:rpc.BuildPacketHead(this.GetAccountId(), 0 ),
		})
	})*/
}

func (this *ItemMgr) CreateItem(ItemId int, Quantity int) (*model.Item, *model.Equip) {
	pItemData := data.ITEMDATA.GetData(ItemId)
	if pItemData == nil{
		return nil, nil
	}

	pItem := &model.Item{}
	pItem.Id = base.UUID.UUID()
	pItem.ItemId = ItemId
	pItem.Quantity = Quantity
	pItem.PlayerId = this.GetPlayerId()

	var pEquip *model.Equip
	if pItemData.IsEquip(){
		pEquip = &model.Equip{}
		pEquip.Id = pItem.Id
		pEquip.ItemId = ItemId
		pEquip.PlayerId = this.GetPlayerId()
	}
	return pItem, pEquip
}

func (this *ItemMgr) AddItem(ItemId int, Quantity int) bool{
	if Quantity > 0 {
		return this.addItem(ItemId, Quantity)
	}

	return  this.reduceItem(ItemId, Quantity)
}

func (this *ItemMgr) CanReduceItem(ItemId int, Quantity int) bool{
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

func (this *ItemMgr) DelEquip(pEquip *model.Equip) bool{
	if pEquip != nil{
		pItem, exist := this.m_ItemMap[pEquip.Id]
		if exist{
			this.GetDB().Exec(db.DeleteSql(pItem,))
		}
		this.GetDB().Exec(db.DeleteSql(pEquip))
		delete(this.m_ItemMap, pEquip.Id)
		delete(this.m_EquipMap, pEquip.Id)
		return true
	}
	return false
}

func (this *ItemMgr) DelEquipById(Id int64) bool{
	pEquip, exist := this.m_EquipMap[Id]
	if exist{
		return this.DelEquip(pEquip)
	}
	return false
}

func (this *ItemMgr) addItem(ItemId int, Quantity int) bool{
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
			pItem, exist := this.m_ItemMap[i]
			if exist{
				pItem.Quantity += v
				this.GetDB().Exec(db.UpdateSql(pItem))
			}
		}

		for _, v := range CreateMap{
			if v.Item != nil{
				this.m_ItemMap[v.Item.Id] = v.Item
				this.GetDB().Exec(db.InsertSql(v.Item))
			}

			if v.Equip != nil{
				this.m_EquipMap[v.Equip.Id] = v.Equip
				this.GetDB().Exec(db.InsertSql(v.Equip))
			}
		}
	}

	return !bEnough
}

func (this *ItemMgr) reduceItem(ItemId int, Quantity int) bool{
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
				break
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
					this.GetDB().Exec(db.DeleteSql(pItem))
				}else{
					this.GetDB().Exec(db.UpdateSql(pItem))
				}
			}

			if bEquip{
				pEquip, exist := this.m_EquipMap[i]
				if exist{
					delete(this.m_EquipMap, i)
					this.GetDB().Exec(db.DeleteSql(pEquip))
				}
			}
		}
	}

	return !bEnough
}
