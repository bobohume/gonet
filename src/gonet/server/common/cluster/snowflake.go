package cluster

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client"
	"gonet/base"
	"log"
	"time"
)

const(
	uuid_dir =  "server/uuid/"
	ttl_time = time.Second * 3
)

type Snowflake struct {
	m_Id int64
	m_Ip string
	m_KeysAPI client.KeysAPI
}

func (this *Snowflake) Key() string{
	return uuid_dir + fmt.Sprintf("%d", this.m_Id)
}

func (this *Snowflake) Value() string{
	return this.m_Ip
}

func (this *Snowflake) Ping(){
	for {
TrySET:
		//设置key
		key := this.Key()
		_, err := this.m_KeysAPI.Set(context.Background(), key, this.Value(), &client.SetOptions{
			TTL: ttl_time * 3, PrevExist:client.PrevNoExist,
		})
		if err != nil{
			this.m_Id++
			this.m_Id = this.m_Id & base.WorkeridMax
			goto TrySET
		}

		base.UUID.Init(this.m_Id)//设置uuid

		//保持ttl
TryTTL:
		resp, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
			TTL: ttl_time * 3, Refresh:true,
		})
		if err != nil || resp.Node.Value != this.Value(){
			goto TrySET
		}else{
			time.Sleep(ttl_time)
			goto TryTTL
		}
	}
}

func (this *Snowflake) Init(IP string, Port int, endpoints []string){
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.m_Id = 0
	this.m_Ip = fmt.Sprintf("%s:%d", IP, Port)
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
}

func (this *Snowflake) Start(){
	go this.Ping()
}

//注册服务器
func NewSnowflake(IP string, Port int, Endpoints []string) *Snowflake{
	uuid := &Snowflake{}
	uuid.Init(IP, Port, Endpoints)
	uuid.Start()
	return uuid
}
