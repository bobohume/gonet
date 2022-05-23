package etv3_test

import (
	"context"
	"fmt"
	"gonet/base"
	"gonet/common/cluster/etv3"
	"log"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type SnowflakeT struct {
	id      int64
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	UUID    base.Snowflake
	status  etv3.STATUS
}

const (
	uuid_dir1   = "uuid1/"
	ttl_time1   = time.Minute
	WorkeridMax = 1<<13 - 1 //mac下要调制最大连接数，默认256，最大 1 << 10
)

func (s *SnowflakeT) Key() string {
	return uuid_dir1 + fmt.Sprintf("%d", s.id)
}

func (s *SnowflakeT) SET() bool {
	//设置key
	key := s.Key()
	tx := s.client.Txn(context.Background())
	//key no exist
	leaseResp, err := s.lease.Grant(context.Background(), 60)
	if err != nil {
		return false
	}
	s.leaseId = leaseResp.ID
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(s.leaseId))).
		Else()
	txnRes, err := tx.Commit()
	if err != nil || !txnRes.Succeeded { //抢锁失败
		s.id = int64(base.RAND.RandI(1, WorkeridMax))
		return false
	}

	s.UUID.Init(s.id) //设置uuid
	s.status = etv3.TTL
	return true
}

func (s *SnowflakeT) TTL() {
	//保持ttl
	_, err := s.lease.KeepAliveOnce(context.Background(), s.leaseId)
	if err != nil {
		s.status = etv3.SET
	} else {
		time.Sleep(time.Second * 20)
	}
}

func (s *SnowflakeT) Run() {
	for {
		switch s.status {
		case etv3.SET:
			s.SET()
		case etv3.TTL:
			s.TTL()
		}
	}
}

func (s *SnowflakeT) Init(endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	s.id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
	s.client = etcdClient
	s.lease = lease
	s.Start()
}

func (s *SnowflakeT) Start() {
	go s.Run()
}

func TestSnowFlakeT(t *testing.T) {
	group := []*SnowflakeT{}
	for i := 0; i < int(1000); i++ {
		v := &SnowflakeT{}
		v.Init([]string{"http://127.0.0.1:2379"})
		group = append(group, v)
	}

	for {
		time.Sleep(time.Second * 1)
	}
}
