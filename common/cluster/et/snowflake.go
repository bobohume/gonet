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
	ttl_time = time.Minute
)

type Snowflake struct {
	m_Id      int64
	m_KeysAPI client.KeysAPI
}

func (this *Snowflake) Key() string {
	return uuid_dir + fmt.Sprintf("%d", this.m_Id)
}

func (this *Snowflake) Run() {
	for {
	TrySET:
		//设置key
		key := this.Key()
		_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
			TTL: ttl_time, PrevExist: client.PrevNoExist, NoValueOnSuccess: true,
		})
		if err != nil {
			resp, err := this.m_KeysAPI.Get(context.Background(), uuid_dir, &client.GetOptions{Quorum: true})
			if err == nil && (resp != nil && resp.Node != nil) {
				Ids := [base.WorkeridMax + 1]bool{}
				for _, v := range resp.Node.Nodes {
					Id := base.Int(v.Key[len(uuid_dir)+1:])
					Ids[Id] = true
				}

				for i, v := range Ids {
					if v == false {
						this.m_Id = int64(i) & base.WorkeridMax
						goto TrySET
					}
				}
			}
			this.m_Id++
			this.m_Id = this.m_Id & base.WorkeridMax
			goto TrySET
		}

		base.UUID.Init(this.m_Id) //设置uuid

		//保持ttl
	TryTTL:
		_, err = this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
			TTL: ttl_time, Refresh: true, NoValueOnSuccess: true,
		})
		if err != nil {
			goto TrySET
		} else {
			time.Sleep(time.Second * 10)
			goto TryTTL
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
	this.m_Id = 0
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
	this.Start()
}

func (this *Snowflake) Start() {
	go this.Run()
}
