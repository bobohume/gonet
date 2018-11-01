package account

import (
	"actor"
	"base"
	"database/sql"
	"db"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"message"
	"server/common"
)

type (
	EventProcess struct {
		actor.Actor
		m_db *sql.DB
	}

	IEventProcess interface {
		actor.IActor
	}
)

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_db = SERVER.GetDB()
	this.RegisterCall("COMMON_RegisterRequest", func(nType int, Ip string, Port int) {
		pServerInfo := new(common.ServerInfo)
		pServerInfo.SocketId = this.GetSocketId()
		pServerInfo.Type = nType
		pServerInfo.Ip = Ip
		pServerInfo.Port = Port

		SERVER.GetServerMgr().SendMsg("CONNECT", nType, Ip, Port)

		switch pServerInfo.Type {
		case int(message.SERVICE_GATESERVER):
			SERVER.GetServer().SendMsgByID(this.GetSocketId(), "COMMON_RegisterResponse")
		case int(message.SERVICE_WORLDSERVER):
			SERVER.GetServer().SendMsgByID(this.GetSocketId(), "COMMON_RegisterResponse")
		}
	})

	//链接断开
	this.RegisterCall("DISCONNECT", func(socketid int) {
		SERVER.GetServerMgr().SendMsg("DISCONNECT", socketid)
	})

	//创建账号
	this.RegisterCall("C_A_RegisterRequest", func(packet *message.C_A_RegisterRequest) {
		accountName := *packet.AccountName
		//password := *packet.Password
		password := "123456"
		socketId := int(*packet.SocketId)
		Error := 1
		var result string
		var accountId int
		rows, err := this.m_db.Query(fmt.Sprintf("call `usp_activeaccount`('%s', '%s')", accountName, password))
		if err == nil {
			rows.Next()
			rows.Next()
			if rows.NextResultSet(){
				rs := db.Query(rows)
				if rs.Next(){
					accountId = rs.Row().Int("@accountId")
					result = rs.Row().String("@result")
					if (result == "0000") {
						SERVER.GetLog().Printf("帐号[%s]创建成功", accountName)
						//登录账号
						SERVER.GetAccountMgr().SendMsg( "Account_Login", accountName, accountId, socketId, this.GetSocketId())
						Error = 0
					}
				}
			}
		}

		if Error != 0 {
			SendToClient(this.GetSocketId(), &message.A_C_RegisterResponse{
				PacketHead: message.BuildPacketHead( accountId, 0),
				Error:      proto.Int32(int32(Error)),
				SocketId:packet.SocketId,
			})
		}
	})

	//登录账号
	this.RegisterCall("C_A_LoginRequest", func(packet *message.C_A_LoginRequest) {
		accountName := *packet.AccountName
		//password := *packet.Password
		password := "123456"
		buildVersion := *packet.BuildNo
		socketId := int(*packet.SocketId)
		error := base.NONE_ERROR

		if base.CVERSION().IsAcceptableBuildVersion(buildVersion) {
			log.Printf("账号[%s]登陆账号服务器", accountName)
			rows, err := this.m_db.Query(fmt.Sprintf("call `usp_login`('%s', '%s')", accountName, password))
			if err == nil{
				rows.Next()
				if(rows.NextResultSet()){//存储过程反馈多个select的时候
					rs := db.Query(rows)
					if rs.Next(){
						accountId := rs.Row().Int("@accountId")
						result := rs.Row().String("@result")
						//register account
						if result == "0001" {
							error = base.ACCOUNT_NOEXIST
						} else if (result == "0000") {
							error = base.NONE_ERROR
							SERVER.GetAccountMgr().SendMsg("Account_Login", accountName, accountId, socketId, this.GetSocketId())
						}
					}
				}
			}
		} else {
			error = base.VERSION_ERROR
			log.Println("版本验证错误 clientVersion=%s,err=%d", buildVersion, error)
		}

		if error != base.NONE_ERROR {
			SendToClient(this.GetSocketId(), &message.A_C_LoginRequest{
				PacketHead:message.BuildPacketHead( 0, 0 ),
				Error:proto.Int32(int32(error)),
				SocketId:packet.SocketId,
				AccountName:packet.AccountName,
			})
		}
	})

	//创建玩家
	this.RegisterCall("W_A_CreatePlayer", func(accountId int, playername string, sex int32) {
		rows, err := this.m_db.Query(fmt.Sprintf("call `usp_createplayer`(%d, '%s')", accountId, playername))
		if err == nil{
			rs := db.Query(rows)
			if rs.Next(){
				err := rs.Row().Int("@err")
				playerId := rs.Row().Int("@playerId")
				if err == 0 && playerId > 0 {
					SERVER.GetServer().SendMsgByID(this.GetSocketId(), "A_W_CreatePlayer", accountId, playerId, playername, sex)
				}
			}
		}
	})

	//删除玩家
	this.RegisterCall("W_A_DeletePlayer", func(accountId int, playerId int) {
		this.m_db.Exec(fmt.Sprintf("update tbl_player set delete_flag = 1 where account_id =%d and player_id=%d", accountId, playerId))
	})

	this.Actor.Start()
}
