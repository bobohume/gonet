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

		m_AccountMap     map[int64]*Account
		m_AccountNameMap map[string]*Account
		m_PlayerMap      map[int64]*Account
	}

	IAccountMgr interface {
		actor.IActor

		GetAccount(int64) *Account
		AddAccount(int64) *Account
		RemoveAccount(int64)
	}
)

func (this *AccountMgr) Init() {
	this.Actor.Init()
	this.m_AccountMap = make(map[int64]*Account)
	this.m_AccountNameMap = make(map[string]*Account)
	this.m_PlayerMap = make(map[int64]*Account)
	//actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *AccountMgr) GetAccount(accountId int64) *Account {
	pAccount, exist := this.m_AccountMap[accountId]
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

func (this *AccountMgr) AddAccount(accountId int64) *Account {
	LoadAccountDB := func(accountId int64) *AccountDB {
		rows, err := orm.DB.Query(fmt.Sprintf("select account_id, account_name, status, login_time, logout_time, login_ip from tbl_account where account_id=%d", accountId))
		rs, err := orm.Query(rows, err)
		if err == nil && rs.Next() {
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
		pAccount.PlayerSimpleDataList = LoadSimplePlayerDatas(accountId)
		this.m_AccountMap[accountId] = pAccount
		this.m_AccountNameMap[pAccount.AccountName] = pAccount
		return pAccount
	}

	return nil
}

func (this *AccountMgr) RemoveAccount(playerId int64) {
	pAccount, exist := this.m_PlayerMap[playerId]
	if exist {
		delete(this.m_AccountNameMap, pAccount.AccountName)
		delete(this.m_AccountMap, pAccount.AccountId)
		for _, v := range pAccount.PlayerSimpleDataList {
			delete(this.m_PlayerMap, v.PlayerId)
		}
		base.LOG.Printf("账号[%d]断开链接", pAccount.AccountId)
	}
}

//账号登录处理
func (this *AccountMgr) Account_Login(ctx context.Context, accountName string, accountId int64, socketId uint32, id uint32, key int64) {
	pAccount := this.GetAccount(accountId)
	if pAccount == nil {
		pAccount = this.AddAccount(accountId)
	}
	pAccount.GateSocketId = socketId
	base.LOG.Printf("帐号[%s]返回登录OK", accountName)
	PlayerDataList := make([]*message.PlayerData, len(pAccount.PlayerSimpleDataList))
	for i, v := range pAccount.PlayerSimpleDataList {
		PlayerDataList[i] = &message.PlayerData{PlayerID: v.PlayerId, PlayerName: v.PlayerName, PlayerGold: int32(v.Gold)}
		this.m_PlayerMap[v.PlayerId] = pAccount
	}

	gm.SendToClient(rpc.RpcHead{ClusterId: id, SocketId: socketId, Id: accountId}, &message.SelectPlayerResponse{PacketHead: message.BuildPacketHead(accountId, rpc.SERVICE_GATE),
		Key:        key,
		PlayerData: PlayerDataList,
		AccountId:  accountId,
	})
}

//account创建玩家反馈
func (this *AccountMgr) CreatePlayerRequest(ctx context.Context, packet *message.CreatePlayerRequest) {
	accountId := this.GetRpcHead(ctx).Id
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
					this.AddAccount(accountId)
					this.LoginPlayerRequset(ctx, &message.LoginPlayerRequset{PlayerId: playerId})
				}
			}
		}
	}

	if error == 1 { //创建失败通知accout删除player
		base.LOG.Printf("账号[%d]创建玩家失败", accountId)
	}
}

func (this *AccountMgr) LoginPlayerRequset(ctx context.Context, packet *message.LoginPlayerRequset) {
	head := this.GetRpcHead(ctx)
	accountId := head.Id
	playerId := packet.GetPlayerId()
	pAccount := this.GetAccount(accountId)
	if pAccount != nil {
		if !pAccount.SetPlayerId(playerId) {
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
			cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: GClusterId}, "game<-PlayerMgr.LoginPlayerRequset", playerId, head.SrcClusterId, pAccount.GateSocketId)
		} else {
			cluster.MGR.SendMsg(rpc.RpcHead{ClusterId: pMailBox.ClusterId}, "game<-PlayerMgr.LoginPlayerRequset", playerId, head.SrcClusterId, pAccount.GateSocketId)
		}
	}
}

//账号断开连接
func (this *AccountMgr) Player_On_UnRegister(ctx context.Context) {
	playerId := this.GetRpcHead(ctx).Id
	this.RemoveAccount(playerId)
}

func LoadSimplePlayerDatas(accountId int64) []*model.SimplePlayerData {
	pList := make([]*model.SimplePlayerData, 0)
	nPlayerNum := 0
	pData := new(model.SimplePlayerData)
	rows, err := orm.DB.Query(orm.LoadSql(pData, orm.WithWhereStr(fmt.Sprintf("account_id=%d", accountId))))
	rs, err := orm.Query(rows, err)
	for rs.Next() {
		loadSimple(rs.Row(), pData)
		pList = append(pList, pData)
		nPlayerNum++
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
