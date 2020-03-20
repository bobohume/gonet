package player

import (
	"database/sql"
	"gonet/actor"
	"gonet/base"
)

type(
	//go文件不能来回闭包，很多模块独立player，这里的IPlayer是player的实现
	IPlayer interface {
		actor.IActor

		GetGateClusterId() uint32//获取网关id
		GetPlayerId() int64//获取playerid
		GetAccountId() int64//获取账号id

		GetDB() *sql.DB//获取db
		GetLog() *base.CLog//获取log

		SetKV(key int, value int64)//设置kv
		DelKV(key int)//删除key
		GetKV(key int) int64//获取key

		GetItemMgr() IItemMgr
	}
)