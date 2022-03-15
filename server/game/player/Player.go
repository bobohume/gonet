package player

import (
	"context"
	"gonet/actor"
	"gonet/common/cluster"
	"gonet/rpc"
	"gonet/server/game"
	"gonet/server/model"
	"time"

	"github.com/golang/protobuf/proto"
)

type (
	Player struct {
		actor.Actor

		model.PlayerData
		Raft         		 rpc.PlayerClusterInfo
		PlayerId             int64
		GateClusterId        uint32
		m_offline_flag bool//离线
		m_InGameFlag   bool//登录游戏
	}
)

func (this *Player) Init() {
	this.Actor.Init()
	this.RegisterTimer((cluster.OFFLINE_TIME/3)*time.Second, this.UpdateLease) //定时器
	this.RegisterTimer(60*time.Second, this.SaveDB) //定时器
	this.Actor.Start()
}

func (this *Player) SendToClient(packet proto.Message) {
	game.SendToClient(this.GetGateClusterId(), packet)
}

func (this *Player) UpdateLease() {
	if ! this.m_offline_flag{
		game.SERVER.GetPlayerRaft().Lease(this.Raft.LeaseId)
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
func (this *Player) Login(ctx context.Context, gateClusterId uint32, clusterInfo rpc.PlayerClusterInfo) {
	this.SetGateClusterId(gateClusterId)
	this.Raft = clusterInfo
	game.SERVER.GetLog().Println("玩家登录成功")
	//加载玩家数据
	game.SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType: rpc.SERVICE_DB, ClusterId:this.Raft.DClusterId}, "Load_Player_DB", this.PlayerId, this.Raft)
}

//断线重连
func (this *Player) ReLogin (ctx context.Context, gateClusterId uint32, clusterInfo rpc.PlayerClusterInfo) {
	game.SERVER.GetLog().Printf("[%d] 重连成功", this.PlayerId)
	this.SetGateClusterId(gateClusterId)
	this.Raft = clusterInfo
	if this.m_InGameFlag{
		this.m_offline_flag = false
		this.UpdateLease()
		this.LoginFinish()
	}else{
		this.Login(ctx, gateClusterId, clusterInfo)
	}
}


//加载玩家结束
func (this *Player) Load_Player_DB_Finish(ctx context.Context, data model.PlayerData) {
	this.m_InGameFlag = true
	this.PlayerData = data
	//加载到地图
	this.LoginFinish()
}

func (this *Player) LoginFinish() {
	//加载到地图
	this.AddMap()
	game.SendToGM(rpc.RpcHead{Id:this.PlayerId}, "AddPlayerToChannel", this.PlayerId, int64(-3000), this.PlayerName, this.GetGateClusterId())
}

//创建玩家
/*func (this *Player) C_W_CreatePlayerRequest(ctx context.Context, packet *message.C_W_CreatePlayerRequest) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Do_CreatePlayer", this.PlayerId, packet.GetPlayerName(), packet.GetSex(), this.GetRpcHead(ctx).SrcClusterId)
}

//account创建玩家反馈
func (this *Player) CreatePlayer(ctx context.Context, playerId int64, playername string, sex int,  err int) {
	//创建成功
	this.SendToClient(&message.W_C_CreatePlayerResponse{
		PacketHead: message.BuildPacketHead(this.PlayerId, 0),
		Error:      int32(err),
		PlayerId:   playerId,
	})
}*/

//玩家断开链接
func (this *Player) Logout (ctx context.Context, playerId int64) {
	game.SERVER.GetLog().Printf("[%d] 断开链接", playerId)
	this.m_offline_flag = true
}

//lease过期
func (this *Player) Player_Lease_Expire(ctx context.Context) {
	game.SERVER.GetLog().Printf("[%d] 过期删除玩家", this.PlayerId)
	actor.MGR.DelPlayer(this.PlayerId)
	this.SetGateClusterId(0)
	this.Stop()
	this.LeaveMap()
}
