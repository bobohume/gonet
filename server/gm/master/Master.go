package login

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/base/cluster"
	"gonet/orm"
	"gonet/rpc"
)

type (
	Master struct {
		actor.Actor
		actor.ActorPool
		cluster.Stub
	}

	IMaster interface {
		actor.IActor
	}
)

var (
	MASTER Master
)

func Init() {
	MASTER.Init()
}

func (this *Master) Init() {
	this.Actor.Init()
	actor.MGR.RegisterActor(this)
	this.Stub.InitStub(rpc.STUB_Master)
	this.Actor.Start()
}

func (this *Master) OnStubRegister(ctx context.Context) {
	//这里可以是加载db数据
	base.LOG.Println("Stub Master register sucess")
}

func (this *Master) OnStubUnRegister(ctx context.Context) {
	//lease一致性这里要清理缓存数据了
	base.LOG.Println("Stub Login unregister sucess")
}

// 登录玩家
func (this *Master) LoginPlayer(accountName string) (int64, error) {

	//查找账号玩家数量
	rs, err := orm.Query(fmt.Sprintf("select player_id from tbl_player where account_name = '%s'", accountName))
	playerId := int64(0)
	if err == nil {
		if !rs.Next() {
			playerId = base.UUID.UUID()
			_, err := orm.DB.Exec(fmt.Sprintf("insert into tbl_player (player_id, player_name, account_name, sex, level, gold, draw_gold)"+
				"values(%d, '%s', '%s', %d, 1, 0,	0)", playerId, "test", accountName, 0))
			if err == nil {
				base.LOG.Printf("创建玩家[%d]", playerId)
			}
		} else {
			playerId = rs.Row().Int64("player_id")
		}
	}

	return playerId, err
}
