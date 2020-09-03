package player

import (
	"database/sql"
	"github.com/golang/protobuf/proto"
	"gonet/actor"
	"gonet/base"
)

type(
	//go文件不能来回闭包，很多模块独立player，这里的IPlayer是player的实现
	IPlayer interface {
		actor.IActor

		GetGateClusterId() uint32//获取网关集群id
		GetZoneClusterId() uint32//获取战斗集群id
		GetPlayerId() int64//获取playerid
		GetAccountId() int64//获取账号id

		GetDB() *sql.DB//获取db
		GetLog() *base.CLog//获取log

		SetKV(key int, value int64)//设置kv
		DelKV(key int)//删除key
		GetKV(key int) int64//获取key

		AddBuff(Orgint int, BuffId int)//添加buff
		RemoveBuff(BuffId int)//删除buff

		AddBuffS(Orgint int, BuffId []int)//批量添加buff
		RemoveBuffS(BuffId []int)//批量删除buff

		SendToZone(funcName string, params  ...interface{})
		SendToClient(packet proto.Message)

		GetItemMgr() IItemMgr
	}
)