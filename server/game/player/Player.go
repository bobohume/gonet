package player

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/common/cluster"
	"gonet/rpc"
	"gonet/server/game"
	"gonet/server/model"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	MAILBOX_TL_TIME = cluster.MAILBOX_TL_TIME / 2
)

type (
	Player struct {
		actor.Actor

		model.PlayerData
		MailBox       rpc.MailBox
		PlayerId      int64
		GateClusterId uint32
		ZoneClusterId uint32
		offline_flag  bool //离线
		online_time   int64
		isInGame      bool //登录游戏
	}
)

func (p *Player) Init() {
	p.Actor.Init()
	p.RegisterTimer((MAILBOX_TL_TIME/2)*time.Second, p.UpdateLease) //定时器
	p.RegisterTimer(60*time.Second, p.SaveDB)                       //定时器
	p.online_time = time.Now().Unix()
	p.Actor.Start()
}

func (p *Player) SendToClient(packet proto.Message) {
	game.SendToClient(p.GetGateClusterId(), packet)
}

func (p *Player) UpdateLease() {
	if !p.offline_flag {
		p.online_time = time.Now().Unix()
		cluster.MGR.MailBox.Lease(p.MailBox.LeaseId)
	} else if time.Now().Unix()-p.online_time >= MAILBOX_TL_TIME {
		cluster.MGR.MailBox.Delete(p.PlayerId)
	}
}

func (p *Player) SaveDB() {
	p.SavePlayerDB()
}

func (p *Player) SetGateClusterId(clusterId uint32) {
	p.GateClusterId = clusterId
}

func (p *Player) GetGateClusterId() uint32 {
	return p.GateClusterId
}

func (p *Player) GetPlayerId() int64 {
	return p.PlayerId
}

//玩家登录
func (p *Player) Login(ctx context.Context, gateClusterId uint32, mailbox rpc.MailBox) {
	p.SetGateClusterId(gateClusterId)
	p.MailBox = mailbox
	base.LOG.Println("玩家登录成功")
	//加载玩家数据
	cluster.MGR.SendMsg(rpc.RpcHead{Id: p.PlayerId}, "db<-PlayerMgr.Load_Player_DB", p.PlayerId, p.MailBox)
}

//断线重连
func (p *Player) ReLogin(ctx context.Context, gateClusterId uint32, mailbox rpc.MailBox) {
	base.LOG.Printf("[%d] 重连成功", p.PlayerId)
	p.SetGateClusterId(gateClusterId)
	p.MailBox = mailbox
	if p.isInGame {
		p.offline_flag = false
		p.online_time = time.Now().Unix()
		p.LoginFinish()
	} else {
		p.Login(ctx, gateClusterId, mailbox)
	}
}

//加载玩家结束
func (p *Player) Load_Player_DB_Finish(ctx context.Context, data model.PlayerData) {
	p.isInGame = true
	p.PlayerData = data
	//加载到地图
	p.LoginFinish()
	p.ZoneClusterId = cluster.MGR.RandomCluster(rpc.RpcHead{Id: p.PlayerId, DestServerType: rpc.SERVICE_ZONE}).ClusterId
}

func (p *Player) LoginFinish() {
	//加载到地图
	p.AddMap()
	game.SendToGM(rpc.RpcHead{Id: p.PlayerId}, "ChatMgr.AddPlayerToChannel", p.PlayerId, int64(-3000), p.PlayerName, p.GetGateClusterId())
}

//玩家断开链接
func (p *Player) Logout(ctx context.Context, playerId int64) {
	base.LOG.Printf("[%d] 断开链接", playerId)
	p.offline_flag = true
	p.SaveDB()
}

//lease过期
func (p *Player) OnUnRegister(ctx context.Context) {
	base.LOG.Printf("[%d] 过期删除玩家", p.PlayerId)
	MGR.DelActor(p.PlayerId)
	p.SetGateClusterId(0)
	p.Stop()
	p.LeaveMap()
	cluster.MGR.SendMsg(rpc.RpcHead{Id: p.AccountId}, "gm<-AccountMgr.OnUnRegister")
}
