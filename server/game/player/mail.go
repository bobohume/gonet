package player

import (
	"gonet/server/model"
)

func (this *Player) loadMialById(mailId int64) *model.MailItem{
	m, bEx := this.MailData.DataMap[mailId]
	if bEx {
		return m
	}
	return nil
}

func (this *Player) deleteMail(playerId int64, mailId int64){
	delete(this.MailData.DataMap, mailId)
	this.SaveMailData()
}

func (this *Player) readMail(playerId int64, mailId int64){
	m := this.loadMialById(mailId)
	m.IsRead = 1

	if m.Recver != playerId{
		return
	}

	//文本邮件看完就删除掉
	if m.ItemId == 0 && m.Money == 0 {
		this.deleteMail(m.Recver, m.Id)
	}

	this.SaveMailData()
}

func (this *Player) recverMail(playerId int64, mailId int64){
	m := this.loadMialById(mailId)
	if m.Recver != playerId{
		return
	}

	if m.RecvFlag == 0{
		m.RecvFlag = 1
		//奖励道具

	}

	this.deleteMail(playerId, mailId)
	this.SaveMailData()
}