package et

import (
	"encoding/json"
	"fmt"
	"gonet/actor"
	"gonet/common"
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

//publish
type (
	MailBox struct {
		*common.ClusterInfo
		m_KeysAPI       client.KeysAPI
		m_MailBoxLocker *sync.RWMutex
		m_MailBoxMap    map[int64]*rpc.MailBox
	}
)

//初始化pub
func (this *MailBox) Init(endpoints []string, info *common.ClusterInfo) {
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
	this.m_MailBoxLocker = &sync.RWMutex{}
	this.m_MailBoxMap = map[int64]*rpc.MailBox{}
	this.Start()
	this.getAll()
}

func (this *MailBox) Start() {
	go this.Run()
}

func (this *MailBox) Create(info *rpc.MailBox) bool {
	info.LeaseId = int64(info.Id)
	key := MAILBOX_DIR + fmt.Sprintf("%d", info.Id)
	data, _ := json.Marshal(info)
	_, err := this.m_KeysAPI.Set(context.Background(), key, string(data), &client.SetOptions{
		TTL: MAILBOX_TL_TIME, PrevExist: client.PrevNoExist,
	})
	return err == nil
}

func (this *MailBox) Lease(Id int64) error {
	key := MAILBOX_DIR + fmt.Sprintf("%d", Id)
	_, err := this.m_KeysAPI.Set(context.Background(), key, "", &client.SetOptions{
		TTL: MAILBOX_TL_TIME, Refresh: true, NoValueOnSuccess: true,
	})
	return err
}

func (this *MailBox) Delete(Id int64) error {
	key := MAILBOX_DIR + fmt.Sprintf("%d", Id)
	_, err := this.m_KeysAPI.Delete(context.Background(), key, &client.DeleteOptions{})
	return err
}

func (this *MailBox) DeleteAll() error {
	_, err := this.m_KeysAPI.Delete(context.Background(), MAILBOX_DIR, &client.DeleteOptions{Recursive:true})
	return err
}


func (this *MailBox) add(info *rpc.MailBox) {
	this.m_MailBoxLocker.Lock()
	pMailBox, bOk := this.m_MailBoxMap[info.Id]
	if !bOk {
		this.m_MailBoxMap[info.Id] = info
	} else {
		*pMailBox = *info
	}
	this.m_MailBoxLocker.Unlock()
}

func (this *MailBox) del(info *rpc.MailBox) {
	this.m_MailBoxLocker.Lock()
	delete(this.m_MailBoxMap, int64(info.Id))
	this.m_MailBoxLocker.Unlock()
	actor.MGR.SendMsg(rpc.RpcHead{Id: info.Id}, fmt.Sprintf("%s.On_UnRegister", info.MailType.String()))
}

func (this *MailBox) Get(Id int64) *rpc.MailBox {
	this.m_MailBoxLocker.RLock()
	pMailBox, bEx := this.m_MailBoxMap[Id]
	this.m_MailBoxLocker.RUnlock()
	if bEx {
		return pMailBox
	}
	return nil
}

// subscribe
func (this *MailBox) Run() {
	watcher := this.m_KeysAPI.Watcher(MAILBOX_DIR, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch service:", err)
			continue
		}
		if res.Action == "expire" || res.Action == "delete" {
			info := nodeToMailBox([]byte(res.PrevNode.Value))
			this.del(info)
		} else if res.Action == "set" || res.Action == "create" {
			info := nodeToMailBox([]byte(res.Node.Value))
			this.add(info)
		}
	}
}

func (this *MailBox) getAll() {
	resp, err := this.m_KeysAPI.Get(context.Background(), MAILBOX_DIR, &client.GetOptions{Recursive: true})
	if err == nil && (resp != nil && resp.Node != nil) {
		for _, v := range resp.Node.Nodes {
			info := nodeToMailBox([]byte(v.Value))
			this.add(info)
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
