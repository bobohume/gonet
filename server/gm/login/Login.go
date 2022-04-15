package login

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/common/cluster"
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

func (this *Login) Init() {
	this.Actor.Init()
	actor.MGR.RegisterActor(this)
	this.InitPool(this, reflect.TypeOf(AccountMgr{}), MAX_ACCOUNT_MGR_COUNT)
	this.Stub.InitStub(this.GetName(), rpc.STUB_LOGIN)
	this.Actor.Start()
}

func (this *Login) Stub_On_Register(ctx context.Context, Id int64) {
	//这里可以是加载db数据
	base.LOG.Println("Stub Login register sucess")
}

func (this *Login) Stub_On_UnRegister(ctx context.Context, Id int64) {
	//lease一致性这里要清理缓存数据了
	base.LOG.Println("Stub Login unregister sucess")
}

//登录玩家
func (this *Login) LoginPlayer(accountName string) (int64, error) {

	//查找账号玩家数量
	rows, err := orm.DB.Query(fmt.Sprintf("select player_id from tbl_player where account_name = '%s'", accountName))
	rs, err := orm.Query(rows, err)
	playerId := int64(0)
	if err == nil {
		if !rs.Next() {
			playerId = base.UUID.UUID()
			_, err := orm.DB.Exec(fmt.Sprintf("insert into tbl_player (player_id, player_name, account_name, sex, level, gold, draw_gold)"+
				"values(%d, '%s', '%s', %d, 1, 0,	0)", playerId, "test", accountName, 0))
			if err == nil {
				base.LOG.Printf("创建玩家[%d]", playerId)
			}
		} else {
			playerId = rs.Row().Int64("player_id")
		}
	}

	return playerId, err
}

//登录账号
func (this *Login) LoginAccountRequest(ctx context.Context, packet *message.LoginAccountRequest) {
	accountName := packet.GetAccountName()
	password := packet.GetPassword()
	buildVersion := packet.GetBuildNo()
	socketId := uint32(this.GetRpcHead(ctx).Id)
	key := packet.GetKey()
	nError := base.NONE_ERROR

	if accountName == "" {
		nError = base.ACCOUNT_NOEXIST
	} else if base.VERSION.IsAcceptableBuildVersion(buildVersion) {
		log.Printf("账号[%s]登陆账号服务器", accountName)
		rows, err := orm.DB.Query(fmt.Sprintf("select account_id, password from tbl_account where account_name = '%s'", accountName))
		if err == nil {
			rs, err := orm.Query(rows, err)
			if err == nil && rs.Next() {
				accountId := rs.Row().Int64("account_id")
				passWd := rs.Row().String("password")
				if password == passWd {
					nError = base.NONE_ERROR
					cluster.MGR.SendMsg(rpc.RpcHead{Id: accountId}, "gm<-AccountMgr.Account_Login", accountName, accountId, socketId, this.GetRpcHead(ctx).SrcClusterId, key)
				} else { //密码错误
					nError = base.PASSWORD_ERROR
				}
			} else {
				accountId := base.UUID.UUID()
				//创建账号
				_, err := orm.DB.Exec(fmt.Sprintf("insert into tbl_account (account_name, password, account_id) values('%s', '%s', %d)", accountName, password, accountId))
				if err == nil {
					base.LOG.Printf("帐号[%s]创建成功", accountName)
					cluster.MGR.SendMsg(rpc.RpcHead{Id: accountId}, "gm<-AccountMgr.Account_Login", accountName, accountId, socketId, this.GetRpcHead(ctx).SrcClusterId, key)
				}
			}
		}
	} else {
		nError = base.VERSION_ERROR
		log.Printf("版本验证错误 clientVersion=%s,err=%d", buildVersion, nError)
	}

	if nError != base.NONE_ERROR {
		gm.SendToClient(rpc.RpcHead{ClusterId: this.GetRpcHead(ctx).SrcClusterId, SocketId: socketId}, &message.LoginAccountResponse{
			PacketHead:  message.BuildPacketHead(0, 0),
			Error:       int32(nError),
			AccountName: packet.AccountName,
		})
	}
}
