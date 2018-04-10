package account

import (
	"actor"
	"base"
	"log"
	"database/sql"
	"fmt"
	"message"
	"github.com/golang/protobuf/proto"
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
	this.RegisterCall("COMMON_RegisterRequest", func(caller *actor.Caller, nType int, GateId int, Ip string, Port int) {
		var ServerInfo stServerInfo
		ServerInfo.SocketId = caller.SocketId
		ServerInfo.Type = nType
		ServerInfo.GateId = GateId
		ServerInfo.Ip = Ip
		ServerInfo.Port = Port

		SERVER.GetServerMgr().SendMsg(caller.SocketId, "CONNECT", nType, GateId, Ip, Port)

		switch ServerInfo.Type {
		case int(message.SERVICE_GATESERVER):
			SERVER.GetServer().SendMsgByID(caller.SocketId, "COMMON_RegisterResponse")
		}
	})

	//链接断开
	this.RegisterCall("DISCONNECT", func(caller *actor.Caller, socketid int) {
		SERVER.GetServerMgr().SendMsg(caller.SocketId, "DISCONNECT", socketid)
	})

	//创建账号
	this.RegisterCall("C_A_RegisterRequest", func(caller *actor.Caller, packet *message.C_A_RegisterRequest) {
		accountName := *packet.AccountName
		//password := *packet.Password
		password := "123456"
		socketId := int(*packet.SocketId)
		tx, _ := this.m_db.Begin()
		Error := 1
		var result string
		var accountId int
		_, err := tx.Exec(fmt.Sprintf("call `usp_activeaccount`('%s', '%s')", accountName, password))
		if err == nil {
			row := tx.QueryRow("select @result, @accountId")
			if row != nil {
				row.Scan(&result, &accountId)
				if (result == "0000") {
					SERVER.GetLog().Printf("帐号[%s]创建成功", accountName)
					//登录账号
					SERVER.GetAccountMgr().SendMsg(caller.SocketId, "Account_Login", accountName, accountId, socketId)
					Error = 0
				}
			}
		}
		tx.Commit()

		if Error != 0 {
			SendToClient(caller.SocketId, &message.A_C_RegisterResponse{
				PacketHead: message.BuildPacketHead( accountId, 0),
				Error:      proto.Int32(int32(Error)),
				SocketId:packet.SocketId,
			})
		}
	})

	//登录账号
	this.RegisterCall("C_A_LoginRequest", func(caller *actor.Caller, packet *message.C_A_LoginRequest) {
		accountName := *packet.AccountName
		//password := *packet.Password
		password := "123456"
		buildVersion := *packet.BuildNo
		socketId := int(*packet.SocketId)
		error := base.NONE_ERROR

		if base.CVERSION().IsAcceptableBuildVersion(buildVersion) {
			log.Printf("账号[%s]登陆账号服务器", accountName)

			tx, _ := this.m_db.Begin()
			_, err := tx.Exec(fmt.Sprintf("call `usp_login`('%s', '%s')", accountName, password))
			if err == nil {
				row := tx.QueryRow("select @result, @accountId")
				if row != nil {
					var result string
					var accountId int
					row.Scan(&result, &accountId)

					//register account
					if result == "0001" {
						error = base.ACCOUNT_NOEXIST
					} else if (result == "0000") {
						error = base.NONE_ERROR
						SERVER.GetAccountMgr().SendMsg(caller.SocketId, "Account_Login", accountName, accountId, socketId)
					}
				}
			}
			tx.Commit()
		} else {
			error = base.VERSION_ERROR
			log.Println("版本验证错误 clientVersion=%s,err=%d", buildVersion, error)
		}

		if error != base.NONE_ERROR {
			SendToClient(caller.SocketId, &message.A_C_LoginRequest{
				PacketHead:message.BuildPacketHead( 0, 0 ),
				Error:proto.Int32(int32(error)),
				SocketId:packet.SocketId,
				AccountName:packet.AccountName,
			})
		}
	})

	//创建玩家
	this.RegisterCall("W_A_CreatePlayer", func(caller *actor.Caller, accountId int, playername string, sex int32) {
		tx, _ := this.m_db.Begin()
		_, err := tx.Exec(fmt.Sprintf("call `usp_createplayer`(%d, '%s')", accountId, playername))
		if err == nil {
			row := tx.QueryRow("select @err, @playerId")
			if row != nil {
				var err int
				var playerId int
				row.Scan(&err, &playerId)

				if err == 0 && playerId > 0 {
					SERVER.GetServer().SendMsgByID(caller.SocketId, "A_W_CreatePlayer", accountId, playerId, playername, sex)
				}
			}
		}
		tx.Commit()
	})

	//删除玩家
	this.RegisterCall("W_A_DeletePlayer", func(caller *actor.Caller, accountId int, playerId int) {
		this.m_db.Exec(fmt.Sprintf("update tbl_player set deleteFlag = 1 where accountId =%d and playerid=%d", accountId, playerId))
	})

	this.Actor.Start()
}
