package common

import (
	"database/sql"
	"gonet/actor"
	"gonet/base"
)

type(
	//go文件不能来回闭包，很多模块独立player，这里的IPlayer是player的实现
	IPlayer interface {
		actor.IActor

		GetGateSocketId() int
		GetPlayerId() int64
		GetAccountId() int64

		GetDB() *sql.DB
		GetLog() *base.CLog
	}
)
