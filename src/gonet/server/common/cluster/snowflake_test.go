package cluster_test

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
	m_Ip string
	m_KeysAPI client.KeysAPI
	m_UUID base.Snowflake
}

const(
	uuid_dir1 =  "server/uuid1/"
	ttl_time1 = time.Second * 3
	WorkeridMax = 1<<7 -1 //mac下要调制最大连接数，默认256，最大 1 << 10
)

func (this *SnowflakeT) Key() string{
	return uuid_dir1 + fmt.Sprintf("%d", this.m_Id)
}

func (this *SnowflakeT) Value() string{
	return this.m_Ip
}

func (this *SnowflakeT) Ping(){
	for {
	TrySET:
		//设置key
		key := this.Key()
		_, err := this.m_KeysAPI.Set(context.Background(), key, this.Value(), &client.SetOptions{
			TTL: ttl_time1 * 3, PrevExist:client.PrevNoExist,
		})
		if err != nil{
			this.m_Id++
			this.m_Id = this.m_Id & WorkeridMax
			goto TrySET
		}

		this.m_UUID.Init(this.m_Id)//设置uuid

		//保持ttl
	TryTTL:
		resp, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
			TTL: ttl_time1 * 3, Refresh:true,
		})
		if err != nil || (resp != nil && resp.Node != nil && resp.Node.Value != this.Value()){
			goto TrySET
		}else{
			time.Sleep(ttl_time1)
			goto TryTTL
		}
	}
}

func (this *SnowflakeT) Init(IP string, Port int, endpoints []string){
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

func (this *SnowflakeT) Start(){
	go this.Ping()
}

//注册服务器
func NewSnowflakeT(IP string, Port int, Endpoints []string) *SnowflakeT{
	uuid := &SnowflakeT{}
	uuid.Init(IP, Port, Endpoints)
	uuid.Start()
	return uuid
}


func TestSnowFlake(t *testing.T){
	group := []*SnowflakeT{}
	for i := 0; i < int(WorkeridMax); i++{
		group = append(group, NewSnowflakeT("127.0.0.1", i, []string{"http://127.0.0.1:2379"}))
	}

	time.Sleep(3*time.Second)
	for i, _ := range group{
		go func(i int) {
			for{
				mm1 := []int64{}
				for _, v := range group{
					_, id, _ := base.ParseUUID(v.m_UUID.UUID())
					mm1 = append(mm1, id)
				}
				for i, v := range mm1 {
					for i1, v1 := range mm1 {
						if i != i1 && v == v1 {
							fmt.Println(mm1)
							break
						}
					}
				}
				time.Sleep(time.Nanosecond * 100)
			}
		}(i)
	}
	for{
		time.Sleep(time.Second * 1)
	}
}

