package player

import (
	"actor"
	"fmt"
	"base"
	"strings"
	"database/sql"
	"message"
	"server/world"
)
//********************************************************
// 玩家管理
//********************************************************
var(
	PLAYERMGR PlayerMgr
	PLAYER Player
)
type(
	PlayerMgr struct{
		actor.Actor
		m_PlayerMap map[int] *Player
		m_db *sql.DB
		m_Log *base.CLog
		//m_Lock sync.Locker
	}

	IPlayerMgr interface {
		actor.IActor
		CreatePlayer(int, string, int, int)

		GetPlayer(accountId int) *Player
		AddPlayer(accountId int) *Player
		RemovePlayer(accountId int)
	}
)

func (this* PlayerMgr) Init(num int){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_PlayerMap = make(map[int] *Player)
	//this.m_Lock = &sync.Mutex{}
	this.Actor.Init(num)
	actor.GetGActorList().RegisterGActorList("playermgr", this)
	//玩家登录
	this.RegisterCall("G_W_CLoginRequest", func(caller *actor.Caller, accountId int) {
		pPlayer := this.GetPlayer(accountId)
		if pPlayer == nil{
			pPlayer = this.AddPlayer(accountId)
		}

		if pPlayer != nil{
			actor.SendMsg(&pPlayer.Actor, caller.SocketId, "Login")
		}
	})

	//玩家断开链接
	this.RegisterCall("G_ClientLost", func(caller *actor.Caller, accountId int) {
		pPlayer := this.GetPlayer(accountId)
		if pPlayer != nil{
			pPlayer.SendMsg(caller.SocketId, "Logout", accountId)
		}

		this.RemovePlayer(accountId)
	})

	//account创建玩家反馈， 考虑到在创建角色的时候退出的情况
	this.RegisterCall("A_W_CreatePlayer", func(caller *actor.Caller, accountId int, playerId int, playername string, sex int32) {
		tx, _ := this.m_db.Begin();
		_, err :=tx.Exec(fmt.Sprintf("call `sp_createPlayer`(%d,'%s',%d, %d)", accountId, playername, sex, playerId))
		if err == nil{
			row := tx.QueryRow("select @err, @playerId")
			if row != nil{
				var err int
				row.Scan(&err, &playerId)
				//register
				if(err == 0) {
					this.m_Log.Printf("账号[%d]创建玩家[%d]", accountId, playerId)
				}else{
					this.m_Log.Printf("账号[%d]创建玩家失败", accountId)
					world.SERVER.GetServer().SendMsgByID(caller.SocketId,"W_A_DeletePlayer", accountId, playerId)
				}

				//通知玩家`
				pPlayer := this.GetPlayer(accountId)
				if pPlayer != nil{
					pPlayer.SendMsg(caller.SocketId, "CreatePlayer", playerId, err)
				}
			}
		}
		tx.Commit()
	})

	this.RegisterCall("testStruct", func(caller *actor.Caller, accountId int, pdata []*SimplePlayerData){
		for i,v := range pdata{
			fmt.Println(i, v)
		}
	})
	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	PLAYER.Init(1)
	this.Actor.Start()
}

func (this *PlayerMgr) GetPlayer(accountId int) *Player{
	//this.m_Lock.Lock()
	pPlayer, exist := this.m_PlayerMap[accountId]
	//this.m_Lock.Unlock()
	if exist{
		return pPlayer
	}
	return nil
}

func (this *PlayerMgr) AddPlayer(accountId int) *Player{
	LoadPlayerDB := func(accountId int) ([]int, int){
		PlayerList := make([]int, 0)
		PlayerNum := 0
		rows, err := this.m_db.Query(fmt.Sprintf("select playerId where accountId=%d", accountId))
		if err == nil{
			for rows.Next(){
				rows.Scan(&PlayerList[PlayerNum])
				PlayerNum++
			}
		}
		return PlayerList, PlayerNum
	}

	fmt.Printf("玩家[%d]登录", accountId)
	PlayerList, PlayerNum := LoadPlayerDB(accountId)
	pPlayer := &Player{}
	pPlayer.AccountId = accountId
	pPlayer.PlayerIdList = PlayerList
	pPlayer.PlayerNum = PlayerNum
	//this.m_Lock.Lock()
	this.m_PlayerMap[accountId] = pPlayer
	//this.m_Lock.Unlock()
	pPlayer.Init(1)
	return pPlayer
}

func (this *PlayerMgr) RemovePlayer(accountId int){
	this.m_Log.Printf("移除帐号数据[%d]", accountId)
	//this.m_Lock.Lock()
	delete(this.m_PlayerMap, accountId)
	//this.m_Lock.Unlock()
}

func (this *PlayerMgr) PacketFunc(id int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("PacketFunc", err)
		}
	}()

	SendToPlayer := func(AccountId int, io actor.CallIO) {
		defer func() {
			if err := recover(); err != nil{
				fmt.Println("SendToPlayer", err)
			}
		}()
		pPlayer := this.GetPlayer(AccountId)
		if pPlayer != nil{
			go actor.SendActor(pPlayer, io)
		}
	}

	var io actor.CallIO
	io.Buff = buff
	io.SocketId = id

	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	funcName = strings.ToLower(funcName)
	pFunc := this.FindCall(funcName)
	if pFunc != nil{
		this.Send(io)
		return true
	}else{
		pFunc := PLAYER.FindCall(funcName)
		if pFunc != nil{
			bitstream.ReadInt(base.Bit8)
			nType := bitstream.ReadInt(base.Bit8);
			if(nType == 8 || nType == 9 || nType == 13){
				nAccountId := bitstream.ReadInt(base.Bit32)
				SendToPlayer(nAccountId, io)
			}else if (nType == 31){
				packet := message.GetPakcetByName(funcName)
				message.UnmarshalText(packet, bitstream.ReadString())
				packetHead := message.GetPakcetHead(packet)
				nAccountId := int(*packetHead.Id)
				SendToPlayer(nAccountId, io)
			}
		}
	}

	return false
}