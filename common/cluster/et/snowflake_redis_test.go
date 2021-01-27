package et_test

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gonet/base"
	"testing"
	"time"
)

type SnowflakeR struct {
	m_Id int64
	m_Ip string
	m_KeysAPI redis.Conn
	m_UUID base.Snowflake
}

const(
	uuid_dir2 =  "uuid/"
	uuid_keys = "uuid/*"
)

func (this *SnowflakeR) Key() string{
	return uuid_dir2 + fmt.Sprintf("%d", this.m_Id)
}

func (this *SnowflakeR) Value() string{
	return this.m_Ip
}

func (this *SnowflakeR) Run(){
	for {
	TrySET:
		//设置key
		key := this.Key()
		val, err := redis.Int(this.m_KeysAPI.Do("setnx", key, this.Value()))
		if err == nil && val == 1{
			this.m_KeysAPI.Do("expire", key, 10)
		}else{
			val, err := redis.Strings(this.m_KeysAPI.Do("keys", uuid_keys))
			if err == nil{
				Ids := [base.WorkeridMax + 1]bool{}
				for _, v := range val {
					Id := base.Int(v[len(uuid_dir2):])
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
			this.m_Id = this.m_Id & WorkeridMax
			goto TrySET
		}
		this.m_UUID.Init(this.m_Id)//设置uuid

		//保持ttl
	TryTTL:
		val, err = redis.Int(this.m_KeysAPI.Do("expire", key, 10))
		if err == nil && val == 1{
			time.Sleep(time.Second * 3)
			goto TryTTL
		}else{
			goto TrySET
		}
	}
}

//uuid生成器
func (this *SnowflakeR) Init(IP string, Port int, endpoints []string){
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
	}
	if _, err := c.Do("AUTH", "Gonet1q2w3e4r()"); err != nil {
		//c.Close()
	}
	this.m_KeysAPI = c
	this.Start()
}

func (this *SnowflakeR) Start(){
	go this.Run()
}

func TestSnowFlakeRedis(t *testing.T){
	group := []*SnowflakeR{}
	for i := 0; i < int(WorkeridMax); i++{
		v := &SnowflakeR{}
		v.Init("127.0.0.1", i, []string{"http://127.0.0.1:2379"})
		group = append(group, v)
	}

	for{
		time.Sleep(time.Second * 1)
	}
}