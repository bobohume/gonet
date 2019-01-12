package account

import (
	"time"
)

type (
	AccountDB struct{
		AccountId int64 `sql:"primary;name:account_id"`//主键
		AccountName string `sql:"name:account_name"`
		Status int `sql:"name:status"`
		LoginTime int64 `sql:"datetime;name:login_time""`//日期
		LogoutTime int64 `sql:"datetime;name:logout_time""`//日期
		LoginIp string `sql:"name:login_ip""`
	}

	Account struct{
		AccountDB
	}

	IAccount interface {
		CheckLoginTime() bool
		UpdateAccountLogoutTime()
	}
)

func (this *Account) CheckLoginTime() bool{
	return  false
}

func (this *Account)  UpdateAccountLogoutTime(){
	this.LogoutTime = time.Now().Unix()
	//db
}