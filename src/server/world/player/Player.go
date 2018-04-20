package player

import (
	"actor"
	"server/common"
	"message"
	"github.com/golang/protobuf/proto"
	"fmt"
	"server/world"
)

type(
	Player struct{
		actor.Actor
		PlayerData

		m_offlineTimer *common.SimpleTimer
	}

	IPlayer interface {
		actor.IActor
		Update()

		IsOffline() bool
		IsLogout() bool
		IsInGame() bool
	}
)

func (this* Player) Init(num int){
	this.Actor.Init(1)
	this.PlayerData.Init()
	this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	this.m_offlineTimer = common.NewSimpleTimer(120 * 1000 * 1000* 1000)

	//玩家登录
	this.RegisterCall("Login", func(caller *actor.Caller) {
		PlayerSimpleList := LoadSimplePlayerDatas(this.AccountId)
		this.PlayerSimpleDataList = PlayerSimpleList

		PlayerDataList := make([]*message.PlayerData, 1)
		for i, v := range PlayerSimpleList{
			PlayerDataList[i] = &message.PlayerData{PlayerID:proto.Int32(int32(v.PlayerId)), PlayerName:proto.String(v.PlayerName),PlayerGold:proto.Int32(int32(v.Gold))}
			this.PlayerIdList = append(this.PlayerIdList, v.PlayerId)
		}

		this.m_Log.Println("玩家登录成功")
		world.SendToClient(caller.SocketId, &message.W_C_SelectPlayerResponse{PacketHead: message.BuildPacketHead( this.AccountId,  int(message.SERVICE_CLIENT)),
			AccountId:proto.Int32(int32(this.AccountId)),
			PlayerData:PlayerDataList,
		})
	})

	//玩家登录到游戏
	this.RegisterCall("C_W_Game_LoginRequset", func(caller *actor.Caller, packet *message.C_W_Game_LoginRequset) {
		nPlayerId := int(*packet.PlayerId)
		if !this.SetPlayerId(nPlayerId){
			this.m_Log.Printf("帐号[%d]登入的玩家[%d]不存在", this.AccountId, nPlayerId)
		}

		//读取玩家数据
	})

	//创建玩家
	this.RegisterCall("C_W_CreatePlayerRequest", func(caller *actor.Caller, packet *message.C_W_CreatePlayerRequest){
		tx, _ := world.SERVER.GetDB().Begin();
		_, err :=tx.Exec(fmt.Sprintf("call `sp_checkcreatePlayer`(%d)", this.AccountId))
		if err == nil{
			row := tx.QueryRow("select @err")
			if row != nil{
				var err int
				row.Scan(&err)
				//register
				if(err == 0) {
					world.SERVER.GetServer().SendMsgByID(caller.SocketId,"W_A_CreatePlayer", this.AccountId, *packet.PlayerName, *packet.Sex)
				}else{
					this.m_Log.Printf("账号[%d]创建玩家上限", this.AccountId)
					world.SendToClient(caller.SocketId, &message.W_C_CreatePlayerResponse{
						PacketHead:message.BuildPacketHead(this.AccountId, 0 ),
						Error:proto.Int32(int32(err)),
						PlayerId:proto.Int32(int32(0)),
					})
				}
			}
		}
		tx.Commit()
	})

	//account创建玩家反馈
	this.RegisterCall("CreatePlayer", func(caller *actor.Caller, playerId int, err int) {
		//创建成功
		if err == 0{
			playerSimpleData := LoadSimplePlayerData(playerId)
			this.PlayerSimpleDataList = append(this.PlayerSimpleDataList, playerSimpleData)
			this.PlayerIdList = append(this.PlayerIdList, playerId)
		}

		world.SendToClient(caller.SocketId, &message.W_C_CreatePlayerResponse{
			PacketHead:message.BuildPacketHead(this.AccountId, 0 ),
			Error:proto.Int32(int32(err)),
			PlayerId:proto.Int32(int32(playerId)),
		})
	})

	//玩家断开链接
	this.RegisterCall("Logout", func(caller *actor.Caller, accountId int) {
		PLAYERMGR.SendMsg(caller.SocketId, "testStruct", accountId,
			[]*SimplePlayerData{&SimplePlayerData{1, 2, "test",4,5,6, 0, 1,1000, 111111111},
			&SimplePlayerData{2, 2, "test222",4,5,6, 0, 1,21000, 222222}})
		world.SERVER.GetLog().Printf("[%d] 断开链接", accountId)
		this.Stop()
	})

	//var val interface{}
	//val = this
	//world.SERVER.GetPlayerMgr().AddActor(val.(*actor.Actor));
	this.Actor.Start()
}

func (this* Player) Update(){

}