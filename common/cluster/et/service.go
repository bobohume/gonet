package et

import (
	"encoding/json"
	"gonet/common"
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
		m_Stats STATUS//状态机
	}
)

func (this *Service) SET() {
	key := ETCD_DIR + this.String() + "/" + this.IpString()
	data, _ := json.Marshal(this.ClusterInfo)
	this.m_KeysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: time.Second * 10,
	})
	this.m_Stats = TTL
	time.Sleep(time.Second * 3)
}

func (this *Service) TTL() {
	//保持ttl
	key := ETCD_DIR + this.String() + "/" + this.IpString()
	_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: time.Second * 10, Refresh: true, NoValueOnSuccess: true,
	})
	if err != nil {
		this.m_Stats = SET
	} else {
		time.Sleep(time.Second * 3)
	}
}

func (this *Service) Run(){
	for {
		switch this.m_Stats {
		case SET:
			this.SET()
		case TTL:
			this.TTL()
		}
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