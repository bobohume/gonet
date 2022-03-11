package player

import(
	"gonet/rpc"
	"gonet/server/game"
)

// 自动生成代码

func (this *Player) SaveSimplePlayerData(){
	this.SimplePlayerData.Dirty = true
}

func (this *Player) __SaveSimplePlayerDataDB(){
	if this.SimplePlayerData.Dirty{
    	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "SaveSimplePlayerData", this.SimplePlayerData.PlayerId, this.SimplePlayerData)
		this.SimplePlayerData.Dirty = false
    	game.SERVER.GetLog().Printf("玩家[%d] SaveSimplePlayerData", this.Raft.Id)
	}
}

func (this *Player) SavePlayerKvData(){
	this.PlayerKvData.Dirty = true
}

func (this *Player) __SavePlayerKvDataDB(){
	if this.PlayerKvData.Dirty{
    	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "SavePlayerKvData", this.PlayerKvData.PlayerId, this.PlayerKvData)
		this.PlayerKvData.Dirty = false
    	game.SERVER.GetLog().Printf("玩家[%d] SavePlayerKvData", this.Raft.Id)
	}
}

func (this *Player) SaveItemData(){
	this.ItemData.Dirty = true
}

func (this *Player) __SaveItemDataDB(){
	if this.ItemData.Dirty{
    	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "SaveItemData", this.ItemData.PlayerId, this.ItemData)
		this.ItemData.Dirty = false
    	game.SERVER.GetLog().Printf("玩家[%d] SaveItemData", this.Raft.Id)
	}
}

func (this *Player) SaveEquipData(){
	this.EquipData.Dirty = true
}

func (this *Player) __SaveEquipDataDB(){
	if this.EquipData.Dirty{
    	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "SaveEquipData", this.EquipData.PlayerId, this.EquipData)
		this.EquipData.Dirty = false
    	game.SERVER.GetLog().Printf("玩家[%d] SaveEquipData", this.Raft.Id)
	}
}

func (this *Player) SaveMailData(){
	this.MailData.Dirty = true
}

func (this *Player) __SaveMailDataDB(){
	if this.MailData.Dirty{
    	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "SaveMailData", this.MailData.PlayerId, this.MailData)
		this.MailData.Dirty = false
    	game.SERVER.GetLog().Printf("玩家[%d] SaveMailData", this.Raft.Id)
	}
}

func (this *Player) SaveSocialData(){
	this.SocialData.Dirty = true
}

func (this *Player) __SaveSocialDataDB(){
	if this.SocialData.Dirty{
    	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "SaveSocialData", this.SocialData.PlayerId, this.SocialData)
		this.SocialData.Dirty = false
    	game.SERVER.GetLog().Printf("玩家[%d] SaveSocialData", this.Raft.Id)
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

