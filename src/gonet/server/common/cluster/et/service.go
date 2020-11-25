package et

import (
	"encoding/json"
	"gonet/server/common"
	"log"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

const(
	ETCD_DIR =  "server/"
)

//注册服务器
type(
	Service struct {
		*common.ClusterInfo
		m_KeysAPI client.KeysAPI
	}
)

func (this *Service) Run(){
	for {
		key := ETCD_DIR + this.String() + "/" + this.IpString()
		data, _ := json.Marshal(this.ClusterInfo)
		this.m_KeysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
			TTL: time.Second * 10,
		})

		time.Sleep(time.Second * 3)
	}
}

//注册服务器
func (this *Service) Init(info *common.ClusterInfo, endpoints []string){
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.ClusterInfo = info
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
	this.Start()
}

func (this *Service) Start(){
	go this.Run()
}