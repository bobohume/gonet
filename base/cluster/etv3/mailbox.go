package etv3

import (
	"fmt"
	"gonet/actor"
	"gonet/rpc"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	"go.etcd.io/etcd/clientv3"

	"golang.org/x/net/context"
)

const (
	MAILBOX_DIR     = "mailbox/"
	MAILBOX_TL_TIME = 20 * 60
)

// publish
type (
	MailBox struct {
		*rpc.ClusterInfo
		client        *clientv3.Client
		lease         clientv3.Lease
		mailBoxLocker *sync.RWMutex
		mailBoxMap    map[int64]*rpc.MailBox
	}
)

// 初始化pub
func (m *MailBox) Init(endpoints []string, info *rpc.ClusterInfo) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	m.ClusterInfo = info
	lease := clientv3.NewLease(etcdClient)
	m.client = etcdClient
	m.lease = lease
	m.mailBoxLocker = &sync.RWMutex{}
	m.mailBoxMap = map[int64]*rpc.MailBox{}
	m.Start()
}

func (m *MailBox) Start() {
	go m.Run()
}

func (m *MailBox) Create(info *rpc.MailBox) bool {
	leaseResp, err := m.lease.Grant(context.Background(), MAILBOX_TL_TIME)
	if err == nil {
		leaseId := leaseResp.ID
		info.LeaseId = int64(leaseId)
		key := MAILBOX_DIR + fmt.Sprintf("%d", info.Id)
		data, _ := proto.Marshal(info)
		//设置key
		tx := m.client.Txn(context.Background())
		tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
			Then(clientv3.OpPut(key, string(data), clientv3.WithLease(leaseId))).
			Else()
		txnRes, err := tx.Commit()
		return err == nil && txnRes.Succeeded
	}
	return false
}

func (m *MailBox) Lease(leaseId int64) error {
	_, err := m.lease.KeepAliveOnce(context.Background(), clientv3.LeaseID(leaseId))
	return err
}

func (m *MailBox) Delete(Id int64) error {
	key := MAILBOX_DIR + fmt.Sprintf("%d", Id)
	_, err := m.client.Delete(context.Background(), key)
	return err
}

func (m *MailBox) DeleteAll() error {
	_, err := m.client.Delete(context.Background(), MAILBOX_DIR, clientv3.WithPrefix())
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
	wch := m.client.Watch(context.Background(), MAILBOX_DIR, clientv3.WithPrefix(), clientv3.WithPrevKV())
	m.getAll()
	for v := range wch {
		for _, v1 := range v.Events {
			if v1.Type.String() == "PUT" {
				info := nodeToMailBox(v1.Kv.Value)
				m.add(info)
			} else {
				info := nodeToMailBox(v1.PrevKv.Value)
				m.del(info)
			}
		}
	}
}

func (m *MailBox) getAll() {
	resp, err := m.client.Get(context.Background(), MAILBOX_DIR, clientv3.WithPrefix())
	if err == nil && (resp != nil && resp.Kvs != nil) {
		for _, v := range resp.Kvs {
			info := nodeToMailBox(v.Value)
			m.add(info)
		}
	}
}

func nodeToMailBox(val []byte) *rpc.MailBox {
	mail := &rpc.MailBox{}
	err := proto.Unmarshal([]byte(val), mail)
	if err != nil {
		log.Print(err)
	}
	return mail
}
