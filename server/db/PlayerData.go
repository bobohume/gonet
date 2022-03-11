package db

import(
	"context"
    "gonet/orm"
    "gonet/server/model"
)

// 自动生成代码

func (this *Player) __SaveSimplePlayerData(data model.SimplePlayerData){
    this.SimplePlayerData = data
	this.SimplePlayerData.Dirty = true
    SERVER.GetLog().Printf("玩家[%d] SaveSimplePlayerData", this.Raft.Id)
}

func (this *Player) __LoadSimplePlayerDataDB(PlayerId int64) error{
    data := &model.SimplePlayerData{PlayerId:PlayerId}
    rows, err := SERVER.GetDB().Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.SimplePlayerData, rs.Row())
    }
	return err
}

func (this *Player) __SaveSimplePlayerDataDB(){
	if this.SimplePlayerData.Dirty{
    	SERVER.GetDB().Exec(orm.SaveSql(this.SimplePlayerData))
		this.SimplePlayerData.Dirty = false
	}
}

func (this *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveSimplePlayerData(data)
	}
}

func (this *Player) __SavePlayerKvData(data model.PlayerKvData){
    this.PlayerKvData = data
	this.PlayerKvData.Dirty = true
    SERVER.GetLog().Printf("玩家[%d] SavePlayerKvData", this.Raft.Id)
}

func (this *Player) __LoadPlayerKvDataDB(PlayerId int64) error{
    data := &model.PlayerKvData{PlayerId:PlayerId}
    rows, err := SERVER.GetDB().Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.PlayerKvData, rs.Row())
    }
	return err
}

func (this *Player) __SavePlayerKvDataDB(){
	if this.PlayerKvData.Dirty{
    	SERVER.GetDB().Exec(orm.SaveSql(this.PlayerKvData))
		this.PlayerKvData.Dirty = false
	}
}

func (this *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SavePlayerKvData(data)
	}
}

func (this *Player) __SaveItemData(data model.ItemData){
    this.ItemData = data
	this.ItemData.Dirty = true
    SERVER.GetLog().Printf("玩家[%d] SaveItemData", this.Raft.Id)
}

func (this *Player) __LoadItemDataDB(PlayerId int64) error{
    data := &model.ItemData{PlayerId:PlayerId}
    rows, err := SERVER.GetDB().Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.ItemData, rs.Row())
    }
	return err
}

func (this *Player) __SaveItemDataDB(){
	if this.ItemData.Dirty{
    	SERVER.GetDB().Exec(orm.SaveSql(this.ItemData))
		this.ItemData.Dirty = false
	}
}

func (this *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveItemData(data)
	}
}

func (this *Player) __SaveEquipData(data model.EquipData){
    this.EquipData = data
	this.EquipData.Dirty = true
    SERVER.GetLog().Printf("玩家[%d] SaveEquipData", this.Raft.Id)
}

func (this *Player) __LoadEquipDataDB(PlayerId int64) error{
    data := &model.EquipData{PlayerId:PlayerId}
    rows, err := SERVER.GetDB().Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.EquipData, rs.Row())
    }
	return err
}

func (this *Player) __SaveEquipDataDB(){
	if this.EquipData.Dirty{
    	SERVER.GetDB().Exec(orm.SaveSql(this.EquipData))
		this.EquipData.Dirty = false
	}
}

func (this *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveEquipData(data)
	}
}

func (this *Player) __SaveMailData(data model.MailData){
    this.MailData = data
	this.MailData.Dirty = true
    SERVER.GetLog().Printf("玩家[%d] SaveMailData", this.Raft.Id)
}

func (this *Player) __LoadMailDataDB(PlayerId int64) error{
    data := &model.MailData{PlayerId:PlayerId}
    rows, err := SERVER.GetDB().Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.MailData, rs.Row())
    }
	return err
}

func (this *Player) __SaveMailDataDB(){
	if this.MailData.Dirty{
    	SERVER.GetDB().Exec(orm.SaveSql(this.MailData))
		this.MailData.Dirty = false
	}
}

func (this *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveMailData(data)
	}
}

func (this *Player) __SaveSocialData(data model.SocialData){
    this.SocialData = data
	this.SocialData.Dirty = true
    SERVER.GetLog().Printf("玩家[%d] SaveSocialData", this.Raft.Id)
}

func (this *Player) __LoadSocialDataDB(PlayerId int64) error{
    data := &model.SocialData{PlayerId:PlayerId}
    rows, err := SERVER.GetDB().Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.SocialData, rs.Row())
    }
	return err
}

func (this *Player) __SaveSocialDataDB(){
	if this.SocialData.Dirty{
    	SERVER.GetDB().Exec(orm.SaveSql(this.SocialData))
		this.SocialData.Dirty = false
	}
}

func (this *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveSocialData(data)
	}
}

func (this *Player) LoadPlayerDB(PlayerId int64) error{
    this.Init(PlayerId)
    if err := this.__LoadSimplePlayerDataDB(PlayerId); err != nil{
        SERVER.GetLog().Printf("__LoadSimplePlayerDataDB() error")
        return err 
    }
    if err := this.__LoadPlayerKvDataDB(PlayerId); err != nil{
        SERVER.GetLog().Printf("__LoadPlayerKvDataDB() error")
        return err 
    }
    if err := this.__LoadItemDataDB(PlayerId); err != nil{
        SERVER.GetLog().Printf("__LoadItemDataDB() error")
        return err 
    }
    if err := this.__LoadEquipDataDB(PlayerId); err != nil{
        SERVER.GetLog().Printf("__LoadEquipDataDB() error")
        return err 
    }
    if err := this.__LoadMailDataDB(PlayerId); err != nil{
        SERVER.GetLog().Printf("__LoadMailDataDB() error")
        return err 
    }
    if err := this.__LoadSocialDataDB(PlayerId); err != nil{
        SERVER.GetLog().Printf("__LoadSocialDataDB() error")
        return err 
    }
    return nil
}


func (this *Player) SavePlayerDB(){
    this.__SaveSimplePlayerDataDB()
    this.__SavePlayerKvDataDB()
    this.__SaveItemDataDB()
    this.__SaveEquipDataDB()
    this.__SaveMailDataDB()
    this.__SaveSocialDataDB()
}

