package et

import (
	"encoding/json"
	"gonet/actor"
	"gonet/rpc"
	"log"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

// 监控服务器
type (
	Master struct {
		keysAPI client.KeysAPI
		rpc.IClusterInfo
	}
)

// 监控服务器
func (m *Master) Init(info rpc.IClusterInfo, Endpoints []string, pActor actor.IActor) {
	cfg := client.Config{
		Endpoints:               Endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	m.keysAPI = client.NewKeysAPI(etcdClient)
	m.Start()
	m.IClusterInfo = info
	m.InitServices()
}

func (m *Master) Start() {
	go m.Run()
}

func (m *Master) addService(info *rpc.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster.Cluster_Add", info)
}

func (m *Master) delService(info *rpc.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster.Cluster_Del", info)
}

func NodeToService(val []byte) *rpc.ClusterInfo {
	info := &rpc.ClusterInfo{}
	err := json.Unmarshal(val, info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (m *Master) Run() {
	watcher := m.keysAPI.Watcher(ETCD_DIR, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" || res.Action == "delete" {
			info := NodeToService([]byte(res.PrevNode.Value))
			m.delService(info)
		} else if res.Action == "set" || res.Action == "create" {
			info := NodeToService([]byte(res.Node.Value))
			m.addService(info)
		}
	}
}

func (m *Master) InitServices() {
	resp, err := m.keysAPI.Get(context.Background(), ETCD_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			if v != nil && v.Nodes != nil {
				for _, v1 := range v.Nodes {
					info := NodeToService([]byte(v1.Value))
					m.addService(info)
				}
			}
		}
	}
}
