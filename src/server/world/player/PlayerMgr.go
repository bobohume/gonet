package player

import (
	"actor"
	"db"
	"fmt"
	"base"
	"strings"
	"database/sql"
	"message"
	"server/world"
	"sync"
	"github.com/golang/protobuf/proto"
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
		m_Lock *sync.RWMutex
	}

	IPlayerMgr interface {
		actor.IActor

		GetPlayer(accountId int) *Player
		AddPlayer(accountId int) *Player
		RemovePlayer(accountId int)
	}
)

func (this* PlayerMgr) Init(num int){
	this.m_db = world.SERVER.GetDB()
	this.m_Log = world.SERVER.GetLog()
	this.m_PlayerMap = make(map[int] *Player)
	this.m_Lock = &sync.RWMutex{}
	this.Actor.Init(num)
	actor.GetGActorList().RegisterGActorList("playermgr", this)
	//玩家登录
	this.RegisterCall("G_W_CLoginRequest", func(accountId int) {
		pPlayer := this.GetPlayer(accountId)
		if pPlayer != nil{
			pPlayer.SendMsg("Logout", accountId)
			this.RemovePlayer(accountId)
		}

		pPlayer = this.AddPlayer(accountId)
		pPlayer.SendMsg("Login", this.GetSocketId())
	})

	//玩家断开链接
	this.RegisterCall("G_ClientLost", func(accountId int) {
		pPlayer := this.GetPlayer(accountId)
		if pPlayer != nil{
			pPlayer.SendMsg("Logout", accountId)
		}

		this.RemovePlayer(accountId)
	})

	//account创建玩家反馈， 考虑到在创建角色的时候退出的情况
	this.RegisterCall("A_W_CreatePlayer", func(accountId int, playerId int, playername string, sex int32) {
		rows, err := this.m_db.Query(fmt.Sprintf("call `sp_createPlayer`(%d,'%s',%d, %d)", accountId, playername, sex, playerId))
		if err == nil{
			rows.Next()
			rows.Next()
			if rows.NextResultSet(){
				rs := db.Query(rows)
				if rs.Next(){
					err := rs.Row().Int("@err")
					playerId := rs.Row().Int("@playerId")
					//register
					if (err == 0) {
						this.m_Log.Printf("账号[%d]创建玩家[%d]", accountId, playerId)
					} else {
						this.m_Log.Printf("账号[%d]创建玩家失败", accountId)
						world.SERVER.GetServer().SendMsgByID(this.GetSocketId(), "W_A_DeletePlayer", accountId, playerId)
					}

					//通知玩家`
					pPlayer := this.GetPlayer(accountId)
					if pPlayer != nil {
						pPlayer.SendMsg("CreatePlayer", playerId, err)
					}
				}
			}
		}
	})

	this.RegisterCall("testStruct", func(accountId int, pdata []*SimplePlayerData){
		for i,v := range pdata{
			fmt.Println(i, v)
		}
	})
	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	PLAYER.Init(1)
	this.Actor.Start()
}

func (this *PlayerMgr) GetPlayer(accountId int) *Player{
	this.m_Lock.RLock()
	pPlayer, exist := this.m_PlayerMap[accountId]
	this.m_Lock.RUnlock()
	if exist{
		return pPlayer
	}
	return nil
}

func (this *PlayerMgr) AddPlayer(accountId int) *Player{
	LoadPlayerDB := func(accountId int) ([]int, int){
		PlayerList := make([]int, 0)
		PlayerNum := 0
		rows, err := this.m_db.Query(fmt.Sprintf("select player_id from tbl_player where account_id=%d", accountId))
		rs := db.Query(rows)
		if err == nil{
			for rs.Next(){
				PlayerId := rs.Row().Int("player_id")
				PlayerList = append(PlayerList, PlayerId)
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
	this.m_Lock.Lock()
	this.m_PlayerMap[accountId] = pPlayer
	this.m_Lock.Unlock()
	pPlayer.Init(1)
	return pPlayer
}

func (this *PlayerMgr) RemovePlayer(accountId int){
	this.m_Log.Printf("移除帐号数据[%d]", accountId)
	this.m_Lock.Lock()
	delete(this.m_PlayerMap, accountId)
	this.m_Lock.Unlock()
}

func (this *PlayerMgr) PacketFunc(id int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("PlayerMgr PacketFunc", err)
		}
	}()

	SendToPlayer := func(AccountId int, io actor.CallIO) {
		pPlayer := this.GetPlayer(AccountId)
		if pPlayer != nil{
			go pPlayer.Send(io)
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
			nType := bitstream.ReadInt(base.Bit8)
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

//--------------发送给客户端----------------------//
func SendToClient(AccountId int, packet proto.Message){
	pPlayer := PLAYERMGR.GetPlayer(AccountId)
	if pPlayer != nil{
		 world.SendToClient(pPlayer.SocketId, packet)
	}
}