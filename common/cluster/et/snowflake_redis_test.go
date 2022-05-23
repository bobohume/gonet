package et_test

import (
	"fmt"
	"gonet/base"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

type SnowflakeR struct {
	id      int64
	ip      string
	keysAPI redis.Conn
	UUID    base.Snowflake
}

const (
	uuid_dir2 = "uuid/"
	uuid_keys = "uuid/*"
)

func (s *SnowflakeR) Key() string {
	return uuid_dir2 + fmt.Sprintf("%d", s.id)
}

func (s *SnowflakeR) Value() string {
	return s.ip
}

func (s *SnowflakeR) Run() {
	for {
	TrySET:
		//设置key
		key := s.Key()
		val, err := redis.Int(s.keysAPI.Do("setnx", key, s.Value()))
		if err == nil && val == 1 {
			s.keysAPI.Do("expire", key, 10)
		} else {
			val, err := redis.Strings(s.keysAPI.Do("keys", uuid_keys))
			if err == nil {
				Ids := [base.WorkeridMax + 1]bool{}
				for _, v := range val {
					Id := base.Int(v[len(uuid_dir2):])
					Ids[Id] = true
				}

				for i, v := range Ids {
					if v == false {
						s.id = int64(i) & base.WorkeridMax
						goto TrySET
					}
				}
			}
			s.id++
			s.id = s.id & WorkeridMax
			goto TrySET
		}
		s.UUID.Init(s.id) //设置uuid

		//保持ttl
	TryTTL:
		val, err = redis.Int(s.keysAPI.Do("expire", key, 10))
		if err == nil && val == 1 {
			time.Sleep(time.Second * 3)
			goto TryTTL
		} else {
			goto TrySET
		}
	}
}

//uuid生成器
func (s *SnowflakeR) Init(IP string, Port int, endpoints []string) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
	}
	if _, err := c.Do("AUTH", "Gonet1q2w3e4r()"); err != nil {
		//c.Close()
	}
	s.keysAPI = c
	s.Start()
}

func (s *SnowflakeR) Start() {
	go s.Run()
}

func TestSnowFlakeRedis(t *testing.T) {
	group := []*SnowflakeR{}
	for i := 0; i < int(WorkeridMax); i++ {
		v := &SnowflakeR{}
		v.Init("127.0.0.1", i, []string{"http://127.0.0.1:2379"})
		group = append(group, v)
	}

	for {
		time.Sleep(time.Second * 1)
	}
}
