package player

import (
	"gonet/server/model"
)

func (p *Player) loadMialById(mailId int64) *model.MailItem {
	m, bEx := p.MailData.DataMap[mailId]
	if bEx {
		return m
	}
	return nil
}

func (p *Player) deleteMail(playerId int64, mailId int64) {
	delete(p.MailData.DataMap, mailId)
	p.SaveMailData()
}

func (p *Player) readMail(playerId int64, mailId int64) {
	m := p.loadMialById(mailId)
	m.IsRead = 1

	if m.Recver != playerId {
		return
	}

	//文本邮件看完就删除掉
	if m.ItemId == 0 && m.Money == 0 {
		p.deleteMail(m.Recver, m.Id)
	}

	p.SaveMailData()
}

func (p *Player) recverMail(playerId int64, mailId int64) {
	m := p.loadMialById(mailId)
	if m.Recver != playerId {
		return
	}

	if m.RecvFlag == 0 {
		m.RecvFlag = 1
		//奖励道具

	}

	p.deleteMail(playerId, mailId)
	p.SaveMailData()
}
