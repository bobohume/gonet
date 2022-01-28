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
		m_ServiceMap map[uint32]*common.ClusterInfo
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

	this.m_ServiceMap = make(map[uint32]*common.ClusterInfo)
	this.m_KeysAPI =  client.NewKeysAPI(etcdClient)
	this.Start()
	this.IClusterInfo = info
}

func (this *Master) Start() {
	go this.Run()
}

func (this *Master) addService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{},"Cluster_Add", info)
	this.m_ServiceMap[info.Id()] = info
}

func (this *Master) delService(info *common.ClusterInfo) {
	delete(this.m_ServiceMap, info.Id())
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
		} else if res.Action == "set" {
			info := NodeToService([]byte(res.Node.Value))
			this.addService(info)
		} else if res.Action == "delete" {
			info := NodeToService([]byte(res.Node.Value))
			this.delService(info)
		}
	}
}