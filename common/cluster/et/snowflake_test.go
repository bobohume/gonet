package et_test

import (
	"context"
	"fmt"
	"gonet/base"
	"gonet/common/cluster/et"
	"log"
	"testing"
	"time"

	"go.etcd.io/etcd/client"
)

type SnowflakeT struct {
	id      int64
	keysAPI client.KeysAPI
	UUID    base.Snowflake
	status  et.STATUS //状态机
}

const (
	uuid_dir1   = "uuid1/"
	ttl_time1   = time.Minute
	WorkeridMax = 1<<13 - 1 //mac下要调制最大连接数，默认256，最大 1 << 10
)

func (s *SnowflakeT) Key() string {
	return uuid_dir1 + fmt.Sprintf("%d", s.id)
}

func (s *SnowflakeT) SET() bool {
	//设置key
	key := s.Key()
	_, err := s.keysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: ttl_time1, PrevExist: client.PrevNoExist, NoValueOnSuccess: true,
	})
	if err != nil {
		s.id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
		return false
	}

	base.UUID.Init(s.id) //设置uuid
	s.status = et.TTL
	return true
}

func (s *SnowflakeT) TTL() {
	//保持ttl
	_, err := s.keysAPI.Set(context.Background(), s.Key(), "", &client.SetOptions{
		TTL: ttl_time1, Refresh: true, NoValueOnSuccess: true,
	})
	if err != nil {
		s.status = et.SET
	} else {
		time.Sleep(time.Second * 20)
	}
}

func (s *SnowflakeT) Run() {
	for {
		switch s.status {
		case et.SET:
			s.SET()
		case et.TTL:
			s.TTL()
		}
	}
}

//uuid生成器
func (s *SnowflakeT) Init(endpoints []string) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 30,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	s.id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
	s.keysAPI = client.NewKeysAPI(etcdClient)
	s.Start()
}

func (s *SnowflakeT) Start() {
	go s.Run()
}

func TestSnowFlake(t *testing.T) {
	group := []*SnowflakeT{}
	for i := 0; i < int(1000); i++ {
		v := &SnowflakeT{}
		v.Init([]string{"http://127.0.0.1:2379"})
		group = append(group, v)
	}

	for {
		time.Sleep(time.Second * 1)
	}
}
