package etv3_test

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"gonet/base"
	"gonet/common/cluster/etv3"
	"log"
	"testing"
	"time"
)

type SnowflakeT struct {
	m_Id int64
	m_Client *clientv3.Client
	m_Lease clientv3.Lease
	m_LeaseId clientv3.LeaseID
	m_UUID base.Snowflake
	m_Stats etv3.STATUS
}

const(
	uuid_dir1 =  "uuid1/"
	ttl_time1 = time.Minute
	WorkeridMax = 1<<13 -1 //mac下要调制最大连接数，默认256，最大 1 << 10
)

func (this *SnowflakeT) Key() string{
	return uuid_dir1 + fmt.Sprintf("%d", this.m_Id)
}

func (this *SnowflakeT) SET() bool{
	//设置key
	key := this.Key()
	tx := this.m_Client.Txn(context.Background())
	//key no exist
	leaseResp,err := this.m_Lease.Grant(context.Background(),60)
	if err != nil{
		return false
	}
	this.m_LeaseId = leaseResp.ID
	tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(this.m_LeaseId))).
		Else()
	txnRes, err := tx.Commit()
	if err != nil || !txnRes.Succeeded{//抢锁失败
		this.m_Id = int64(base.RAND.RandI(1, WorkeridMax ))
		return false
	}

	this.m_UUID.Init(this.m_Id)//设置uuid
	this.m_Stats = etv3.TTL
	return true
}

func (this *SnowflakeT) TTL(){
	//保持ttl
	_, err := this.m_Lease.KeepAliveOnce(context.Background(), this.m_LeaseId)
	if err != nil{
		this.m_Stats = etv3.SET
	}else{
		time.Sleep(time.Second * 20)
	}
}

func (this *SnowflakeT) Run(){
	for {
		switch this.m_Stats {
		case etv3.SET:
			this.SET()
		case etv3.TTL:
			this.TTL()
		}
	}
}

func (this *SnowflakeT) Init(endpoints []string){
	cfg := clientv3.Config{
		Endpoints:               endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.m_Id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
	this.m_Client = etcdClient
	this.m_Lease = lease
	this.Start()
}

func (this *SnowflakeT) Start(){
	go this.Run()
}

func TestSnowFlakeT(t *testing.T){
	group := []*SnowflakeT{}
	for i := 0; i < int(1000); i++{
		v := &SnowflakeT{}
		v.Init([]string{"http://127.0.0.1:2379"})
		group = append(group, v)
	}

	for{
		time.Sleep(time.Second * 1)
	}
}