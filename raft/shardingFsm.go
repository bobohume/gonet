package raft

import (
	"fmt"
	"github.com/hashicorp/raft"
	jsoniter "github.com/json-iterator/go"
	"gonet/actor"
	"io"
	"sync"
)

type(
	ShardingInfo struct{
		Op string `json:"Op"`
		Id int64 `json:"Id"`
		ClusterId uint32 `json:"ClusterId"`
	}

	ShardingSnapshot struct {
		m_Locker *sync.RWMutex
		DataMap map[int64] uint32
	}

	ShardingFsm struct{
		ShardingSnapshot
		m_Actor actor.IActor
	}
)

func (this *ShardingInfo) ToByte() []byte{
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data,_ := json.Marshal(this)
	return data
}

func (this *ShardingSnapshot) Persist(sink raft.SnapshotSink) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	this.m_Locker.Lock()
	data, err := json.Marshal(this.DataMap)
	this.m_Locker.Unlock()
	if err != nil {
		return err
	}
	sink.Write(data)
	sink.Close()
	return nil
}

func (this *ShardingSnapshot) Release() {}

func (this *ShardingFsm) Apply(l *raft.Log) interface{} {
	info := &ShardingInfo{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(l.Data, info)
	this.m_Locker.Lock()
	if info.Op == "add" {
		this.DataMap[info.Id] = info.ClusterId
	} else if info.Op == "del" {
		delete(this.DataMap, info.Id)
	}
	this.m_Locker.Unlock()
	fmt.Println(info)
	return nil
}

func (this *ShardingFsm) Snapshot() (raft.FSMSnapshot, error) {
	return &this.ShardingSnapshot, nil
}

func (this *ShardingFsm) Restore(reader io.ReadCloser) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	buf := []byte{}
	io.ReadFull(reader, buf)
	this.m_Locker.Lock()
	json.Unmarshal(buf, this.DataMap)
	this.m_Locker.Unlock()

	return nil
}

func (this *ShardingFsm) Init(pActor actor.IActor) {
	this.DataMap = map[int64]uint32{}
	this.m_Locker = &sync.RWMutex{}
	this.m_Actor = pActor
}

func (this *ShardingFsm) Get(Id int64) uint32{
	this.m_Locker.RLock()
	defer this.m_Locker.RUnlock()
	ClusterId, bEx := this.DataMap[Id]
	if bEx{
		return ClusterId
	}
	return 0
}
