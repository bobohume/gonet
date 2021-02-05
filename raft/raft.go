package raft

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/common/cluster/etv3"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type(
	Raft struct {
		*raft.Raft
		actor.Actor
		*common.ClusterInfo
		m_Master  *cluster.Master
		m_HashRing	*base.HashRing//hash一致性
		m_ClusterInfoMap map[uint32] *common.ClusterInfo
	}
)

func (this *Raft) RegisterRaftCall(){
	//集群新加member
	this.RegisterCall("Cluster_Add", func(ctx context.Context, info *common.ClusterInfo){
		_, bEx := this.m_ClusterInfoMap[info.Id()]
		if !bEx {
			this.m_HashRing.Add(info.IpString())
			if (this.RaftIp() != info.RaftIp()) && this.IsLeader(){
				this.AddVoter(raft.ServerID(info.IpString()), raft.ServerAddress(info.RaftIp()), 0, -1)
			}
		}
	})

	//集群删除member
	this.RegisterCall("Cluster_Del", func(ctx context.Context, info *common.ClusterInfo){
		delete(this.m_ClusterInfoMap, info.Id())
		this.m_HashRing.Remove(info.IpString())
		if (this.RaftIp() != info.RaftIp()) && this.IsLeader(){
			this.RemoveServer(raft.ServerID(info.IpString()), 0, -1)
		}
	})
}

func (this *Raft) InitRaft(info *common.ClusterInfo, Endpoints []string, fsm raft.FSM) {
	this.ClusterInfo = info
	this.m_Master = cluster.NewMaster(info, Endpoints, &this.Actor)
	this.m_HashRing = base.NewHashRing()
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)

	this.Raft, _ = NewRaft(info.RaftIp(), info.IpString(), "./node", fsm)
	services := (*etv3.Master)(this.m_Master).GetServices()
	services = append(services, this.ClusterInfo)

	var configuration raft.Configuration
	for _, v := range services{
		server := raft.Server{ID:raft.ServerID(v.IpString()), Address: raft.ServerAddress(v.RaftIp()),}
		configuration.Servers = append(configuration.Servers, server)
	}

	this.BootstrapCluster(configuration)
}

func (this *Raft) IsLeader() bool{
	return string(this.Leader()) ==  this.RaftIp()
}

func (this *Raft) GetHashRing(Id int64) (error, uint32) {
	return this.m_HashRing.Get64(Id)
}


func NewRaft(raftAddr, raftId, raftDir string, fsm raft.FSM) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftId)
	// config.HeartbeatTimeout = 1000 * time.Millisecond
	// config.ElectionTimeout = 1000 * time.Millisecond
	// config.CommitTimeout = 1000 * time.Millisecond

	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(raftAddr, addr, 3, 5*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(raftDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}

	rf, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	//var configuration raft.Configuration
	//	server := raft.Server{ID:raft.ServerID(info.IpString()), Address: raft.ServerAddress(info.IpString()),}
	//	configuration.Servers = append(configuration.Servers, server)
	//	rf.BootstrapCluster(configuration)

	return rf, nil
}