package mail

import (
	"actor"
	"database/sql"
	"db"
	"fmt"
	"server/world"
)

type (
	MailItem struct{
		Id int`primary`
		Sender int
		SenderName string
		Recver int
		RecverName string
		Money int
		ItemId int
		ItemCount int
		IsRead int8
		IsSystem int8
		RecvFlag int8
		Title string
		Content string
	}

	CMailMgr struct {
		actor.Actor
		m_db *sql.DB
	}

	IMailMgr interface {
		actor.IActor
		SendMail(int, int, int, int, int, string, string, bool)
		LoadMail(int, []*MailItem, int, int)
		LoadMialById(int) *MailItem
		DeleteMail(int, int)
		ReadMail(int, int)
		RecverMail(int, int)
	}
)

var(
	MAILMGR CMailMgr
)

func (this *CMailMgr) Init(num int) {
	this.m_db = world.SERVER.GetDB()
	this.Actor.Init(num)
	actor.GetGActorList().RegisterGActorList("mail", this)

	this.Actor.Start()
}

func (this *CMailMgr) SendMail(sender int, recver int, money int, itemId int, itemNum int, title string, content string, isSystem int8){
	m := &MailItem{}
	m.Sender = sender
	m.Recver = recver
	m.ItemId = itemId
	m.ItemCount = itemNum
	m.Money = money
	m.IsSystem = isSystem
	m.Title = title
	m.Content = content

	tx, _ := this.m_db.Begin()
	_, err :=tx.Exec(fmt.Sprintf("call `sp_updatemail`(%d,%d,'%s',%d,%d,%d,%d,'%s',%d,'%s','%s')", 0, sender, "",money, itemId, itemNum, recver, "", isSystem, title, content))
	if err == nil{
		row := tx.QueryRow("select @err, @mailid, @recver")
		if row != nil{
			var err int
			row.Scan(&err, &m.Id, &m.Recver)
			//register
			if(err == 0) {
				world.SERVER.GetLog().Printf("邮件发送给[%d]玩家成功", recver)
			}else{
				world.SERVER.GetLog().Printf("账号[%d]创建玩家失败", recver)
			}
			/*world.SendToClient(caller.SocketId, &message.W_C_CreatePlayerResponse{
				PacketHead:message.BuildPacketHead(this.AccountId, 0 ),
				Error:proto.Int32(int32(err)),
				PlayerId:proto.Int32(int32(playerId)),
			})*/
		}
	}
	tx.Commit()
}

func (this *CMailMgr) LoadMail(playerid int, mailList []*MailItem, recvCount int, noReadCount int){
	//this.SendMail(50000055, 50000055, 1000, 60010, 10, "test", "test1111", 1)
	rows, err := this.m_db.Query(db.LoadSql(MailItem{}, "tbl_mail", fmt.Sprintf("recver=%d", playerid)))
	if err == nil{
		for rows.Next(){
			m := &MailItem{}
			err = rows.Scan(&m.Id, &m.Sender, &m.SenderName, &m.Recver, &m.RecverName, &m.Money, &m.ItemId,
				&m.ItemCount, &m.IsRead, &m.IsSystem, &m.RecvFlag, &m.Title, &m.Content)
			if err != nil{
				world.SERVER.GetLog().Printf("load mail err[%s]", err.Error())
			}else{
				mailList = append(mailList, m)
				recvCount++
				if m.IsRead == 0{
					noReadCount++
				}
				//fmt.Println(m)
				world.SERVER.GetLog().Printf("读取玩家[%d]邮件成功", playerid)
			}
		}
	}
}

func (this *CMailMgr) LoadMialById(mailid int) *MailItem{
	m := &MailItem{}
	row := this.m_db.QueryRow(db.LoadSql(m, "tbl_mail", fmt.Sprintf("id=%d", mailid)))
	if row != nil{
		err := row.Scan(&m.Id, &m.Sender, &m.SenderName, &m.Recver, &m.RecverName, &m.Money, &m.ItemId,
			&m.ItemCount, &m.IsRead, &m.IsSystem, &m.RecvFlag, &m.Title, &m.Content)
		if err == nil{
			return m
		}
	}
	return nil
}

func (this *CMailMgr) DeleteMail(playerid int, mailid int){
	this.m_db.Exec("delete form tbl_mail where playerid=%d and id =%d", playerid, mailid)
}

func (this *CMailMgr) ReadMail(playerid int, mailid int){
	m := this.LoadMialById(mailid)
	m.IsRead = 1

	if m.Recver != playerid{
		return
	}

	//文本邮件看完就删除掉
	if m.ItemId == 0 && m.Money == 0 {
		this.DeleteMail(m.Recver, m.Id)
	}else{
		this.m_db.Exec(db.UpdateSqlEx(m, "tb_mail", "Id", "IsRead"))
	}
}

func (this *CMailMgr) RecverMail(playerid int, mailid int){
	m := this.LoadMialById(mailid)
	if m.Recver != playerid{
		return
	}

	if m.RecvFlag == 0{
		m.RecvFlag = 1
		this.m_db.Exec(db.UpdateSqlEx(m, "tb_mail", "Id", "RecvFlag"))
		//奖励道具

	}

	this.DeleteMail(playerid, mailid)
}