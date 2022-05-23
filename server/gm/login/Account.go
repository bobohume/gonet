package login

import (
	"gonet/server/model"
	"time"
)

type (
	AccountDB struct {
		AccountName string `sql:"primary;name:account_name"` //主键
		AccountId   int64  `sql:"name:account_id"`
		Status      int    `sql:"name:status"`
		LoginTime   int64  `sql:"datetime;name:login_time"`  //日期
		LogoutTime  int64  `sql:"datetime;name:logout_time"` //日期
		LoginIp     string `sql:"name:login_ip"`
	}

	Account struct {
		AccountDB
		PlayerSimpleDataList []*model.SimplePlayerData
		PlayerId             int64
		GateSocketId         uint32
	}

	IAccount interface {
		CheckLoginTime() bool
		UpdateAccountLogoutTime()
	}
)

func (a *Account) CheckLoginTime() bool {
	return false
}

func (a *Account) UpdateAccountLogoutTime() {
	a.LogoutTime = time.Now().Unix()
	//db
}

func (a *Account) SetPlayerId(PlayerId int64) bool {
	for i := 0; i < len(a.PlayerSimpleDataList); i++ {
		if a.PlayerSimpleDataList[i].PlayerId == PlayerId {
			a.PlayerId = PlayerId
			return true
		}
	}
	return false
}
