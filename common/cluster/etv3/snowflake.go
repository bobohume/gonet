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
	m_Id      int64
	m_Client  *clientv3.Client
	m_Lease   clientv3.Lease
	m_LeaseId clientv3.LeaseID
	m_Stats   STATUS //状态机
}

func (this *Snowflake) Key() string {
	return uuid_dir + fmt.Sprintf("%d", this.m_Id)
}

func (this *Snowflake) SET() bool {
	//设置key
	key := this.Key()
	tx := this.m_Client.Txn(context.Background())
	//key no exist
	leaseResp, err := this.m_Lease.Grant(context.Background(), ttl_time)
	if err != nil {
		return false
	}
	this.m_LeaseId = leaseResp.ID
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(this.m_LeaseId))).
		Else()
	txnRes, err := tx.Commit()
	if err != nil || !txnRes.Succeeded { //抢锁失败
		this.m_Id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
		return false
	}

	base.UUID.Init(this.m_Id) //设置uuid
	this.m_Stats = TTL
	return true
}

func (this *Snowflake) TTL() {
	//保持ttl
	_, err := this.m_Lease.KeepAliveOnce(context.Background(), this.m_LeaseId)
	if err != nil {
		this.m_Stats = SET
	} else {
		time.Sleep(ttl_time / 3)
	}
}

func (this *Snowflake) Run() {
	for {
		switch this.m_Stats {
		case SET:
			this.SET()
		case TTL:
			this.TTL()
		}
	}
}

//uuid生成器
func (this *Snowflake) Init(endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.m_Id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
	this.m_Client = etcdClient
	this.m_Lease = lease
	for !this.SET() {
	}
	this.Start()
}

func (this *Snowflake) Start() {
	go this.Run()
}
