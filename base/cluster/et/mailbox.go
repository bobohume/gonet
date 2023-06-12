package et

import (
	"encoding/json"
	"fmt"
	"gonet/actor"
	"gonet/rpc"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/client"
	"golang.org/x/net/context"
)

const (
	MAILBOX_DIR     = "mailbox/"
	MAILBOX_TL_TIME = 20 * 60 * time.Second
)

// publish
type (
	MailBox struct {
		*rpc.ClusterInfo
		keysAPI       client.KeysAPI
		mailBoxLocker *sync.RWMutex
		mailBoxMap    map[int64]*rpc.MailBox
	}
)

// 初始化pub
func (m *MailBox) Init(endpoints []string, info *rpc.ClusterInfo) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	m.ClusterInfo = info
	m.keysAPI = client.NewKeysAPI(etcdClient)
	m.mailBoxLocker = &sync.RWMutex{}
	m.mailBoxMap = map[int64]*rpc.MailBox{}
	m.Start()
}

func (m *MailBox) Start() {
	go m.Run()
}

func (m *MailBox) Create(info *rpc.MailBox) bool {
	info.LeaseId = int64(info.Id)
	key := MAILBOX_DIR + fmt.Sprintf("%d", info.Id)
	data, _ := json.Marshal(info)
	_, err := m.keysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: MAILBOX_TL_TIME, PrevExist: client.PrevNoExist,
	})
	return err == nil
}

func (m *MailBox) Lease(Id int64) error {
	key := MAILBOX_DIR + fmt.Sprintf("%d", Id)
	_, err := m.keysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: MAILBOX_TL_TIME, Refresh: true, NoValueOnSuccess: true,
	})
	return err
}

func (m *MailBox) Delete(Id int64) error {
	key := MAILBOX_DIR + fmt.Sprintf("%d", Id)
	_, err := m.keysAPI.Delete(context.Background(), key, &client.DeleteOptions{})
	return err
}

func (m *MailBox) DeleteAll() error {
	_, err := m.keysAPI.Delete(context.Background(), MAILBOX_DIR, &client.DeleteOptions{Recursive: true})
	return err
}

func (m *MailBox) add(info *rpc.MailBox) {
	m.mailBoxLocker.Lock()
	mail, bOk := m.mailBoxMap[info.Id]
	if !bOk {
		m.mailBoxMap[info.Id] = info
	} else {
		*mail = *info
	}
	m.mailBoxLocker.Unlock()
}

func (m *MailBox) del(info *rpc.MailBox) {
	m.mailBoxLocker.Lock()
	delete(m.mailBoxMap, int64(info.Id))
	m.mailBoxLocker.Unlock()
	actor.MGR.SendMsg(rpc.RpcHead{Id: info.Id}, fmt.Sprintf("%s.OnUnRegister", info.MailType.String()))
}

func (m *MailBox) Get(Id int64) *rpc.MailBox {
	m.mailBoxLocker.RLock()
	mail, bEx := m.mailBoxMap[Id]
	m.mailBoxLocker.RUnlock()
	if bEx {
		return mail
	}
	return nil
}

// subscribe
func (m *MailBox) Run() {
	watcher := m.keysAPI.Watcher(MAILBOX_DIR, &client.WatcherOptions{
		Recursive: true,
	})
	m.getAll()

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" || res.Action == "delete" {
			info := nodeToMailBox([]byte(res.PrevNode.Value))
			m.del(info)
		} else if res.Action == "set" || res.Action == "create" {
			info := nodeToMailBox([]byte(res.Node.Value))
			m.add(info)
		}
	}
}

func (m *MailBox) getAll() {
	resp, err := m.keysAPI.Get(context.Background(), MAILBOX_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			info := nodeToMailBox([]byte(v.Value))
			m.add(info)
		}
	}
}

func nodeToMailBox(val []byte) *rpc.MailBox {
	info := &rpc.MailBox{}
	err := json.Unmarshal([]byte(val), info)
	if err != nil {
		log.Print(err)
	}
	return info
}
