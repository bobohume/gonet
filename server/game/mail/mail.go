package mail

import (
	"database/sql"
	"gonet/actor"
	"gonet/base"
	"gonet/orm"
	"gonet/rpc"
	"gonet/server/game"
	"gonet/server/model"
)

type (
	CMailMgr struct {
		actor.Actor
		m_db *sql.DB
	}

	IMailMgr interface {
		actor.IActor

		sendMail(sender int64, recver int64, money int, itemId int, itemNum int, title string, content string, isSystem int8)
		loadMail(playerId int64, mailList []*model.MailItem, recvCount int, noReadCount int)
		loadMialById(mailId int64) *model.MailItem
		deleteMail(playerId int64, mailId int64)
		readMail(playerId int64, mailId int64)
		recverMail(playerId int64, mailId int64)
	}
)

var(
	MGR CMailMgr
)

func (this *CMailMgr) Init() {
	this.Actor.Init()
	this.m_db = game.SERVER.GetDB()
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
	//this.sendMail(10000238, 10000238, 1000, 60010, 10, "test", "我是大剌剌", 1)
	//this.loadMialById(2)
}

func (this *CMailMgr) sendMail(sender int64, recver int64, money int, itemId int, itemNum int, title string, content string, isSystem int8){
	m := &model.MailItem{}
	m.Id = base.UUID.UUID()
	m.Sender = sender
	m.Recver = recver
	m.ItemId = itemId
	m.ItemCount = itemNum
	m.Money = money
	m.IsSystem = isSystem
	m.Title = title
	m.Content = content
	//离线
	if game.SERVER.GetPlayerRaft().GetPlayer(recver) == nil{
		this.m_db.Exec(orm.InsertSql(m))
		game.SERVER.GetLog().Printf("邮件发送给[%d]玩家成功", recver)
	}else{
		game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType: rpc.SERVICE_GAME, Id:recver}, "Add_Player_Mail", m)
	}
}