package et

import (
	"encoding/json"
	"gonet/actor"
	"gonet/rpc"
	"gonet/server/common"
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
		m_Actor actor.IActor
		*common.ClusterInfo
	}

	IMaster interface {
		Start()
	}
)

//监控服务器
func (this *Master) Init(info *common.ClusterInfo, Endpoints []string, pActor actor.IActor) {
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
	this.BindActor(pActor)
	this.Start()
	this.ClusterInfo = info
}

func (this *Master) Start() {
	go this.Run()
}

func (this *Master) BindActor(pActor actor.IActor) {
	this.m_Actor = pActor
}

func (this *Master) AddService(info *common.ClusterInfo) {
	_, bEx := this.m_ServiceMap[info.Id()]
	if !bEx{
		this.m_Actor.SendMsg(rpc.RpcHead{},"Cluster_Add", info)
	}
	this.m_ServiceMap[info.Id()] = info
}

func (this *Master) DelService(info *common.ClusterInfo) {
	delete(this.m_ServiceMap, info.Id())
	this.m_Actor.SendMsg(rpc.RpcHead{},"Cluster_Del", info)
}

func (this *Master) InitService(info *common.ClusterInfo) {
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
	watcher := this.m_KeysAPI.Watcher(ETCD_DIR + this.Type.String(), &client.WatcherOptions{
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
			this.DelService(info)
		} else if res.Action == "set" {
			info := NodeToService([]byte(res.Node.Value))
			this.AddService(info)
		} else if res.Action == "delete" {
			info := NodeToService([]byte(res.Node.Value))
			this.DelService(info)
		}
	}
}