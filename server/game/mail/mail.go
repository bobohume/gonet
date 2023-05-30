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
	MailMgr struct {
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
	MGR MailMgr
)

func (m *MailMgr) Init() {
	m.Actor.Init()
	actor.MGR.RegisterActor(m)
	m.Actor.Start()
	//m.sendMail(10000238, 10000238, 1000, 60010, 10, "test", "我是大剌剌", 1)
	//m.loadMialById(2)
}

func (m *MailMgr) sendMail(sender int64, recver int64, money int, itemId int, itemNum int, title string, content string, isSystem int8) {
	mail := &model.MailItem{}
	mail.Id = base.UUID.UUID()
	mail.Sender = sender
	mail.Recver = recver
	mail.ItemId = itemId
	mail.ItemCount = itemNum
	mail.Money = money
	mail.IsSystem = isSystem
	mail.Title = title
	mail.Content = content
	//离线
	if cluster.MGR.MailBox.Get(recver) == nil {
		orm.InsertSql(mail)
		base.LOG.Printf("邮件发送给[%d]玩家成功", recver)
	} else {
		cluster.MGR.SendMsg(rpc.RpcHead{Id: recver}, "game<-Player.Add_Player_Mail", mail)
	}
}
