package et

import (
	"encoding/json"
	"gonet/rpc"
	"log"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

const (
	ETCD_DIR = "server/"
)

// 注册服务器
type (
	Service struct {
		*rpc.ClusterInfo
		keysAPI client.KeysAPI
		status  STATUS //状态机
	}
)

func (s *Service) SET() {
	key := ETCD_DIR + s.ServiceName() + "/" + s.IpString()
	data, _ := json.Marshal(s.ClusterInfo)
	s.keysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: time.Second * 10,
	})
	s.status = TTL
	time.Sleep(time.Second * 3)
}

func (s *Service) TTL() {
	//保持ttl
	key := ETCD_DIR + s.ServiceName() + "/" + s.IpString()
	_, err := s.keysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: time.Second * 10, Refresh: true, NoValueOnSuccess: true,
	})
	if err != nil {
		s.status = SET
	} else {
		time.Sleep(time.Second * 3)
	}
}

func (s *Service) Run() {
	for {
		switch s.status {
		case SET:
			s.SET()
		case TTL:
			s.TTL()
		}
	}
}

// 注册服务器
func (s *Service) Init(info *rpc.ClusterInfo, endpoints []string) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	s.ClusterInfo = info
	s.keysAPI = client.NewKeysAPI(etcdClient)
	s.Start()
}

func (s *Service) Start() {
	go s.Run()
}
