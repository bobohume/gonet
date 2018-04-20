package account

import (
	"time"
)

type (
	AccountDB struct{
		AccountId int `primary`//主键
		AccountName string
		Status int
		LoginTime int64 `datetime`//日期
		LogoutTime int64 `datetime`//日期
		LoginIp string
	}

	Account struct{
		AccountDB
	}

	IAccount interface {
		UpdateAccountLogoutTime()
	}
)

func (this *Account)  UpdateAccountLogoutTime(){
	this.LogoutTime = time.Now().Unix()
	//db
}