package etv3

import (
	"context"
	"fmt"
	"gonet/base"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	uuid_dir = "uuid/"
	ttl_time = 30 * 60
)

type STATUS uint32

const (
	SET STATUS = iota
	TTL STATUS = iota
)

type Snowflake struct {
	id      int64
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	status  STATUS //状态机
}

func (s *Snowflake) Key() string {
	return uuid_dir + fmt.Sprintf("%d", s.id)
}

func (s *Snowflake) SET() bool {
	//设置key
	key := s.Key()
	tx := s.client.Txn(context.Background())
	//key no exist
	leaseResp, err := s.lease.Grant(context.Background(), ttl_time)
	if err != nil {
		return false
	}
	s.leaseId = leaseResp.ID
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(s.leaseId))).
		Else()
	txnRes, err := tx.Commit()
	if err != nil || !txnRes.Succeeded { //抢锁失败
		s.id = int64(base.RandI(1, int(base.WorkeridMax)))
		return false
	}

	base.UUID.Init(s.id) //设置uuid
	s.status = TTL
	return true
}

func (s *Snowflake) TTL() {
	//保持ttl
	_, err := s.lease.KeepAliveOnce(context.Background(), s.leaseId)
	if err != nil {
		s.status = SET
	} else {
		time.Sleep(ttl_time / 3)
	}
}

func (s *Snowflake) Run() {
	for {
		switch s.status {
		case SET:
			s.SET()
		case TTL:
			s.TTL()
		}
	}
}

// uuid生成器
func (s *Snowflake) Init(endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	s.id = int64(base.RandI(1, int(base.WorkeridMax)))
	s.client = etcdClient
	s.lease = lease
	for !s.SET() {
	}
	s.Start()
}

func (s *Snowflake) Start() {
	go s.Run()
}
