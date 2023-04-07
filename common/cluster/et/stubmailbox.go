package et

import (
	"encoding/json"
	"fmt"
	"gonet/common"
	"gonet/rpc"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

const (
	STUB_DIR      = "stubmailbox/"
	STUB_TTL_TIME = 30 * time.Second
)

// publish
type (
	StubMailBoxMap map[int64]*common.StubMailBox
	StubMailBox    struct {
		*common.ClusterInfo
		keysAPI           client.KeysAPI
		stubMailBoxMap    [rpc.STUB_END]StubMailBoxMap
		stubMailBoxLocker [rpc.STUB_END]*sync.RWMutex
	}
)

// 初始化pub
func (s *StubMailBox) Init(endpoints []string, info *common.ClusterInfo) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	s.ClusterInfo = info
	s.keysAPI = client.NewKeysAPI(etcdClient)
	for i := 0; i < int(rpc.STUB_END); i++ {
		s.stubMailBoxLocker[i] = &sync.RWMutex{}
		s.stubMailBoxMap[i] = make(StubMailBoxMap)
	}
	s.Start()
}

func (s *StubMailBox) Start() {
	go s.Run()
}

func (s *StubMailBox) Create(info *common.StubMailBox) bool {
	key := fmt.Sprintf("%s%s", STUB_DIR, info.Key())
	data, _ := json.Marshal(info)
	_, err := s.keysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: STUB_TTL_TIME, PrevExist: client.PrevNoExist,
	})
	return err == nil
}

func (s *StubMailBox) Lease(info *common.StubMailBox) error {
	key := fmt.Sprintf("%s%s", STUB_DIR, info.Key())
	_, err := s.keysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: STUB_TTL_TIME, Refresh: true, NoValueOnSuccess: true,
	})
	return err
}

func (s *StubMailBox) add(info *common.StubMailBox) {
	s.stubMailBoxLocker[info.StubType].Lock()
	stub, bOk := s.stubMailBoxMap[info.StubType][info.Id]
	if !bOk {
		s.stubMailBoxMap[info.StubType][info.Id] = info
	} else {
		*stub = *info
	}
	s.stubMailBoxLocker[info.StubType].Unlock()
}

func (s *StubMailBox) del(info *common.StubMailBox) {
	s.stubMailBoxLocker[info.StubType].Lock()
	delete(s.stubMailBoxMap[info.StubType], info.Id)
	s.stubMailBoxLocker[info.StubType].Unlock()
}

func (s *StubMailBox) Get(stubType rpc.STUB, Id int64) *common.StubMailBox {
	s.stubMailBoxLocker[stubType].RLock()
	stub, bEx := s.stubMailBoxMap[stubType][Id]
	s.stubMailBoxLocker[stubType].RUnlock()
	if bEx {
		return stub
	}
	return nil
}

func (s *StubMailBox) Count(stubType rpc.STUB) int64 {
	s.stubMailBoxLocker[stubType].RLock()
	nLen := len(s.stubMailBoxMap[stubType])
	s.stubMailBoxLocker[stubType].RUnlock()
	return int64(nLen)
}

// subscribe
func (s *StubMailBox) Run() {
	watcher := s.keysAPI.Watcher(STUB_DIR, &client.WatcherOptions{
		Recursive: true,
	})
	s.getAll()

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" || res.Action == "delete" {
			info := nodeToStubMailBox([]byte(res.PrevNode.Value))
			s.del(info)
		} else if res.Action == "set" || res.Action == "create" {
			info := nodeToStubMailBox([]byte(res.Node.Value))
			s.add(info)
		}
	}
}

func (s *StubMailBox) getAll() {
	resp, err := s.keysAPI.Get(context.Background(), STUB_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			for _, v1 := range v.Nodes {
				info := nodeToStubMailBox([]byte(v1.Value))
				s.add(info)
			}
		}
	}
}

func nodeToStubMailBox(val []byte) *common.StubMailBox {
	info := &common.StubMailBox{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}
