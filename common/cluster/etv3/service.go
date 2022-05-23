package etv3

import (
	"encoding/json"
	"gonet/common"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"

	"golang.org/x/net/context"
)

const (
	ETCD_DIR = "server/"
)

//注册服务器
type Service struct {
	*common.ClusterInfo
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	status  STATUS //状态机
}

func (s *Service) SET() {
	leaseResp, _ := s.lease.Grant(context.Background(), 10)
	s.leaseId = leaseResp.ID
	key := ETCD_DIR + s.String() + "/" + s.IpString()
	data, _ := json.Marshal(s.ClusterInfo)
	s.client.Put(context.Background(), key, string(data), clientv3.WithLease(s.leaseId))
	s.status = TTL
	time.Sleep(time.Second * 3)
}

func (s *Service) TTL() {
	//保持ttl
	_, err := s.lease.KeepAliveOnce(context.Background(), s.leaseId)
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

//注册服务器
func (s *Service) Init(info *common.ClusterInfo, endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	s.client = etcdClient
	s.lease = lease
	s.ClusterInfo = info
	s.Start()
}

func (s *Service) Start() {
	go s.Run()
}
