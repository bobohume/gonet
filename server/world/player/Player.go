package player

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gonet/actor"
	"gonet/base"
	"gonet/db"
	"gonet/rpc"
	"gonet/common"
	"gonet/server/message"
	"gonet/server/world"
)

type(
	Player struct{
		actor.Actor

		PlayerData
		m_ItemMgr      IItemMgr
		m_db           *sql.DB
		m_Log          *base.CLog
		m_offlineTimer *common.SimpleTimer
	}
)

func (this *Player) Init(num int){
	this.Actor.Init(num)
	this.PlayerData.Init()
	this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	this.m_offlineTimer = common.NewSimpleTimer(5 *60)
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_ItemMgr = &ItemMgr{}
	this.m_ItemMgr.Init(this)

	//玩家登录
	this.RegisterCall("Login", func(ctx context.Context, gateClusterId uint32, zoneClusterId uint32) {
		PlayerSimpleList := LoadSimplePlayerDatas(this.AccountId)
		this.PlayerSimpleDataList = PlayerSimpleList

		PlayerDataList := make([]*message.PlayerData, len(PlayerSimpleList))
		this.PlayerIdList = []int64{}
		for i, v := range PlayerSimpleList{
			PlayerDataList[i] = &message.PlayerData{PlayerID:v.PlayerId, PlayerName:v.PlayerName,PlayerGold:int32(v.Gold)}
			this.PlayerIdList = append(this.PlayerIdList, v.PlayerId)
		}

		this.m_Log.Println("玩家登录成功")
		this.SetGateClusterId(gateClusterId)
		this.SetZoneClusterId(zoneClusterId)
		this.SendToClient(&message.W_C_SelectPlayerResponse{PacketHead: message.BuildPacketHead( this.AccountId,  rpc.SERVICE_CLIENT),
			AccountId:this.AccountId,
			PlayerData:PlayerDataList,
		})
	})

	//玩家登录到游戏
	this.RegisterCall("C_W_Game_LoginRequset", func(ctx context.Context, packet *message.C_W_Game_LoginRequset) {
		nPlayerId := packet.GetPlayerId()
		if !this.SetPlayerId(nPlayerId){
			this.m_Log.Printf("帐号[%d]登入的玩家[%d]不存在", this.AccountId, nPlayerId)
			return
		}

		//读取玩家数据
		this.LoadPlayerData()
		//加载到地图
		this.AddMap()
		//添加到世界频道
		actor.MGR.SendMsg(rpc.RpcHead{ActorName:"chatmgr"}, "AddPlayerToChannel", this.AccountId, this.GetPlayerId(), int64(-3000), this.GetPlayerName(), this.GetGateClusterId())
	})

	//创建玩家
	this.RegisterCall("C_W_CreatePlayerRequest", func(ctx context.Context, packet *message.C_W_CreatePlayerRequest){
		rows, err := this.m_db.Query(fmt.Sprintf("select count(player_id) as player_count from tbl_player where account_id = %d", this.AccountId))
		if err == nil {
			rs := db.Query(rows, err)
			if rs.Next() {
				player_count := rs.Row().Int("player_count")
				if player_count >= 1 {
					this.m_Log.Printf("账号[%d]创建玩家上限", this.AccountId)
					world.SendToClientBySocketId(this.GetRpcHead(ctx).SocketId, &message.W_C_CreatePlayerResponse{
						PacketHead:message.BuildPacketHead(this.AccountId, 0 ),
						Error:int32(1),
						PlayerId:0,
					})
				}else{
					world.SendToAccount("W_A_CreatePlayer", this.AccountId, packet.GetPlayerName(), packet.GetSex(), this.GetRpcHead(ctx).SocketId)
				}
			}
		}
	})

	//account创建玩家反馈
	this.RegisterCall("CreatePlayer", func(ctx context.Context, playerId int64, socketId uint32, err int) {
		//创建成功
		if err == 0{
			this.PlayerIdList = []int64{}
			playerSimpleData := LoadSimplePlayerData(playerId)
			this.PlayerSimpleDataList = append(this.PlayerSimpleDataList, playerSimpleData)
			this.PlayerIdList = append(this.PlayerIdList, playerId)
		}

		world.SendToClientBySocketId(socketId, &message.W_C_CreatePlayerResponse{
			PacketHead:message.BuildPacketHead(this.AccountId, 0 ),
			Error:int32(err),
			PlayerId:playerId,
		})
	})

	//玩家断开链接
	this.RegisterCall("Logout", func(ctx context.Context, accountId int64) {
		world.SERVER.GetLog().Printf("[%d] 断开链接", accountId)
		this.SetGateClusterId(0)
		this.Stop()
		this.LeaveMap()
	})

	this.Actor.Start()
}

func (this *Player)GetDB() *sql.DB{
	return this.m_db
}

func (this *Player) GetLog() *base.CLog{
	return this.m_Log
}

func (this *Player) SendToClient(packet proto.Message){
	world.SendToClient(this.GetGateClusterId(), packet)
}

func (this *Player) Update(){

}