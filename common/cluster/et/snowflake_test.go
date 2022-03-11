package et_test

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client"
	"gonet/base"
	"gonet/common/cluster/et"
	"log"
	"testing"
	"time"
)

type SnowflakeT struct {
	m_Id int64
	m_KeysAPI client.KeysAPI
	m_UUID base.Snowflake
	m_Stats et.STATUS//状态机
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
	_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: ttl_time1, PrevExist: client.PrevNoExist, NoValueOnSuccess: true,
	})
	if err != nil {
		this.m_Id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
		return false
	}

	base.UUID.Init(this.m_Id) //设置uuid
	this.m_Stats = et.TTL
	return true
}

func (this *SnowflakeT) TTL(){
	//保持ttl
	_, err := this.m_KeysAPI.Set(context.Background(), this.Key(), "", &client.SetOptions{
		TTL: ttl_time1, Refresh: true, NoValueOnSuccess: true,
	})
	if err != nil {
		this.m_Stats = et.SET
	} else {
		time.Sleep(time.Second * 20)
	}
}

func (this *SnowflakeT) Run(){
	for {
		switch this.m_Stats {
		case et.SET:
			this.SET()
		case et.TTL:
			this.TTL()
		}
	}
}

//uuid生成器
func (this *SnowflakeT) Init(endpoints []string){
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 30,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.m_Id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
	this.Start()
}

func (this *SnowflakeT) Start(){
	go this.Run()
}

func TestSnowFlake(t *testing.T){
	group := []*SnowflakeT{}
	for i := 0; i < int(WorkeridMax); i++{
		v := &SnowflakeT{}
		v.Init([]string{"http://127.0.0.1:2379"})
		group = append(group, v)
	}

	for{
		time.Sleep(time.Second * 1)
	}
}