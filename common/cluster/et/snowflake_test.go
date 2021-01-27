package et_test

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client"
	"gonet/base"
	"log"
	"testing"
	"time"
)

type SnowflakeT struct {
	m_Id int64
	m_KeysAPI client.KeysAPI
	m_UUID base.Snowflake
}

const(
	uuid_dir1 =  "uuid1/"
	ttl_time1 = time.Minute
	WorkeridMax = 1<<9 -1 //mac下要调制最大连接数，默认256，最大 1 << 10
)

func (this *SnowflakeT) Key() string{
	return uuid_dir1 + fmt.Sprintf("%d", this.m_Id)
}

func (this *SnowflakeT) Run(){
	for {
	TrySET:
		//设置key
		key := this.Key()
		_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
			TTL: ttl_time1, PrevExist:client.PrevNoExist, NoValueOnSuccess:true,
		})
		if err != nil{
			resp, err := this.m_KeysAPI.Get(context.Background(), uuid_dir1, &client.GetOptions{Quorum:true})
			if err == nil && (resp != nil && resp.Node != nil){
				Ids := [base.WorkeridMax+1]bool{}
				for _, v := range resp.Node.Nodes{
					Id := base.Int(v.Key[len(uuid_dir1) + 1:])
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
			this.m_Id = this.m_Id & WorkeridMax
			goto TrySET
		}

		this.m_UUID.Init(this.m_Id)//设置uuid

		//保持ttl
	TryTTL:
		_, err = this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
			TTL: ttl_time1, Refresh:true, NoValueOnSuccess:true,
		})
		if err != nil{
			goto TrySET
		}else{
			time.Sleep(time.Second * 10)
			goto TryTTL
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
	this.m_Id = 0
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