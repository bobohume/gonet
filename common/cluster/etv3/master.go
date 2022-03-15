package etv3

import (
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"gonet/actor"
	"gonet/common"
	"gonet/rpc"
	"log"

	"golang.org/x/net/context"
)

//监控服务器
type Master struct {
	m_Client     *clientv3.Client
	m_Actor      actor.IActor
	common.IClusterInfo
}

//监控服务器
func (this *Master) Init(info common.IClusterInfo, Endpoints []string, pActor actor.IActor) {
	cfg := clientv3.Config{
		Endpoints: Endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	this.m_Client = etcdClient
	this.BindActor(pActor)
	this.Start()
	this.IClusterInfo = info
	this.InitServices()
}

func (this *Master) Start() {
	go this.Run()
}

func (this *Master) BindActor(pActor actor.IActor) {
	this.m_Actor = pActor
}

func (this *Master) addService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster_Add", info)
}

func (this *Master) delService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster_Del", info)
}

func NodeToService(val []byte) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (this *Master) Run() {
	wch := this.m_Client.Watch(context.Background(), ETCD_DIR+this.String(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := NodeToService(v1.Kv.Value)
				this.addService(info)
			} else {
				info := NodeToService(v1.PrevKv.Value)
				this.delService(info)
			}
		}
	}
}

func (this *Master) InitServices() {
	resp, err := this.m_Client.Get(context.Background(), ETCD_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := NodeToService(v.Value)
			this.addService(info)
		}
	}
}
