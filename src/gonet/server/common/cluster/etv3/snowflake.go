package etv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"gonet/base"
	"log"
	"time"
)

const(
	uuid_dir =  "uuid/"
	ttl_time = time.Minute
)

type Snowflake struct {
	m_Id int64
	m_Client *clientv3.Client
	m_Lease clientv3.Lease
	m_LeaseId clientv3.LeaseID
	m_UUID base.Snowflake
}

func (this *Snowflake) Key() string{
	return uuid_dir + fmt.Sprintf("%d", this.m_Id)
}

func (this *Snowflake) Run(){
	for {
	TrySET:
		//设置key
		key := this.Key()
		tx := this.m_Client.Txn(context.Background())
		//key no exist
		leaseResp,err := this.m_Lease.Grant(context.Background(),60)
		if err != nil{
			goto TrySET
		}
		this.m_LeaseId = leaseResp.ID
		tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
			Then(clientv3.OpPut(key, "", clientv3.WithLease(this.m_LeaseId))).
			Else()
		txnRes, err := tx.Commit()
		if err != nil || !txnRes.Succeeded{//抢锁失败
			resp, err := this.m_Client.Get(context.Background(), uuid_dir)
			if err == nil && (resp != nil && resp.Kvs != nil){
				Ids := [base.WorkeridMax+1]bool{}
				for _, v := range resp.Kvs{
					Id := base.Int(string(v.Value[len(uuid_dir) + 1:]))
					Ids[Id] = true
				}

				for i, v := range Ids{
					if v == false{
						this.m_Id = int64(i) & base.WorkeridMax
						goto TrySET
					}
				}
			}
			this.m_Id++
			this.m_Id = this.m_Id & base.WorkeridMax
			goto TrySET
		}

		this.m_UUID.Init(this.m_Id)//设置uuid

		//保持ttl
	TryTTL:
		_, err = this.m_Lease.KeepAliveOnce(context.Background(), this.m_LeaseId)
		if err != nil{
			goto TrySET
		}else{
			time.Sleep(time.Second * 10)
			goto TryTTL
		}
	}
}

//uuid生成器
func (this *Snowflake) Init(endpoints []string){
	cfg := clientv3.Config{
		Endpoints:               endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.m_Id = 0
	this.m_Client = etcdClient
	this.m_Lease = lease
	this.Start()
}

func (this *Snowflake) Start(){
	go this.Run()
}