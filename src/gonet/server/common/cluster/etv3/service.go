package etv3

import (
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"gonet/server/common"
	"log"
	"time"

	"golang.org/x/net/context"
)

const(
	ETCD_DIR =  "server/"
)

//注册服务器
type Service struct {
	*common.ClusterInfo
	m_Client *clientv3.Client
	m_Lease clientv3.Lease
	m_LeaseId clientv3.LeaseID
}

func (this *Service) Run(){
	for {
		leaseResp, _ := this.m_Lease.Grant(context.Background(),10)
		this.m_LeaseId = leaseResp.ID
		key := ETCD_DIR + this.String() + "/" + this.IpString()
		data, _ := json.Marshal(this.ClusterInfo)
		this.m_Client.Put(context.Background(), key, string(data),clientv3.WithLease(this.m_LeaseId))
		time.Sleep(time.Second * 3)
	}
}

//注册服务器
func (this *Service) Init(info *common.ClusterInfo, endpoints []string){
	cfg := clientv3.Config{
		Endpoints:               endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.m_Client = etcdClient
	this.m_Lease = lease
	this.ClusterInfo = info
	this.Start()
}

func (this *Service) Start(){
	go this.Run()
}