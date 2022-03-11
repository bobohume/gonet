package etv3

import (
	"encoding/json"
	"fmt"
	"gonet/actor"
	"gonet/rpc"
	"log"
	"sync"

	"go.etcd.io/etcd/clientv3"

	"golang.org/x/net/context"
)

const (
	GM_DIR   = "gm/"
)

//publish
type (
	GmRaft struct {
		m_Client       *clientv3.Client
		m_Lease        clientv3.Lease
		m_PlayerLocker *sync.RWMutex
		m_PlayerMap    map[int64]*rpc.GmClusterInfo
	}
)

//初始化pub
func (this *GmRaft) Init(endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.m_Client = etcdClient
	this.m_Lease = lease
	this.m_PlayerLocker = &sync.RWMutex{}
	this.m_PlayerMap = map[int64]*rpc.GmClusterInfo{}
	this.Start()
	this.InitPlayers()
}

func (this *GmRaft) Start() {
	go this.Run()
}

func (this *GmRaft) Publish(info *rpc.GmClusterInfo, ttl int64) bool {
	leaseResp, err := this.m_Lease.Grant(context.Background(), ttl)
	if err == nil {
		leaseId := leaseResp.ID
		info.LeaseId = int64(leaseId)
		key := GM_DIR + fmt.Sprintf("%d", info.Id)
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

func (this *GmRaft) Lease(leaseId int64) error {
	_, err := this.m_Lease.KeepAliveOnce(context.Background(), clientv3.LeaseID(leaseId))
	return err
}

func (this *GmRaft) addPlayer(info *rpc.GmClusterInfo) {
	this.m_PlayerLocker.Lock()
	pPlayer, bOk := this.m_PlayerMap[info.Id]
	if !bOk {
		this.m_PlayerMap[info.Id] = info
	} else {
		*pPlayer = *info
	}
	this.m_PlayerLocker.Unlock()
}

func (this *GmRaft) delPlayer(info *rpc.GmClusterInfo) {
	this.m_PlayerLocker.Lock()
	delete(this.m_PlayerMap, int64(info.Id))
	this.m_PlayerLocker.Unlock()
	actor.MGR.SendMsg(rpc.RpcHead{Id: info.Id}, "GM_Lease_Expire")
}

func (this *GmRaft) GetPlayer(Id int64) *rpc.GmClusterInfo {
	this.m_PlayerLocker.RLock()
	pPlayer, bEx := this.m_PlayerMap[Id]
	this.m_PlayerLocker.RUnlock()
	if bEx {
		return pPlayer
	}
	return nil
}

// subscribe
func (this *GmRaft) Run() {
	wch := this.m_Client.Watch(context.Background(), GM_DIR, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := nodeToGmCluster(v1.Kv.Value)
				this.addPlayer(info)
			} else {
				info := nodeToGmCluster(v1.PrevKv.Value)
				this.delPlayer(info)
			}
		}
	}
}

func (this *GmRaft) InitPlayers() {
	resp, err := this.m_Client.Get(context.Background(), GM_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := nodeToGmCluster(v.Value)
			this.addPlayer(info)
		}
	}
}

func nodeToGmCluster(val []byte) *rpc.GmClusterInfo {
	info := &rpc.GmClusterInfo{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}
