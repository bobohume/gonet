package db

import(
	"context"
	"gonet/base"
    "gonet/orm"
    "gonet/server/model"
)

// 自动生成代码

func (this *Player) __LoadSimplePlayerDataDB(PlayerId int64) error{
    data := &model.SimplePlayerData{PlayerId:PlayerId}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.SimplePlayerData, rs.Row())
    }
	return err
}

func (this *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveSimplePlayerData", playerId)
}
/*
func (this *Player) __SaveSimplePlayerDataDB(){
	if this.SimplePlayerData.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.SimplePlayerData))
		this.SimplePlayerData.Dirty = false
	}
}

func (this *Player) __SaveSimplePlayerData(data model.SimplePlayerData){
    this.SimplePlayerData = data
	this.SimplePlayerData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveSimplePlayerData", this.MailBox.Id)
}

func (this *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveSimplePlayerData(data)
	}
}
*/

func (this *Player) __LoadPlayerKvDataDB(PlayerId int64) error{
    data := &model.PlayerKvData{PlayerId:PlayerId}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.PlayerKvData, rs.Row())
    }
	return err
}

func (this *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SavePlayerKvData", playerId)
}
/*
func (this *Player) __SavePlayerKvDataDB(){
	if this.PlayerKvData.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.PlayerKvData))
		this.PlayerKvData.Dirty = false
	}
}

func (this *Player) __SavePlayerKvData(data model.PlayerKvData){
    this.PlayerKvData = data
	this.PlayerKvData.Dirty = true
    base.LOG.Printf("玩家[%d] SavePlayerKvData", this.MailBox.Id)
}

func (this *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SavePlayerKvData(data)
	}
}
*/

func (this *Player) __LoadItemDataDB(PlayerId int64) error{
    data := &model.ItemData{PlayerId:PlayerId}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.ItemData, rs.Row())
    }
	return err
}

func (this *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveItemData", playerId)
}
/*
func (this *Player) __SaveItemDataDB(){
	if this.ItemData.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.ItemData))
		this.ItemData.Dirty = false
	}
}

func (this *Player) __SaveItemData(data model.ItemData){
    this.ItemData = data
	this.ItemData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveItemData", this.MailBox.Id)
}

func (this *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveItemData(data)
	}
}
*/

func (this *Player) __LoadEquipDataDB(PlayerId int64) error{
    data := &model.EquipData{PlayerId:PlayerId}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.EquipData, rs.Row())
    }
	return err
}

func (this *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveEquipData", playerId)
}
/*
func (this *Player) __SaveEquipDataDB(){
	if this.EquipData.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.EquipData))
		this.EquipData.Dirty = false
	}
}

func (this *Player) __SaveEquipData(data model.EquipData){
    this.EquipData = data
	this.EquipData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveEquipData", this.MailBox.Id)
}

func (this *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveEquipData(data)
	}
}
*/

func (this *Player) __LoadMailDataDB(PlayerId int64) error{
    data := &model.MailData{PlayerId:PlayerId}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.MailData, rs.Row())
    }
	return err
}

func (this *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveMailData", playerId)
}
/*
func (this *Player) __SaveMailDataDB(){
	if this.MailData.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.MailData))
		this.MailData.Dirty = false
	}
}

func (this *Player) __SaveMailData(data model.MailData){
    this.MailData = data
	this.MailData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveMailData", this.MailBox.Id)
}

func (this *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveMailData(data)
	}
}
*/

func (this *Player) __LoadSocialDataDB(PlayerId int64) error{
    data := &model.SocialData{PlayerId:PlayerId}
    rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
    rs, err := orm.Query(rows, err)
    if err == nil && rs.Next() {
        orm.LoadObjSql(&this.SocialData, rs.Row())
    }
	return err
}

func (this *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData){
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveSocialData", playerId)
}
/*
func (this *Player) __SaveSocialDataDB(){
	if this.SocialData.Dirty{
    	orm.DB.Exec(orm.SaveSql(this.SocialData))
		this.SocialData.Dirty = false
	}
}

func (this *Player) __SaveSocialData(data model.SocialData){
    this.SocialData = data
	this.SocialData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveSocialData", this.MailBox.Id)
}

func (this *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData){
	pPlayer, bEx := this.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveSocialData(data)
	}
}
*/

func (this *Player) LoadPlayerDB(PlayerId int64) error{
    this.Init(PlayerId)
    if err := this.__LoadSimplePlayerDataDB(PlayerId); err != nil{
        base.LOG.Printf("__LoadSimplePlayerDataDB() error")
        return err 
    }
    if err := this.__LoadPlayerKvDataDB(PlayerId); err != nil{
        base.LOG.Printf("__LoadPlayerKvDataDB() error")
        return err 
    }
    if err := this.__LoadItemDataDB(PlayerId); err != nil{
        base.LOG.Printf("__LoadItemDataDB() error")
        return err 
    }
    if err := this.__LoadEquipDataDB(PlayerId); err != nil{
        base.LOG.Printf("__LoadEquipDataDB() error")
        return err 
    }
    if err := this.__LoadMailDataDB(PlayerId); err != nil{
        base.LOG.Printf("__LoadMailDataDB() error")
        return err 
    }
    if err := this.__LoadSocialDataDB(PlayerId); err != nil{
        base.LOG.Printf("__LoadSocialDataDB() error")
        return err 
    }
    return nil
}

