package cluster

import (
	"encoding/json"
	"gonet/actor"
	"gonet/server/common"
	"log"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

//监控服务器
type Master struct {
	m_ServiceMap map[uint32]*common.ClusterInfo
	m_KeysAPI client.KeysAPI
	m_Actor actor.IActor
	m_MasterType int
}

//监控服务器
func NewMaster(Type int, Endpoints []string, pActor actor.IActor) *Master {
	master := &Master{}
	master.Init(Endpoints, pActor)
	master.Start()
	master.m_MasterType = Type
	return master
}

func (this *Master) Init(Endpoints []string, pActor actor.IActor) {
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
}

func (this *Master) Start() {
	go this.WatchService()
}

func (this *Master) BindActor(pActor actor.IActor) {
	this.m_Actor = pActor
}

func (this *Master) AddService(info *common.ClusterInfo) {
	_, bEx := this.m_ServiceMap[info.Id()]
	if !bEx{
		this.m_Actor.SendMsg("Cluster_Add", info)
	}
	this.m_ServiceMap[info.Id()] = info
}

func (this *Master) DelService(info *common.ClusterInfo) {
	delete(this.m_ServiceMap, info.Id())
	this.m_Actor.SendMsg("Cluster_Del", info)
}

func (this *Master) InitService(info *common.ClusterInfo) {
	res, err :=this.m_KeysAPI.Get(context.Background(), "workers/service/", nil)
	if err == nil{
		log.Println(res.Node.Value)
		list := []common.ClusterInfo{}
		json.Unmarshal([]byte(res.Node.Value), list)
		for _, v := range list{
			this.m_Actor.SendMsg("Cluster_Socket_Add", v)
		}
	}
}

func NodeToService(node *client.Node) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal([]byte(node.Value), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (this *Master) WatchService() {
	watcher := this.m_KeysAPI.Watcher(ETCD_DIR + common.ToServiceString(this.m_MasterType), &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			break
		}
		if res.Action == "expire" {
			info := NodeToService(res.PrevNode)
			this.DelService(info)
		} else if res.Action == "set" {
			info := NodeToService(res.Node)
			this.AddService(info)
		} else if res.Action == "delete" {
			info := NodeToService(res.Node)
			this.DelService(info)
		}
	}
}