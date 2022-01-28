package account

import (
	"context"
	"database/sql"
	"fmt"
	"gonet/actor"
	"gonet/db"
	"gonet/rpc"
)

type (
	AccountMgr struct {
		actor.Actor

		m_AccountMap     map[int64]*Account
		m_AccountNameMap map[string]*Account
		m_db             *sql.DB
	}

	IAccountMgr interface {
		actor.IActor

		GetAccount(int64) *Account
		AddAccount(int64) *Account
		RemoveAccount(int64, bool)
		KickAccount(int64)
	}
)

var (
	ACCOUNTMGR AccountMgr
)

func (this *AccountMgr) Init() {
	this.Actor.Init()
	this.m_db = SERVER.GetDB()
	this.m_AccountMap = make(map[int64]*Account)
	this.m_AccountNameMap = make(map[string]*Account)
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *AccountMgr) GetAccount(accountId int64) *Account {
	pAccount, exist := this.m_AccountMap[accountId]
	if exist {
		return pAccount
	}
	return nil
}

func loadAccount(row db.IRow, a *AccountDB) {
	a.AccountId = row.Int64("account_id")
	a.AccountName = row.String("account_name")
	a.LoginIp = row.String("login_ip")
	a.Status = row.Int("status")
	a.LoginTime = row.Time("login_time")
	a.LogoutTime = row.Time("logout_time")
}

func (this *AccountMgr) AddAccount(accountId int64) *Account {
	LoadAccountDB := func(accountId int64) *AccountDB {
		rows, err := this.m_db.Query(fmt.Sprintf("select account_id, account_name, status, login_time, logout_time, login_ip from tbl_account where account_id=%d", accountId))
		rs := db.Query(rows, err)
		if rs.Next() {
			pAccountDB := &AccountDB{}
			pAccountDB.AccountId = accountId
			loadAccount(rs.Row(), pAccountDB)
			return pAccountDB
		}
		return nil
	}

	pAccountDB := LoadAccountDB(accountId)
	if pAccountDB != nil {
		pAccount := &Account{}
		pAccount.AccountDB = *pAccountDB
		this.m_AccountMap[accountId] = pAccount
		this.m_AccountNameMap[pAccount.AccountName] = pAccount
		return pAccount
	}

	return nil
}

func (this *AccountMgr) RemoveAccount(accountId int64, bLogin bool) {
	pAccount := this.GetAccount(accountId)
	if pAccount != nil {
		delete(this.m_AccountNameMap, pAccount.AccountName)
		delete(this.m_AccountMap, accountId)
		SERVER.GetLog().Printf("账号[%d]断开链接", accountId)
	}
	//假如账号服务器分布式，只要踢出world世界服务器即可
	//这里要登录的时候就同步到踢人world
	if bLogin || pAccount != nil {
		KickWorldPlayer(accountId)
	}
}

func (this *AccountMgr) KickAccount(accountId int64) {

}

//账号登录处理
func (this *AccountMgr) Account_Login(ctx context.Context, accountName string, accountId int64, socketId uint32, id uint32) {
	pPlayer := SERVER.GetPlayerRaft().GetPlayer(accountId)
	if pPlayer == nil {
		info := &rpc.PlayerClusterInfo{}
		info.Id = accountId
		info.WClusterId = SERVER.GetCluster().RandomCluster(rpc.RpcHead{Id: accountId, DestServerType: rpc.SERVICE_WORLDSERVER}).ClusterId
		info.ZClusterId = SERVER.GetCluster().RandomCluster(rpc.RpcHead{Id: accountId, DestServerType: rpc.SERVICE_ZONESERVER}).ClusterId
		if info.WClusterId != 0 {
			if SERVER.GetPlayerRaft().Publish(info) {
				pPlayer = info
			}
		} else {
			SERVER.GetLog().Println("没有可用的集群")
		}
	}

	if pPlayer != nil {
		//踢出其他账号服务器
		this.RemoveAccount(accountId, true)
		pAccount := this.AddAccount(accountId)
		if pAccount != nil {
			SERVER.GetLog().Printf("帐号[%s]返回登录OK", accountName)
			SERVER.GetCluster().SendMsg(rpc.RpcHead{ClusterId: id, DestServerType: rpc.SERVICE_GATESERVER}, "A_G_Account_Login", socketId, *pPlayer)
		}
	}
}

//账号断开连接
func (this *AccountMgr) G_ClientLost(ctx context.Context, accountId int64) {
	SERVER.GetLog().Printf("账号[%d] 断开链接", accountId)
	this.RemoveAccount(accountId, false)
}
