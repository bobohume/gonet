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

//publish
type (
	StubMailBoxMap map[int64]*common.StubMailBox
	StubMailBox    struct {
		*common.ClusterInfo
		m_KeysAPI           client.KeysAPI
		m_StubMailBoxMap    [rpc.STUB_END]StubMailBoxMap
		m_StubMailBoxLocker [rpc.STUB_END]*sync.RWMutex
	}
)

//初始化pub
func (this *StubMailBox) Init(endpoints []string, info *common.ClusterInfo) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	this.ClusterInfo = info
	this.m_KeysAPI = client.NewKeysAPI(etcdClient)
	for i := 0; i < int(rpc.STUB_END); i++ {
		this.m_StubMailBoxLocker[i] = &sync.RWMutex{}
		this.m_StubMailBoxMap[i] = make(StubMailBoxMap)
	}
	this.Start()
	this.getAll()
}

func (this *StubMailBox) Start() {
	go this.Run()
}

func (this *StubMailBox) Publish(info *common.StubMailBox) bool {
	key := fmt.Sprintf("%s%s", STUB_DIR, info.Key())
	data, _ := json.Marshal(info)
	_, err := this.m_KeysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: STUB_TTL_TIME, PrevExist: client.PrevNoExist,
	})
	return err == nil
}

func (this *StubMailBox) Lease(info *common.StubMailBox) error {
	key := fmt.Sprintf("%s%s", STUB_DIR, info.Key())
	_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: STUB_TTL_TIME, Refresh: true, NoValueOnSuccess: true,
	})
	return err
}

func (this *StubMailBox) add(info *common.StubMailBox) {
	this.m_StubMailBoxLocker[info.StubType].Lock()
	pStubMailBox, bOk := this.m_StubMailBoxMap[info.StubType][info.Id]
	if !bOk {
		this.m_StubMailBoxMap[info.StubType][info.Id] = info
	} else {
		*pStubMailBox = *info
	}
	this.m_StubMailBoxLocker[info.StubType].Unlock()
}

func (this *StubMailBox) del(info *common.StubMailBox) {
	this.m_StubMailBoxLocker[info.StubType].Lock()
	delete(this.m_StubMailBoxMap[info.StubType], info.Id)
	this.m_StubMailBoxLocker[info.StubType].Unlock()
}

func (this *StubMailBox) Get(stubType rpc.STUB, Id int64) *common.StubMailBox {
	this.m_StubMailBoxLocker[stubType].RLock()
	pStubMailBox, bEx := this.m_StubMailBoxMap[stubType][Id]
	this.m_StubMailBoxLocker[stubType].RUnlock()
	if bEx {
		return pStubMailBox
	}
	return nil
}

func (this *StubMailBox) Count(stubType rpc.STUB) int {
	this.m_StubMailBoxLocker[stubType].RLock()
	nLen := len(this.m_StubMailBoxMap[stubType])
	this.m_StubMailBoxLocker[stubType].RUnlock()
	return nLen
}

// subscribe
func (this *StubMailBox) Run() {
	watcher := this.m_KeysAPI.Watcher(STUB_DIR, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" || res.Action == "delete" {
			info := nodeToStubMailBox([]byte(res.PrevNode.Value))
			this.del(info)
		} else if res.Action == "set" || res.Action == "create" {
			info := nodeToStubMailBox([]byte(res.Node.Value))
			this.add(info)
		}
	}
}

func (this *StubMailBox) getAll() {
	resp, err := this.m_KeysAPI.Get(context.Background(), STUB_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			for _, v1 := range v.Nodes {
				info := nodeToStubMailBox([]byte(v1.Value))
				this.add(info)
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
