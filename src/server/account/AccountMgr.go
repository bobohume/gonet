package account

import (
	"actor"
	"database/sql"
	"fmt"
)

type (
	AccountMgr struct {
		actor.Actor

		m_AccountMap map[int] *Account
		m_AccountNameMap map[string] *Account
		m_db *sql.DB
	}

	IAccountMgr interface {
		actor.IActor

		GetAccount(int) *Account
		AddAccount(int) *Account
		RemoveAccount(int)
		KickAccount(int)
	}
)

var (
	ACCOUNTMGR AccountMgr
)

func (this* AccountMgr) Init(num int){
	this.m_db = SERVER.GetDB()
	this.Actor.Init(1000)
	this.m_AccountMap = make(map[int] *Account)
	this.m_AccountNameMap = make(map[string] *Account)
	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	//账号登录处理
	this.RegisterCall("Account_Login", func(caller *actor.Caller, accountName string, accountId int, socketId int) {
		LoginAccount := func(pAccount *Account) {
			if pAccount != nil {
				SERVER.GetLog().Printf("帐号[%s]返回登录OK", accountName)
				SERVER.GetServer().SendMsgByID(caller.SocketId, "A_G_Account_Login", accountId, int(socketId))
			}
		}

		pAccount := this.GetAccount(accountId)
		if pAccount != nil {
			this.RemoveAccount(accountId)
		}

		pAccount = this.AddAccount(accountId)
		LoginAccount(pAccount)
	})

	//账号断开连接
	this.RegisterCall("G_ClientLost", func(caller *actor.Caller, accountId int) {
		SERVER.GetLog().Printf("账号[%d] 断开链接", accountId)
		this.RemoveAccount(accountId)
	})

	this.Actor.Start()
}

func (this *AccountMgr) GetAccount(accountId int) *Account{
	pAccount, exist := this.m_AccountMap[accountId]
	if exist{
		return pAccount
	}
	return nil
}

func (this *AccountMgr) AddAccount(accountId int) *Account{
	LoadAccountDB := func(accountId int) *AccountDB {
		row := this.m_db.QueryRow(fmt.Sprintf("select accountName, status, loginTime, logoutTime, loginIp where accountId=%d", accountId))
		if row != nil{
			pAccountDB := &AccountDB{}
			pAccountDB.AccountId = accountId
			row.Scan(&pAccountDB.AccountName, &pAccountDB.Status, &pAccountDB.LoginTime, &pAccountDB.LogoutTime, &pAccountDB.LoginIp)
			return pAccountDB
		}
		return  nil
	}

	pAccountDB := LoadAccountDB(accountId)
	if pAccountDB != nil{
		pAccount := &Account{}
		pAccount.AccountDB = *pAccountDB
		this.m_AccountMap[accountId] = pAccount
		this.m_AccountNameMap[pAccount.AccountName] = pAccount
		return pAccount
	}

	return nil
}

func (this *AccountMgr) RemoveAccount(accountId int){
	pAccount := this.GetAccount(accountId)
	if pAccount != nil{
		delete(this.m_AccountNameMap, pAccount.AccountName)
		delete(this.m_AccountMap, accountId)
		SERVER.GetLog().Printf("账号[%d]断开链接", accountId)
		SERVER.GetServerMgr().SendMsg(0, "G_ClientLost", accountId)
	}
}
