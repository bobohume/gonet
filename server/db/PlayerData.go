package db

import (
	"context"
	"gonet/base"
	"gonet/orm"
	"gonet/server/model"
)

// 自动生成代码

func (this *Player) __LoadSimplePlayerDataDB(PlayerId int64) error {
	data := &model.SimplePlayerData{PlayerId: PlayerId}
	rs, err := orm.LoadSql(data, orm.WithWhere(data))
	if err == nil && rs.Next() {
		orm.LoadObjSql(&this.SimplePlayerData, rs.Row())
	}
	return err
}

/*
	func (this *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData){
		orm.SaveSql(&data)
		base.LOG.Printf("玩家[%d] SaveSimplePlayerData", playerId)
	}
*/
func (this *Player) __SaveSimplePlayerDataDB() {
	if this.SimplePlayerData.Dirty {
		orm.SaveSql(this.SimplePlayerData)
		this.SimplePlayerData.Dirty = false
	}
}

func (this *Player) __SaveSimplePlayerData(data model.SimplePlayerData) {
	this.SimplePlayerData = data
	this.SimplePlayerData.Dirty = true
	base.LOG.Printf("玩家[%d] SaveSimplePlayerData", this.SimplePlayerData.PlayerId)
}

func (this *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData) {
	player := this.GetPlayer(playerId)
	if player != nil {
		player.__SaveSimplePlayerData(data)
	}
}

func (this *Player) __LoadPlayerKvDataDB(PlayerId int64) error {
	data := &model.PlayerKvData{PlayerId: PlayerId}
	rs, err := orm.LoadSql(data, orm.WithWhere(data))
	if err == nil && rs.Next() {
		orm.LoadObjSql(&this.PlayerKvData, rs.Row())
	}
	return err
}

/*
	func (this *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData){
		orm.SaveSql(&data)
		base.LOG.Printf("玩家[%d] SavePlayerKvData", playerId)
	}
*/
func (this *Player) __SavePlayerKvDataDB() {
	if this.PlayerKvData.Dirty {
		orm.SaveSql(this.PlayerKvData)
		this.PlayerKvData.Dirty = false
	}
}

func (this *Player) __SavePlayerKvData(data model.PlayerKvData) {
	this.PlayerKvData = data
	this.PlayerKvData.Dirty = true
	base.LOG.Printf("玩家[%d] SavePlayerKvData", this.PlayerKvData.PlayerId)
}

func (this *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData) {
	player := this.GetPlayer(playerId)
	if player != nil {
		player.__SavePlayerKvData(data)
	}
}

func (this *Player) __LoadItemDataDB(PlayerId int64) error {
	data := &model.ItemData{PlayerId: PlayerId}
	rs, err := orm.LoadSql(data, orm.WithWhere(data))
	if err == nil && rs.Next() {
		orm.LoadObjSql(&this.ItemData, rs.Row())
	}
	return err
}

/*
	func (this *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData){
		orm.SaveSql(&data)
		base.LOG.Printf("玩家[%d] SaveItemData", playerId)
	}
*/
func (this *Player) __SaveItemDataDB() {
	if this.ItemData.Dirty {
		orm.SaveSql(this.ItemData)
		this.ItemData.Dirty = false
	}
}

func (this *Player) __SaveItemData(data model.ItemData) {
	this.ItemData = data
	this.ItemData.Dirty = true
	base.LOG.Printf("玩家[%d] SaveItemData", this.ItemData.PlayerId)
}

func (this *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData) {
	player := this.GetPlayer(playerId)
	if player != nil {
		player.__SaveItemData(data)
	}
}

func (this *Player) __LoadEquipDataDB(PlayerId int64) error {
	data := &model.EquipData{PlayerId: PlayerId}
	rs, err := orm.LoadSql(data, orm.WithWhere(data))
	if err == nil && rs.Next() {
		orm.LoadObjSql(&this.EquipData, rs.Row())
	}
	return err
}

/*
	func (this *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData){
		orm.SaveSql(&data)
		base.LOG.Printf("玩家[%d] SaveEquipData", playerId)
	}
*/
func (this *Player) __SaveEquipDataDB() {
	if this.EquipData.Dirty {
		orm.SaveSql(this.EquipData)
		this.EquipData.Dirty = false
	}
}

func (this *Player) __SaveEquipData(data model.EquipData) {
	this.EquipData = data
	this.EquipData.Dirty = true
	base.LOG.Printf("玩家[%d] SaveEquipData", this.EquipData.PlayerId)
}

func (this *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData) {
	player := this.GetPlayer(playerId)
	if player != nil {
		player.__SaveEquipData(data)
	}
}

func (this *Player) __LoadMailDataDB(PlayerId int64) error {
	data := &model.MailData{PlayerId: PlayerId}
	rs, err := orm.LoadSql(data, orm.WithWhere(data))
	if err == nil && rs.Next() {
		orm.LoadObjSql(&this.MailData, rs.Row())
	}
	return err
}

/*
	func (this *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData){
		orm.SaveSql(&data)
		base.LOG.Printf("玩家[%d] SaveMailData", playerId)
	}
*/
func (this *Player) __SaveMailDataDB() {
	if this.MailData.Dirty {
		orm.SaveSql(this.MailData)
		this.MailData.Dirty = false
	}
}

func (this *Player) __SaveMailData(data model.MailData) {
	this.MailData = data
	this.MailData.Dirty = true
	base.LOG.Printf("玩家[%d] SaveMailData", this.MailData.PlayerId)
}

func (this *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData) {
	player := this.GetPlayer(playerId)
	if player != nil {
		player.__SaveMailData(data)
	}
}

func (this *Player) __LoadSocialDataDB(PlayerId int64) error {
	data := &model.SocialData{PlayerId: PlayerId}
	rs, err := orm.LoadSql(data, orm.WithWhere(data))
	if err == nil && rs.Next() {
		orm.LoadObjSql(&this.SocialData, rs.Row())
	}
	return err
}

/*
	func (this *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData){
		orm.SaveSql(&data)
		base.LOG.Printf("玩家[%d] SaveSocialData", playerId)
	}
*/
func (this *Player) __SaveSocialDataDB() {
	if this.SocialData.Dirty {
		orm.SaveSql(this.SocialData)
		this.SocialData.Dirty = false
	}
}

func (this *Player) __SaveSocialData(data model.SocialData) {
	this.SocialData = data
	this.SocialData.Dirty = true
	base.LOG.Printf("玩家[%d] SaveSocialData", this.SocialData.PlayerId)
}

func (this *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData) {
	player := this.GetPlayer(playerId)
	if player != nil {
		player.__SaveSocialData(data)
	}
}

func (this *Player) LoadPlayerDB(PlayerId int64) error {
	this.Init(PlayerId)
	if err := this.__LoadSimplePlayerDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadSimplePlayerDataDB() error")
		return err
	}
	if err := this.__LoadPlayerKvDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadPlayerKvDataDB() error")
		return err
	}
	if err := this.__LoadItemDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadItemDataDB() error")
		return err
	}
	if err := this.__LoadEquipDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadEquipDataDB() error")
		return err
	}
	if err := this.__LoadMailDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadMailDataDB() error")
		return err
	}
	if err := this.__LoadSocialDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadSocialDataDB() error")
		return err
	}
	return nil
}

func (this *Player) SavePlayerDB() {
	this.__SaveSimplePlayerDataDB()
	this.__SavePlayerKvDataDB()
	this.__SaveItemDataDB()
	this.__SaveEquipDataDB()
	this.__SaveMailDataDB()
	this.__SaveSocialDataDB()
}
