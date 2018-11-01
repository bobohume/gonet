package mail

import (
	"actor"
	"database/sql"
	"db"
	"fmt"
	"server/world"
)

const(
	sqlTable = "tbl_mail"
)

type (
	MailItem struct{
		Id int`sql:"primary;name:id"`
		Sender int `sql:"name:sender"`
		SenderName string `sql:"name:sender_name"`
		Recver int `sql:"name:recver"`
		RecverName string `sql:"name:recver_name"`
		Money int `sql:"name:money"`
		ItemId int `sql:"name:item_id"`
		ItemCount int `sql:"name:item_count"`
		IsRead int8 `sql:"name:is_read"`
		IsSystem int8 `sql:"name:is_system"`
		RecvFlag int8 `sql:"name:recv_flag"`
		Title string `sql:"name:title"`
		Content string `sql:"name:content"`
	}

	CMailMgr struct {
		actor.Actor
		m_db *sql.DB
	}

	IMailMgr interface {
		actor.IActor

		sendMail(int, int, int, int, int, string, string, int8)
		loadMail(int, []*MailItem, int, int)
		loadMialById(int) *MailItem
		deleteMail(int, int)
		readMail(int, int)
		recverMail(int, int)
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
	//this.sendMail(10000238, 10000238, 1000, 60010, 10, "test", "我是大剌剌", 1)
	//this.loadMialById(2)
}

func (this *CMailMgr) sendMail(sender int, recver int, money int, itemId int, itemNum int, title string, content string, isSystem int8){
	m := &MailItem{}
	m.Sender = sender
	m.Recver = recver
	m.ItemId = itemId
	m.ItemCount = itemNum
	m.Money = money
	m.IsSystem = isSystem
	m.Title = title
	m.Content = content


	rows, err := this.m_db.Query(fmt.Sprintf("call `sp_updatemail`(%d,%d,'%s',%d,%d,%d,%d,'%s',%d,'%s','%s')", 0, sender, "",money, itemId, itemNum, recver, "", isSystem, title, content))
	if err == nil{
		rows.Next()
		rows.Next()
		rows.Next()
		if rows.NextResultSet(){
			rs := db.Query(rows)
			if rs.Next(){
				err := rs.Row().Int("@err")
				m.Id = rs.Row().Int("@mailid")
				m.Recver = rs.Row().Int("@recver")
				//register
				if(err == 0) {
					world.SERVER.GetLog().Printf("邮件发送给[%d]玩家成功", recver)
				}else{
					world.SERVER.GetLog().Printf("玩家[%d]邮件发送给失败", recver)
				}
				/*world.SendToClient(caller.SocketId, &message.W_C_CreatePlayerResponse{
					PacketHead:message.BuildPacketHead(this.AccountId, 0 ),
					Error:proto.Int32(int32(err)),
					PlayerId:proto.Int32(int32(playerId)),
				})*/
			}
		}
	}
}

func loadMail(row db.IRow, m *MailItem){
	m.Id = row.Int("id")
	m.Sender = row.Int("sender")
	m.SenderName = row.String("sender_name")
	m.Recver = row.Int("recver")
	m.RecverName = row.String("recver_name")
	m.Money = row.Int("money")
	m.ItemId = row.Int("item_id")
	m.ItemCount = row.Int("item_count")
	m.IsRead = int8(row.Int("is_read"))
	m.IsSystem = int8(row.Int("is_system"))
	m.RecvFlag = int8(row.Int("recv_flag"))
	m.Title = row.String("title")
	m.Content = row.String("content")
}

func (this *CMailMgr) loadMail(playerid int, mailList []*MailItem, recvCount int, noReadCount int){
	rows, err := this.m_db.Query(db.LoadSql(MailItem{}, "tbl_mail", fmt.Sprintf("recver=%d", playerid)))
	if err == nil{
		rs := db.Query(rows)
		if rs.Next(){
			m := &MailItem{}
			loadMail(rs.Row(), m)
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

func (this *CMailMgr) loadMialById(mailid int) *MailItem{
	m := &MailItem{}
	rows, err := this.m_db.Query(db.LoadSql(m, "tbl_mail", fmt.Sprintf("id=%d", mailid)))
	rs := db.Query(rows)
	if err == nil && rs.Next() {
		loadMail(rs.Row(), m)
		return m
	}
	return nil
}

func (this *CMailMgr) deleteMail(playerid int, mailid int){
	this.m_db.Exec("delete form tbl_mail where playerid=%d and id =%d", playerid, mailid)
}

func (this *CMailMgr) readMail(playerid int, mailid int){
	m := this.loadMialById(mailid)
	m.IsRead = 1

	if m.Recver != playerid{
		return
	}

	//文本邮件看完就删除掉
	if m.ItemId == 0 && m.Money == 0 {
		this.deleteMail(m.Recver, m.Id)
	}else{
		this.m_db.Exec(db.UpdateSqlEx(m, "tb_mail", "id", "is_read"))
	}
}

func (this *CMailMgr) recverMail(playerid int, mailid int){
	m := this.loadMialById(mailid)
	if m.Recver != playerid{
		return
	}

	if m.RecvFlag == 0{
		m.RecvFlag = 1
		this.m_db.Exec(db.UpdateSqlEx(m, "tb_mail", "id", "recv_flag"))
		//奖励道具

	}

	this.deleteMail(playerid, mailid)
}