package etv3

import (
	"encoding/json"
	"fmt"
	"gonet/rpc"
	"log"
	"sync"

	"go.etcd.io/etcd/clientv3"

	"golang.org/x/net/context"
)

const (
	PLAYER_DIR   = "player/"
	OFFLINE_TIME = 15 * 60
)

//publish
type (
	PlayerRaft struct {
		m_Client       *clientv3.Client
		m_Lease        clientv3.Lease
		m_PlayerLocker *sync.RWMutex
		m_PlayerMap    map[int64]*rpc.PlayerClusterInfo
	}
)

//初始化pub
func (this *PlayerRaft) Init(endpoints []string) {
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
	this.m_PlayerMap = map[int64]*rpc.PlayerClusterInfo{}
	this.Start()
	this.InitPlayers()
}

func (this *PlayerRaft) Start() {
	go this.Run()
}

func (this *PlayerRaft) Publish(info *rpc.PlayerClusterInfo) bool {
	leaseResp, err := this.m_Lease.Grant(context.Background(), OFFLINE_TIME)
	if err == nil {
		leaseId := leaseResp.ID
		info.LeaseId = int64(leaseId)
		key := PLAYER_DIR + fmt.Sprintf("%d", info.Id)
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

func (this *PlayerRaft) Lease(leaseId int64) error {
	_, err := this.m_Lease.KeepAliveOnce(context.Background(), clientv3.LeaseID(leaseId))
	return err
}

func (this *PlayerRaft) addPlayer(info *rpc.PlayerClusterInfo) {
	this.m_PlayerLocker.Lock()
	pPlayer, bOk := this.m_PlayerMap[info.Id]
	if !bOk {
		this.m_PlayerMap[info.Id] = info
	} else {
		*pPlayer = *info
	}
	this.m_PlayerLocker.Unlock()
}

func (this *PlayerRaft) delPlayer(info *rpc.PlayerClusterInfo) {
	this.m_PlayerLocker.Lock()
	delete(this.m_PlayerMap, int64(info.Id))
	this.m_PlayerLocker.Unlock()
}

func (this *PlayerRaft) GetPlayer(Id int64) *rpc.PlayerClusterInfo {
	this.m_PlayerLocker.RLock()
	pPlayer, bEx := this.m_PlayerMap[Id]
	this.m_PlayerLocker.RUnlock()
	if bEx {
		return pPlayer
	}
	return nil
}

// subscribe
func (this *PlayerRaft) Run() {
	wch := this.m_Client.Watch(context.Background(), PLAYER_DIR, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := NodeToPlayer(v1.Kv.Value)
				this.addPlayer(info)
			} else {
				info := NodeToPlayer(v1.PrevKv.Value)
				this.delPlayer(info)
			}
		}
	}
}

func (this *PlayerRaft) InitPlayers() {
	resp, err := this.m_Client.Get(context.Background(), PLAYER_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := NodeToPlayer(v.Value)
			this.addPlayer(info)
		}
	}
}

func NodeToPlayer(val []byte) *rpc.PlayerClusterInfo {
	info := &rpc.PlayerClusterInfo{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}
