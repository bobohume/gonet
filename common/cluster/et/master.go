package et

import (
	"encoding/json"
	"gonet/actor"
	"gonet/common"
	"gonet/rpc"
	"log"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

//监控服务器
type (
	Master struct {
		m_KeysAPI client.KeysAPI
		common.IClusterInfo
	}
)

//监控服务器
func (this *Master) Init(info common.IClusterInfo, Endpoints []string, pActor actor.IActor) {
	cfg := client.Config{
		Endpoints:               Endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	this.m_KeysAPI =  client.NewKeysAPI(etcdClient)
	this.Start()
	this.IClusterInfo = info
	this.InitServices()
}

func (this *Master) Start() {
	go this.Run()
}

func (this *Master) addService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{},"Cluster_Add", info)
}

func (this *Master) delService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{},"Cluster_Del", info)
}

func NodeToService(val []byte) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal(val, info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (this *Master) Run() {
	watcher := this.m_KeysAPI.Watcher(ETCD_DIR+ this.String(), &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" {
			info := NodeToService([]byte(res.PrevNode.Value))
			this.delService(info)
		} else if res.Action == "set" || res.Action == "create"{
			info := NodeToService([]byte(res.Node.Value))
			this.addService(info)
		} else if res.Action == "delete" {
			info := NodeToService([]byte(res.Node.Value))
			this.delService(info)
		}
	}
}

func (this *Master) InitServices() {
	resp, err := this.m_KeysAPI.Get(context.Background(), ETCD_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			if v != nil && v.Nodes != nil{
				for _, v1 := range v.Nodes{
					info := NodeToService([]byte(v1.Value))
					this.addService(info)
				}
			}
		}
	}
}