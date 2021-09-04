package raft

import (
	"gonet/base"
	"gonet/common"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type (
	Raft struct {
		*raft.Raft
		*common.ClusterInfo
		m_HashRing       *base.HashRing //hash一致性
		m_ClusterInfoMap map[uint32]*common.ClusterInfo
	}
)

func (this *Raft) InitRaft(info *common.ClusterInfo, Endpoints []string, fsm raft.FSM) {
	this.ClusterInfo = info
	this.m_HashRing = base.NewHashRing()
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)

	this.Raft, _ = NewRaft(info.IpString(), info.IpString(), "./node", fsm)
	var configuration raft.Configuration
	for _, v := range Endpoints {
		server := raft.Server{ID: raft.ServerID(v), Address: raft.ServerAddress(v)}
		configuration.Servers = append(configuration.Servers, server)
	}

	this.BootstrapCluster(configuration)
}

func (this *Raft) IsLeader() bool {
	return string(this.Leader()) == this.IpString()
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
