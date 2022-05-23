package etv3

import (
	"encoding/json"
	"gonet/actor"
	"gonet/common"
	"gonet/rpc"
	"log"

	"go.etcd.io/etcd/clientv3"

	"golang.org/x/net/context"
)

//监控服务器
type Master struct {
	client *clientv3.Client
	common.IClusterInfo
}

//监控服务器
func (m *Master) Init(info common.IClusterInfo, Endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: Endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	m.client = etcdClient
	m.Start()
	m.IClusterInfo = info
	m.InitServices()
}

func (m *Master) Start() {
	go m.Run()
}

func (m *Master) addService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster.Cluster_Add", info)
}

func (m *Master) delService(info *common.ClusterInfo) {
	actor.MGR.SendMsg(rpc.RpcHead{}, "Cluster.Cluster_Del", info)
}

func NodeToService(val []byte) *common.ClusterInfo {
	info := &common.ClusterInfo{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}

func (m *Master) Run() {
	wch := m.client.Watch(context.Background(), ETCD_DIR+m.String(), clientv3.WithPrefix(), clientv3.WithPrevKV())
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := NodeToService(v1.Kv.Value)
				m.addService(info)
			} else {
				info := NodeToService(v1.PrevKv.Value)
				m.delService(info)
			}
		}
	}
}

func (m *Master) InitServices() {
	resp, err := m.client.Get(context.Background(), ETCD_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := NodeToService(v.Value)
			m.addService(info)
		}
	}
}
