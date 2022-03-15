package et

import (
	"encoding/json"
	"fmt"
	"gonet/actor"
	"gonet/rpc"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

const (
	GM_DIR   = "gm/"
)

//publish
type (
	GmRaft struct {
		m_KeysAPI      client.KeysAPI
		m_PlayerLocker *sync.RWMutex
		m_PlayerMap    map[int64]*rpc.GmClusterInfo
	}
)

//初始化pub
func (this *GmRaft) Init(endpoints []string) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
	this.m_PlayerLocker = &sync.RWMutex{}
	this.m_PlayerMap = map[int64]*rpc.GmClusterInfo{}
	this.Start()
	this.InitPlayers()
}

func (this *GmRaft) Start() {
	go this.Run()
}

func (this *GmRaft) Publish(info *rpc.GmClusterInfo, ttl int64) bool {
	info.LeaseId = int64(info.Id)
	key := GM_DIR + fmt.Sprintf("%d", info.Id)
	data, _ := json.Marshal(info)
	_, err := this.m_KeysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: time.Duration(ttl), PrevExist: client.PrevNoExist,
	})
	return err == nil
}

func (this *GmRaft) Lease(Id int64, ttl int64) error {
	key := GM_DIR + fmt.Sprintf("%d", Id)
	_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: time.Duration(ttl), Refresh: true, NoValueOnSuccess: true,
	})
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
	watcher := this.m_KeysAPI.Watcher(GM_DIR, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" {
			info := nodeToGmCluster([]byte(res.PrevNode.Value))
			this.delPlayer(info)
		} else if res.Action == "set" || res.Action == "create"{
			info := nodeToGmCluster([]byte(res.Node.Value))
			this.addPlayer(info)
		} else if res.Action == "delete" {
			info := nodeToGmCluster([]byte(res.Node.Value))
			this.delPlayer(info)
		}
	}
}

func (this *GmRaft) InitPlayers() {
	resp, err := this.m_KeysAPI.Get(context.Background(), GM_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			info := nodeToGmCluster([]byte(v.Value))
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
