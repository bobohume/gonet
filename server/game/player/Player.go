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
		MailBox        rpc.MailBox
		PlayerId       int64
		GateClusterId  uint32
		ZoneClusterId  uint32
		m_offline_flag bool //离线
		m_online_time  int64
		m_InGameFlag   bool //登录游戏
	}
)

func (this *Player) Init() {
	this.Actor.Init()
	this.RegisterTimer((MAILBOX_TL_TIME/2)*time.Second, this.UpdateLease) //定时器
	this.RegisterTimer(60*time.Second, this.SaveDB)                       //定时器
	this.m_online_time = time.Now().Unix()
	this.Actor.Start()
}

func (this *Player) SendToClient(packet proto.Message) {
	game.SendToClient(this.GetGateClusterId(), packet)
}

func (this *Player) UpdateLease() {
	if !this.m_offline_flag {
		this.m_online_time = time.Now().Unix()
		cluster.MGR.MailBox.Lease(this.MailBox.LeaseId)
	} else if time.Now().Unix()-this.m_online_time >= MAILBOX_TL_TIME {
		cluster.MGR.MailBox.Delete(this.PlayerId)
	}
}

func (this *Player) SaveDB() {
	this.SavePlayerDB()
}

func (this *Player) SetGateClusterId(clusterId uint32) {
	this.GateClusterId = clusterId
}

func (this *Player) GetGateClusterId() uint32 {
	return this.GateClusterId
}

func (this *Player) GetPlayerId() int64 {
	return this.PlayerId
}

//玩家登录
func (this *Player) Login(ctx context.Context, gateClusterId uint32, mailbox rpc.MailBox) {
	this.SetGateClusterId(gateClusterId)
	this.MailBox = mailbox
	base.LOG.Println("玩家登录成功")
	//加载玩家数据
	cluster.MGR.SendMsg(rpc.RpcHead{Id: this.PlayerId}, "db<-PlayerMgr.Load_Player_DB", this.PlayerId, this.MailBox)
}

//断线重连
func (this *Player) ReLogin(ctx context.Context, gateClusterId uint32, mailbox rpc.MailBox) {
	base.LOG.Printf("[%d] 重连成功", this.PlayerId)
	this.SetGateClusterId(gateClusterId)
	this.MailBox = mailbox
	if this.m_InGameFlag {
		this.m_offline_flag = false
		this.m_online_time = time.Now().Unix()
		this.LoginFinish()
	} else {
		this.Login(ctx, gateClusterId, mailbox)
	}
}

//加载玩家结束
func (this *Player) Load_Player_DB_Finish(ctx context.Context, data model.PlayerData) {
	this.m_InGameFlag = true
	this.PlayerData = data
	//加载到地图
	this.LoginFinish()
	this.ZoneClusterId = cluster.MGR.RandomCluster(rpc.RpcHead{Id: this.PlayerId, DestServerType: rpc.SERVICE_ZONE}).ClusterId
}

func (this *Player) LoginFinish() {
	//加载到地图
	this.AddMap()
	game.SendToGM(rpc.RpcHead{Id: this.PlayerId}, "ChatMgr.AddPlayerToChannel", this.PlayerId, int64(-3000), this.PlayerName, this.GetGateClusterId())
}

//玩家断开链接
func (this *Player) Logout(ctx context.Context, playerId int64) {
	base.LOG.Printf("[%d] 断开链接", playerId)
	this.m_offline_flag = true
	this.SaveDB()
}

//lease过期
func (this *Player) Player_On_UnRegister(ctx context.Context) {
	base.LOG.Printf("[%d] 过期删除玩家", this.PlayerId)
	MGR.DelPlayer(this.PlayerId)
	this.SetGateClusterId(0)
	this.Stop()
	this.LeaveMap()
}
