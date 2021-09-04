package world

import (
	"gonet/actor"
	"gonet/common"
	"gonet/rpc"
	"time"
)

type (
	RaftManager struct {
		actor.Actor
		raft.Raft
		m_Fsm raft.ShardingFsm
	}

	IRaftManager interface {
		actor.IActor
	}
)

func (this *RaftManager) Init() {
	this.Actor.Init()
	//注册到集群
	this.m_Fsm.Init(this)
	this.InitRaft(&common.ClusterInfo{Type: rpc.SERVICE_WORLDSERVER, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)}, CONF.Etcd.Endpoints, &this.m_Fsm)
	this.RegisterTimer(10*1000*1000, this.Update)

	this.Actor.Start()
}

func (this *RaftManager) Add(Id int64) {
	if !this.IsLeader() {
		return
	}

	_, clusterId := this.GetHashRing(Id)
	info := raft.ShardingInfo{Op: "add", Id: Id, ClusterId: clusterId}
	this.Apply(info.ToByte(), 10*time.Microsecond)
}

func (this *RaftManager) Del(Id int64) {
	if !this.IsLeader() {
		return
	}
	info := raft.ShardingInfo{Op: "del", Id: Id, ClusterId: 0}
	this.Apply(info.ToByte(), 10*time.Microsecond)
}

func (this *RaftManager) IsSharding(Id int64) bool {
	clusterId := this.m_Fsm.Get(Id)
	if this.Id() == clusterId {
		return true
	}

	return false
}

func (this *RaftManager) Update() {
	//this.Add(base.UUID.UUID())
}
