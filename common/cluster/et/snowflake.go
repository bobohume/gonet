package et

import (
	"context"
	"fmt"
	"gonet/base"
	"log"
	"time"

	"go.etcd.io/etcd/client"
)

const (
	uuid_dir = "uuid/"
	ttl_time = 30 * 60 * time.Second
)

type STATUS uint32

const (
	SET STATUS = iota
	TTL STATUS = iota
)

type Snowflake struct {
	m_Id      int64
	m_KeysAPI client.KeysAPI
	m_Stats   STATUS //状态机
}

func (this *Snowflake) Key() string {
	return uuid_dir + fmt.Sprintf("%d", this.m_Id)
}

func (this *Snowflake) SET() bool {
	//设置key
	key := this.Key()
	_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: ttl_time, PrevExist: client.PrevNoExist, NoValueOnSuccess: true,
	})
	if err != nil {
		this.m_Id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
		return false
	}

	base.UUID.Init(this.m_Id) //设置uuid
	this.m_Stats = TTL
	return true
}

func (this *Snowflake) TTL() {
	//保持ttl
	_, err := this.m_KeysAPI.Set(context.Background(), this.Key(), "", &client.SetOptions{
		TTL: ttl_time, Refresh: true, NoValueOnSuccess: true,
	})
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
	for !this.SET() {
	}
	this.Start()
}

func (this *Snowflake) Start() {
	go this.Run()
}
