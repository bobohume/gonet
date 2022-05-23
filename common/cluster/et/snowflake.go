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
	id      int64
	keysAPI client.KeysAPI
	status  STATUS //状态机
}

func (s *Snowflake) Key() string {
	return uuid_dir + fmt.Sprintf("%d", s.id)
}

func (s *Snowflake) SET() bool {
	//设置key
	key := s.Key()
	_, err := s.keysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: ttl_time, PrevExist: client.PrevNoExist, NoValueOnSuccess: true,
	})
	if err != nil {
		s.id = int64(base.RAND.RandI(1, int(base.WorkeridMax)))
		return false
	}

	base.UUID.Init(s.id) //设置uuid
	s.status = TTL
	return true
}

func (s *Snowflake) TTL() {
	//保持ttl
	_, err := s.keysAPI.Set(context.Background(), s.Key(), "", &client.SetOptions{
		TTL: ttl_time, Refresh: true, NoValueOnSuccess: true,
	})
	if err != nil {
		s.status = SET
	} else {
		time.Sleep(ttl_time / 3)
	}
}

func (s *Snowflake) Run() {
	for {
		switch s.status {
		case SET:
			s.SET()
		case TTL:
			s.TTL()
		}
	}
}

//uuid生成器
func (s *Snowflake) Init(endpoints []string) {
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
	for !s.SET() {
	}
	s.Start()
}

func (s *Snowflake) Start() {
	go s.Run()
}
