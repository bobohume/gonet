package db

import (
	"context"
	"gonet/base"
	"gonet/orm"
	"gonet/server/model"
)

// 自动生成代码

func (p *Player) __LoadSimplePlayerDataDB(PlayerId int64) error {
	data := &model.SimplePlayerData{PlayerId: PlayerId}
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		orm.LoadObjSql(&p.SimplePlayerData, rs.Row())
	}
	return err
}

func (p *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData) {
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveSimplePlayerData", playerId)
}

/*
func (p *Player) __SaveSimplePlayerDataDB(){
	if p.SimplePlayerData.Dirty{
    	orm.DB.Exec(orm.SaveSql(p.SimplePlayerData))
		p.SimplePlayerData.Dirty = false
	}
}

func (p *Player) __SaveSimplePlayerData(data model.SimplePlayerData){
    p.SimplePlayerData = data
	p.SimplePlayerData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveSimplePlayerData", p.MailBox.Id)
}

func (p *PlayerMgr) SaveSimplePlayerData(ctx context.Context, playerId int64, data model.SimplePlayerData){
	pPlayer, bEx := p.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveSimplePlayerData(data)
	}
}
*/

func (p *Player) __LoadPlayerKvDataDB(PlayerId int64) error {
	data := &model.PlayerKvData{PlayerId: PlayerId}
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		orm.LoadObjSql(&p.PlayerKvData, rs.Row())
	}
	return err
}

func (p *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData) {
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SavePlayerKvData", playerId)
}

/*
func (p *Player) __SavePlayerKvDataDB(){
	if p.PlayerKvData.Dirty{
    	orm.DB.Exec(orm.SaveSql(p.PlayerKvData))
		p.PlayerKvData.Dirty = false
	}
}

func (p *Player) __SavePlayerKvData(data model.PlayerKvData){
    p.PlayerKvData = data
	p.PlayerKvData.Dirty = true
    base.LOG.Printf("玩家[%d] SavePlayerKvData", p.MailBox.Id)
}

func (p *PlayerMgr) SavePlayerKvData(ctx context.Context, playerId int64, data model.PlayerKvData){
	pPlayer, bEx := p.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SavePlayerKvData(data)
	}
}
*/

func (p *Player) __LoadItemDataDB(PlayerId int64) error {
	data := &model.ItemData{PlayerId: PlayerId}
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		orm.LoadObjSql(&p.ItemData, rs.Row())
	}
	return err
}

func (p *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData) {
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveItemData", playerId)
}

/*
func (p *Player) __SaveItemDataDB(){
	if p.ItemData.Dirty{
    	orm.DB.Exec(orm.SaveSql(p.ItemData))
		p.ItemData.Dirty = false
	}
}

func (p *Player) __SaveItemData(data model.ItemData){
    p.ItemData = data
	p.ItemData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveItemData", p.MailBox.Id)
}

func (p *PlayerMgr) SaveItemData(ctx context.Context, playerId int64, data model.ItemData){
	pPlayer, bEx := p.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveItemData(data)
	}
}
*/

func (p *Player) __LoadEquipDataDB(PlayerId int64) error {
	data := &model.EquipData{PlayerId: PlayerId}
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		orm.LoadObjSql(&p.EquipData, rs.Row())
	}
	return err
}

func (p *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData) {
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveEquipData", playerId)
}

/*
func (p *Player) __SaveEquipDataDB(){
	if p.EquipData.Dirty{
    	orm.DB.Exec(orm.SaveSql(p.EquipData))
		p.EquipData.Dirty = false
	}
}

func (p *Player) __SaveEquipData(data model.EquipData){
    p.EquipData = data
	p.EquipData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveEquipData", p.MailBox.Id)
}

func (p *PlayerMgr) SaveEquipData(ctx context.Context, playerId int64, data model.EquipData){
	pPlayer, bEx := p.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveEquipData(data)
	}
}
*/

func (p *Player) __LoadMailDataDB(PlayerId int64) error {
	data := &model.MailData{PlayerId: PlayerId}
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		orm.LoadObjSql(&p.MailData, rs.Row())
	}
	return err
}

func (p *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData) {
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveMailData", playerId)
}

/*
func (p *Player) __SaveMailDataDB(){
	if p.MailData.Dirty{
    	orm.DB.Exec(orm.SaveSql(p.MailData))
		p.MailData.Dirty = false
	}
}

func (p *Player) __SaveMailData(data model.MailData){
    p.MailData = data
	p.MailData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveMailData", p.MailBox.Id)
}

func (p *PlayerMgr) SaveMailData(ctx context.Context, playerId int64, data model.MailData){
	pPlayer, bEx := p.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveMailData(data)
	}
}
*/

func (p *Player) __LoadSocialDataDB(PlayerId int64) error {
	data := &model.SocialData{PlayerId: PlayerId}
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhere(data)))
	rs, err := orm.Query(rows, err)
	if err == nil && rs.Next() {
		orm.LoadObjSql(&p.SocialData, rs.Row())
	}
	return err
}

func (p *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData) {
	orm.DB.Exec(orm.SaveSql(&data))
	base.LOG.Printf("玩家[%d] SaveSocialData", playerId)
}

/*
func (p *Player) __SaveSocialDataDB(){
	if p.SocialData.Dirty{
    	orm.DB.Exec(orm.SaveSql(p.SocialData))
		p.SocialData.Dirty = false
	}
}

func (p *Player) __SaveSocialData(data model.SocialData){
    p.SocialData = data
	p.SocialData.Dirty = true
    base.LOG.Printf("玩家[%d] SaveSocialData", p.MailBox.Id)
}

func (p *PlayerMgr) SaveSocialData(ctx context.Context, playerId int64, data model.SocialData){
	pPlayer, bEx := p.m_PlayerMap[playerId]
	if bEx{
		pPlayer.__SaveSocialData(data)
	}
}
*/

func (p *Player) LoadPlayerDB(PlayerId int64) error {
	p.Init(PlayerId)
	if err := p.__LoadSimplePlayerDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadSimplePlayerDataDB() error")
		return err
	}
	if err := p.__LoadPlayerKvDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadPlayerKvDataDB() error")
		return err
	}
	if err := p.__LoadItemDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadItemDataDB() error")
		return err
	}
	if err := p.__LoadEquipDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadEquipDataDB() error")
		return err
	}
	if err := p.__LoadMailDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadMailDataDB() error")
		return err
	}
	if err := p.__LoadSocialDataDB(PlayerId); err != nil {
		base.LOG.Printf("__LoadSocialDataDB() error")
		return err
	}
	return nil
}
