package raft

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

type(
	Snapshot struct {
		m_Locker *sync.RWMutex
		DataMap map[string] string
	}

	Fsm struct{
		m_Data Snapshot
	}
)

//-------------------must implent--------------------//
func (this *Snapshot) Persist(sink raft.SnapshotSink) error {
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

func (this *Snapshot) Release() {}

func (this *Fsm) Apply(l *raft.Log) interface{} {
	fmt.Println("apply data:", string(l.Data))
	data := strings.Split(string(l.Data), ",")
	op := data[0]
	this.m_Data.m_Locker.Lock()
	if op == "set" {
		key := data[1]
		value := data[2]
		this.m_Data.DataMap[key] = value
	}
	this.m_Data.m_Locker.Unlock()
	return nil
}

func (this *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &this.m_Data, nil
}

func (this *Fsm) Restore(reader io.ReadCloser) error {
	buf := []byte{}
	io.ReadFull(reader, buf)
	this.m_Data.m_Locker.Lock()
	json.Unmarshal(buf, this.m_Data)
	this.m_Data.m_Locker.Unlock()
	return nil
}

func (this *Fsm) Init() {
	this.m_Data.DataMap = map[string]string{}
	this.m_Data.m_Locker = &sync.RWMutex{}
}
//-------------------must implent--------------------//