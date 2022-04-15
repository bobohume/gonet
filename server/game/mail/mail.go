package mail

import (
	"gonet/actor"
	"gonet/base"
	"gonet/common/cluster"
	"gonet/orm"
	"gonet/rpc"
	"gonet/server/model"
)

type (
	CMailMgr struct {
		actor.Actor
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

var (
	MGR CMailMgr
)

func (this *CMailMgr) Init() {
	this.Actor.Init()
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
	//this.sendMail(10000238, 10000238, 1000, 60010, 10, "test", "我是大剌剌", 1)
	//this.loadMialById(2)
}

func (this *CMailMgr) sendMail(sender int64, recver int64, money int, itemId int, itemNum int, title string, content string, isSystem int8) {
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
	if cluster.MGR.MailBox.Get(recver) == nil {
		orm.DB.Exec(orm.InsertSql(m))
		base.LOG.Printf("邮件发送给[%d]玩家成功", recver)
	} else {
		cluster.MGR.SendMsg(rpc.RpcHead{Id: recver}, "game<-Player.Add_Player_Mail", m)
	}
}
