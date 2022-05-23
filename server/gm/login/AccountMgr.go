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
	"gonet/server/model"
)

type (
	AccountMgr struct {
		actor.Actor

		accountMap       map[int64]*Account
		m_AccountNameMap map[string]*Account
	}

	IAccountMgr interface {
		actor.IActor

		GetAccount(int64) *Account
		AddAccount(int64) *Account
		RemoveAccount(int64)
	}
)

func (a *AccountMgr) Init() {
	a.Actor.Init()
	a.accountMap = make(map[int64]*Account)
	a.m_AccountNameMap = make(map[string]*Account)
	//actor.MGR.RegisterActor(a)
	a.Actor.Start()
}

func (a *AccountMgr) GetAccount(accountId int64) *Account {
	pAccount, exist := a.accountMap[accountId]
	if exist {
		return pAccount
	}
	return nil
}

func loadAccount(row orm.IRow, a *AccountDB) {
	a.AccountId = row.Int64("account_id")
	a.AccountName = row.String("account_name")
	a.LoginIp = row.String("login_ip")
	a.Status = row.Int("status")
	a.LoginTime = row.Time("login_time")
	a.LogoutTime = row.Time("logout_time")
}

func (a *AccountMgr) AddAccount(accountId int64) *Account {
	LoadAccountDB := func(accountId int64) *AccountDB {
		rows, err := orm.DB.Query(fmt.Sprintf("select account_id, account_name, status, login_time, logout_time, login_ip from tbl_account where account_id=%d", accountId))
		rs, err := orm.Query(rows, err)
		if err == nil && rs.Next() {
			accountDb := &AccountDB{}
			accountDb.AccountId = accountId
			loadAccount(rs.Row(), accountDb)
			return accountDb
		}
		return nil
	}

	accountDb := LoadAccountDB(accountId)
	if accountDb != nil {
		account := &Account{}
		account.AccountDB = *accountDb
		account.PlayerSimpleDataList = LoadSimplePlayerDatas(accountId)
		a.accountMap[accountId] = account
		a.m_AccountNameMap[account.AccountName] = account
		return account
	}

	return nil
}

func (a *AccountMgr) RemoveAccount(accountId int64) {
	account := a.GetAccount(accountId)
	if account != nil {
		delete(a.m_AccountNameMap, account.AccountName)
		delete(a.accountMap, account.AccountId)
		base.LOG.Printf("账号[%d]断开链接", account.AccountId)
	}
}

//账号登录处理
func (a *AccountMgr) Account_Login(ctx context.Context, accountName string, accountId int64, socketId uint32, id uint32, key int64) {
	account := a.GetAccount(accountId)
	if account == nil {
		account = a.AddAccount(accountId)
	}
	account.GateSocketId = socketId
	base.LOG.Printf("帐号[%s]返回登录OK", accountName)
	PlayerDataList := make([]*message.PlayerData, len(account.PlayerSimpleDataList))
	for i, v := range account.PlayerSimpleDataList {
		PlayerDataList[i] = &message.PlayerData{PlayerID: v.PlayerId, PlayerName: v.PlayerName, PlayerGold: int32(v.Gold)}
	}

	gm.SendToClient(rpc.RpcHead{ClusterId: id, SocketId: socketId, Id: accountId}, &message.SelectPlayerResponse{PacketHead: message.BuildPacketHead(accountId, rpc.SERVICE_GATE),
		Key:        key,
		PlayerData: PlayerDataList,
		AccountId:  accountId,
	})
}

//account创建玩家反馈
func (a *AccountMgr) CreatePlayerRequest(ctx context.Context, packet *message.CreatePlayerRequest) {
	accountId := a.GetRpcHead(ctx).Id
	//查找账号玩家数量
	error := 0
	rows, err := orm.DB.Query(fmt.Sprintf("select count(player_id) as player_count from tbl_player where account_id = %d", accountId))
	if err == nil {
		rs, err := orm.Query(rows, err)
		if err == nil && rs.Next() {
			player_count := rs.Row().Int("player_count")
			if player_count >= 1 { //创建玩家上限
				base.LOG.Printf("账号[%d]创建玩家数量上限", accountId)
			} else { //创建玩家
				playerId := base.UUID.UUID()
				_, err = orm.DB.Exec(fmt.Sprintf("insert into tbl_player (account_id, player_id, player_name, sex, level, gold, draw_gold)"+
					"values(%d, %d, '%s', %d, 1, 0,	0)", accountId, playerId, packet.PlayerName, packet.Sex))
				if err == nil {
					base.LOG.Printf("账号[%d]创建玩家[%d]", accountId, playerId)
					//登录游戏
					a.AddAccount(accountId)
					a.LoginPlayerRequset(ctx, &message.LoginPlayerRequset{PlayerId: playerId})
				}
			}
		}
	}

	if error == 1 { //创建失败通知accout删除player
		base.LOG.Printf("账号[%d]创建玩家失败", accountId)
	}
}

func (a *AccountMgr) LoginPlayerRequset(ctx context.Context, packet *message.LoginPlayerRequset) {
	head := a.GetRpcHead(ctx)
	accountId := head.Id
	playerId := packet.GetPlayerId()
	account := a.GetAccount(accountId)
	if account != nil {
		if !account.SetPlayerId(playerId) {
			base.LOG.Printf("帐号[%d]登入的玩家[%d]不存在", accountId, playerId)
			return
		}

		pMailBox := cluster.MGR.MailBox.Get(playerId)
		if pMailBox == nil {
			GClusterId := cluster.MGR.RandomCluster(rpc.RpcHead{Id: playerId, DestServerType: rpc.SERVICE_GAME}).ClusterId
			if GClusterId == 0 {
				base.LOG.Println("没有可用的GAME集群")
				return
			}
			cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: GClusterId}, "game<-PlayerMgr.LoginPlayerRequset", playerId, head.SrcClusterId, account.GateSocketId)
		} else {
			cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: pMailBox.ClusterId}, "game<-PlayerMgr.LoginPlayerRequset", playerId, head.SrcClusterId, account.GateSocketId)
		}
	}
}

//账号断开连接
func (a *AccountMgr) OnUnRegister(ctx context.Context) {
	accountId := a.GetRpcHead(ctx).Id
	a.RemoveAccount(accountId)
}

func LoadSimplePlayerDatas(accountId int64) []*model.SimplePlayerData {
	pList := make([]*model.SimplePlayerData, 0)
	playerNum := 0
	data := new(model.SimplePlayerData)
	rows, err := orm.DB.Query(orm.LoadSql(data, orm.WithWhereStr(fmt.Sprintf("account_id=%d", accountId))))
	rs, err := orm.Query(rows, err)
	for rs.Next() {
		loadSimple(rs.Row(), data)
		pList = append(pList, data)
		playerNum++
	}
	return pList
}

func loadSimple(row orm.IRow, s *model.SimplePlayerData) {
	s.PlayerId = row.Int64("player_id")
	s.PlayerName = row.String("player_name")
	s.AccountId = row.Int64("account_id")
	s.Level = row.Int("level")
	s.Sex = row.Int("sex")
	s.Gold = row.Int("gold")
	s.DrawGold = row.Int("draw_gold")
	s.Vip = row.Int("vip")
	s.LastLoginTime = row.Time("last_login_time")
	s.LastLogoutTime = row.Time("last_logout_time")
}

func (a *AccountMgr) OnStubRegister(ctx context.Context) {
	//这里可以是加载db数据
	base.LOG.Println("Stub Login register sucess")
}

func (a *AccountMgr) OnStubUnRegister(ctx context.Context) {
	//lease一致性这里要清理缓存数据了
	base.LOG.Println("Stub Login unregister sucess")
}
