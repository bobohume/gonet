package cluster

import (
	"encoding/json"
	"gonet/server/common"
	"log"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

const(
	ETCD_DIR =  "service/"
)

type Service struct {
	*common.ClusterInfo
	m_KeysAPI client.KeysAPI
}

func (this *Service) Ping(){
	for {
		key := ETCD_DIR + this.String() + "/" + this.IpString()
		data, _ := json.Marshal(this.ClusterInfo)
		this.m_KeysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
			TTL: time.Second * 10,
		})
		time.Sleep(time.Second * 3)
	}
}

func (this *Service) Init(Type int, IP string, Port int, endpoints []string){
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.ClusterInfo = &common.ClusterInfo{Type, IP, Port, 0}
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
}

func (this *Service) Start(){
	go this.Ping()
}

//注册服务器
func NewService(Type int, IP string, Port int, Endpoints []string) *Service{
	service := &Service{}
	service.Init(Type, IP, Port, Endpoints)
	service.Start()
	return service
}