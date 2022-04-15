package etv3

import (
	"encoding/json"
	"fmt"
	"gonet/common"
	"gonet/rpc"
	"log"
	"sync"

	"go.etcd.io/etcd/clientv3"

	"golang.org/x/net/context"
)

const (
	STUB_DIR      = "stub/"
	STUB_TTL_TIME = 30
)

//publish
type (
	StubMailBoxMap map[int64]*common.StubMailBox
	StubMailBox    struct {
		*common.ClusterInfo
		m_Client            *clientv3.Client
		m_Lease             clientv3.Lease
		m_StubMailBoxMap    [rpc.STUB_END]StubMailBoxMap
		m_StubMailBoxLocker [rpc.STUB_END]*sync.RWMutex
	}
)

//初始化pub
func (this *StubMailBox) Init(endpoints []string, info *common.ClusterInfo) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.ClusterInfo = info
	lease := clientv3.NewLease(etcdClient)
	this.m_Client = etcdClient
	this.m_Lease = lease
	for i := 0; i < int(rpc.STUB_END); i++ {
		this.m_StubMailBoxLocker[i] = &sync.RWMutex{}
		this.m_StubMailBoxMap[i] = make(StubMailBoxMap)
	}
	this.Start()
	this.getAll()
}

func (this *StubMailBox) Start() {
	go this.Run()
}

func (this *StubMailBox) Publish(info *common.StubMailBox) bool {
	leaseResp, err := this.m_Lease.Grant(context.Background(), STUB_TTL_TIME)
	if err == nil {
		leaseId := leaseResp.ID
		info.LeaseId = int64(leaseId)
		key := fmt.Sprintf("%s%s", STUB_DIR, info.Key())
		data, _ := json.Marshal(info)
		//设置key
		tx := this.m_Client.Txn(context.Background())
		tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
			Then(clientv3.OpPut(key, string(data), clientv3.WithLease(leaseId))).
			Else()
		txnRes, err := tx.Commit()
		return err == nil && txnRes.Succeeded
	}
	return false
}

func (this *StubMailBox) Lease(info *common.StubMailBox) error {
	_, err := this.m_Lease.KeepAliveOnce(context.Background(), clientv3.LeaseID(info.LeaseId))
	return err
}

func (this *StubMailBox) add(info *common.StubMailBox) {
	this.m_StubMailBoxLocker[info.StubType].Lock()
	pStubMailBox, bOk := this.m_StubMailBoxMap[info.StubType][info.Id]
	if !bOk {
		this.m_StubMailBoxMap[info.StubType][info.Id] = info
	} else {
		*pStubMailBox = *info
	}
	this.m_StubMailBoxLocker[info.StubType].Unlock()
}

func (this *StubMailBox) del(info *common.StubMailBox) {
	this.m_StubMailBoxLocker[info.StubType].Lock()
	delete(this.m_StubMailBoxMap[info.StubType], info.Id)
	this.m_StubMailBoxLocker[info.StubType].Unlock()
}

func (this *StubMailBox) Get(stubType rpc.STUB, Id int64) *common.StubMailBox {
	this.m_StubMailBoxLocker[stubType].RLock()
	pStubMailBox, bEx := this.m_StubMailBoxMap[stubType][Id]
	this.m_StubMailBoxLocker[stubType].RUnlock()
	if bEx {
		return pStubMailBox
	}
	return nil
}

func (this *StubMailBox) Count(stubType rpc.STUB) int {
	this.m_StubMailBoxLocker[stubType].RLock()
	nLen := len(this.m_StubMailBoxMap[stubType])
	this.m_StubMailBoxLocker[stubType].RUnlock()
	return nLen
}

// subscribe
func (this *StubMailBox) Run() {
	wch := this.m_Client.Watch(context.Background(), STUB_DIR, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := nodeToStubMailBox(v1.Kv.Value)
				this.add(info)
			} else {
				info := nodeToStubMailBox(v1.PrevKv.Value)
				this.del(info)
			}
		}
	}
}

func (this *StubMailBox) getAll() {
	resp, err := this.m_Client.Get(context.Background(), STUB_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := nodeToStubMailBox(v.Value)
			this.add(info)
		}
	}
}

func nodeToStubMailBox(val []byte) *common.StubMailBox {
	info := &common.StubMailBox{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}
