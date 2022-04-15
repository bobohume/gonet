package player

import(
	"gonet/base"
	"gonet/common/cluster"
	"gonet/rpc"
)

// 自动生成代码

func (this *Player) SaveSimplePlayerData(){
	this.SimplePlayerData.Dirty = true
}

func (this *Player) __SaveSimplePlayerDataDB(){
	if this.SimplePlayerData.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.SaveSimplePlayerData", this.SimplePlayerData.PlayerId, this.SimplePlayerData)
		this.SimplePlayerData.Dirty = false
    	base.LOG.Printf("玩家[%d] SaveSimplePlayerData", this.MailBox.Id)
	}
}

func (this *Player) SavePlayerKvData(){
	this.PlayerKvData.Dirty = true
}

func (this *Player) __SavePlayerKvDataDB(){
	if this.PlayerKvData.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.SavePlayerKvData", this.PlayerKvData.PlayerId, this.PlayerKvData)
		this.PlayerKvData.Dirty = false
    	base.LOG.Printf("玩家[%d] SavePlayerKvData", this.MailBox.Id)
	}
}

func (this *Player) SaveItemData(){
	this.ItemData.Dirty = true
}

func (this *Player) __SaveItemDataDB(){
	if this.ItemData.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.SaveItemData", this.ItemData.PlayerId, this.ItemData)
		this.ItemData.Dirty = false
    	base.LOG.Printf("玩家[%d] SaveItemData", this.MailBox.Id)
	}
}

func (this *Player) SaveEquipData(){
	this.EquipData.Dirty = true
}

func (this *Player) __SaveEquipDataDB(){
	if this.EquipData.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.SaveEquipData", this.EquipData.PlayerId, this.EquipData)
		this.EquipData.Dirty = false
    	base.LOG.Printf("玩家[%d] SaveEquipData", this.MailBox.Id)
	}
}

func (this *Player) SaveMailData(){
	this.MailData.Dirty = true
}

func (this *Player) __SaveMailDataDB(){
	if this.MailData.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.SaveMailData", this.MailData.PlayerId, this.MailData)
		this.MailData.Dirty = false
    	base.LOG.Printf("玩家[%d] SaveMailData", this.MailBox.Id)
	}
}

func (this *Player) SaveSocialData(){
	this.SocialData.Dirty = true
}

func (this *Player) __SaveSocialDataDB(){
	if this.SocialData.Dirty{
    	cluster.MGR.SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, Id:this.MailBox.Id}, "PlayerMgr.SaveSocialData", this.SocialData.PlayerId, this.SocialData)
		this.SocialData.Dirty = false
    	base.LOG.Printf("玩家[%d] SaveSocialData", this.MailBox.Id)
	}
}

func (this *Player) SavePlayerDB(){
    this.__SaveSimplePlayerDataDB()
    this.__SavePlayerKvDataDB()
    this.__SaveItemDataDB()
    this.__SaveEquipDataDB()
    this.__SaveMailDataDB()
    this.__SaveSocialDataDB()
}

