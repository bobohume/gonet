package login

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/base/cluster"
	"gonet/orm"
	"gonet/rpc"
	"gonet/server/gm"
	"gonet/server/message"
	"log"
	"reflect"
)

type (
	Login struct {
		actor.Actor
		actor.ActorPool
		cluster.Stub
	}

	ILogin interface {
		actor.IActor
	}
)

var (
	LOGIN Login
)

const (
	MAX_ACCOUNT_MGR_COUNT = 3
)

func ToSlat(accountName string, pwd string) string {
	return fmt.Sprintf("%s__%s", accountName, pwd)
}

func ToCrc(accountName string, pwd string, buildNo string, nKey int64) uint32 {
	return base.GetMessageCode1(fmt.Sprintf("%s_%s_%s_%d", accountName, pwd, buildNo, nKey))
}

func Init() {
	LOGIN.Init()
}

func (l *Login) Init() {
	l.Actor.Init()
	actor.MGR.RegisterActor(l)
	l.InitPool(l, reflect.TypeOf(AccountMgr{}), MAX_ACCOUNT_MGR_COUNT)
	l.Stub.InitStub(rpc.STUB_AccountMgr)
	l.Actor.Start()
}

// 登录账号
func (l *AccountMgr) LoginAccountRequest(ctx context.Context, packet *message.LoginAccountRequest, gateSocketId uint32) {
	accountName := packet.GetAccountName()
	password := packet.GetPassword()
	buildVersion := packet.GetBuildNo()
	socketId := gateSocketId
	key := packet.GetKey()
	nError := base.NONE_ERROR

	if accountName == "" {
		nError = base.ACCOUNT_NOEXIST
	} else if base.VERSION.IsAcceptableBuildVersion(buildVersion) {
		log.Printf("账号[%s]登陆账号服务器", accountName)
		rs, err := orm.Query(fmt.Sprintf("select account_id, password from tbl_account where account_name = '%s'", accountName))
		if err == nil {
			if err == nil && rs.Next() {
				accountId := rs.Row().Int64("account_id")
				passWd := rs.Row().String("password")
				if password == passWd {
					nError = base.NONE_ERROR
					cluster.MGR.SendMsg(rpc.RpcHead{Id: accountId}, "gm<-AccountMgr.Account_Login", accountName, accountId, socketId, l.GetRpcHead(ctx).SrcClusterId, key)
				} else { //密码错误
					nError = base.PASSWORD_ERROR
				}
			} else {
				accountId := base.UUID.UUID()
				//创建账号
				_, err := orm.DB.Exec(fmt.Sprintf("insert into tbl_account (account_name, password, account_id) values('%s', '%s', %d)", accountName, password, accountId))
				if err == nil {
					base.LOG.Printf("帐号[%s]创建成功", accountName)
					cluster.MGR.SendMsg(rpc.RpcHead{Id: accountId}, "gm<-AccountMgr.Account_Login", accountName, accountId, socketId, l.GetRpcHead(ctx).SrcClusterId, key)
				}
			}
		}
	} else {
		nError = base.VERSION_ERROR
		log.Printf("版本验证错误 clientVersion=%s,err=%d", buildVersion, nError)
	}

	if nError != base.NONE_ERROR {
		gm.SendToClient(rpc.RpcHead{ClusterId: l.GetRpcHead(ctx).SrcClusterId, SocketId: socketId}, &message.LoginAccountResponse{
			PacketHead:  message.BuildPacketHead(0, 0),
			Error:       int32(nError),
			AccountName: packet.AccountName,
		})
	}
}
