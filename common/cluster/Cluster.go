package cluster

import (
	"context"
	"errors"
	"gonet/actor"
	"gonet/base"
	"gonet/base/vector"
	"gonet/common"
	"gonet/common/cluster/etv3"
	"gonet/network"
	"gonet/rpc"
	"reflect"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	MAX_CLUSTER_NUM = int(rpc.SERVICE_DB) + 1
	CALL_TIME_OUT   = 500 * time.Millisecond
)

type (
	HashClusterMap       map[uint32]*common.ClusterInfo
	HashClusterSocketMap map[uint32]*common.ClusterInfo

	Op struct {
		m_mailBoxEndpoints     []string
		m_stubMailBoxEndpoints []string
		m_stub                 common.Stub
	}

	OpOption func(*Op)

	//集群服务器
	Cluster struct {
		actor.Actor
		*Service             //集群注册
		m_ClusterMap         [MAX_CLUSTER_NUM]HashClusterMap
		m_ClusterLocker      [MAX_CLUSTER_NUM]*sync.RWMutex
		m_HashRing           [MAX_CLUSTER_NUM]*base.HashRing //hash一致性
		m_Conn               *nats.Conn
		m_DieChan            chan bool
		m_Master             *Master
		m_ClusterInfoMap     map[uint32]*common.ClusterInfo
		m_PacketFuncList     *vector.Vector //call back
		MailBox              etv3.MailBox
		StubMailBox          etv3.StubMailBox
		Stub                 common.Stub
		m_clusterDelCallBack func()
	}

	ICluster interface {
		actor.IActor
		InitCluster(info *common.ClusterInfo, Endpoints []string, natsUrl string, params ...OpOption)
		RegisterClusterCall() //注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *common.ClusterInfo

		BindPacketFunc(packetFunc network.PacketFunc)
		CallMsg(interface{}, rpc.RpcHead, string, ...interface{}) error //同步给集群特定服务器

		RandomCluster(head rpc.RpcHead) rpc.RpcHead //随机分配
		IsEnoughStub(stub rpc.STUB) bool
	}

	EmptyClusterInfo struct {
		common.ClusterInfo
	}

	CallFunc struct {
		Func       interface{}
		FuncType   reflect.Type
		FuncVal    reflect.Value
		FuncParams string
	}
)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func WithMailBoxEtcd(Endpoints []string) OpOption {
	return func(op *Op) {
		op.m_mailBoxEndpoints = Endpoints
	}
}

func WithStubMailBoxEtcd(Endpoints []string, stub *common.Stub) OpOption {
	return func(op *Op) {
		stub.Init()
		op.m_stubMailBoxEndpoints = Endpoints
		op.m_stub = *stub
	}
}

func (this *EmptyClusterInfo) String() string {
	return ""
}

func (this *Cluster) InitCluster(info *common.ClusterInfo, Endpoints []string, natsUrl string, params ...OpOption) {
	this.Actor.Init()
	for i := 0; i < MAX_CLUSTER_NUM; i++ {
		this.m_ClusterLocker[i] = &sync.RWMutex{}
		this.m_ClusterMap[i] = make(HashClusterMap)
		this.m_HashRing[i] = base.NewHashRing()
	}
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)
	this.m_PacketFuncList = vector.NewVector()

	conn, err := setupNatsConn(
		natsUrl,
		this.m_DieChan,
	)
	if err != nil {
		base.LOG.Fatalln("nats connect error!!!!")
	}
	this.m_Conn = conn

	this.m_Conn.Subscribe(getChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff: msg.Data})
	})

	this.m_Conn.Subscribe(getTopicChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff: msg.Data})
	})

	this.m_Conn.Subscribe(getCallChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff: msg.Data, Reply: msg.Reply})
	})

	op := Op{}
	op.applyOpts(params)
	if len(op.m_mailBoxEndpoints) > 0 {
		this.MailBox.Init(op.m_mailBoxEndpoints, info)
	}
	if len(op.m_stubMailBoxEndpoints) > 0 {
		this.StubMailBox.Init(op.m_stubMailBoxEndpoints, info)
		this.Stub = op.m_stub
	}

	rpc.GCall = reflect.ValueOf(this.call)
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
	//注册服务器
	this.Service = NewService(info, Endpoints)
	this.m_Master = NewMaster(&EmptyClusterInfo{}, Endpoints, &this.Actor)
}

//params[0]:rpc.RpcHead
//params[1]:error
func (this *Cluster) call(parmas ...interface{}) {
	head := *parmas[0].(*rpc.RpcHead)
	reply := head.Reply
	head.Reply = ""
	head.ClusterId = head.SrcClusterId
	if parmas[1] == nil {
		parmas[1] = ""
	} else {
		parmas[1] = parmas[1].(error).Error()
	}
	funcName := ""
	packet := rpc.Marshal(&head, &funcName, parmas[1:]...)
	this.m_Conn.Publish(reply, packet.Buff)
}

func (this *Cluster) AddCluster(info *common.ClusterInfo) {
	this.m_ClusterLocker[info.Type].Lock()
	this.m_ClusterMap[info.Type][info.Id()] = info
	this.m_ClusterLocker[info.Type].Unlock()
	this.m_HashRing[info.Type].Add(info.IpString())
	base.LOG.Printf("服务器[%s:%s:%d]建立连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) DelCluster(info *common.ClusterInfo) {
	this.m_ClusterLocker[info.Type].RLock()
	_, bEx := this.m_ClusterMap[info.Type][info.Id()]
	this.m_ClusterLocker[info.Type].RUnlock()
	if bEx {
		this.m_ClusterLocker[info.Type].Lock()
		delete(this.m_ClusterMap[info.Type], info.Id())
		this.m_ClusterLocker[info.Type].Unlock()
	}

	this.m_HashRing[info.Type].Remove(info.IpString())
	if this.m_clusterDelCallBack != nil {
		this.m_clusterDelCallBack()
	}
	base.LOG.Printf("服务器[%s:%s:%d]断开连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) GetCluster(head rpc.RpcHead) *common.ClusterInfo {
	this.m_ClusterLocker[head.DestServerType].RLock()
	defer this.m_ClusterLocker[head.DestServerType].RUnlock()
	pClient, bEx := this.m_ClusterMap[head.DestServerType][head.ClusterId]
	if bEx {
		return pClient
	}
	return nil
}

func (this *Cluster) BindPacketFunc(callfunc network.PacketFunc) {
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Cluster) HandlePacket(packet rpc.Packet) {
	for _, v := range this.m_PacketFuncList.Values() {
		if v.(network.PacketFunc)(packet) {
			break
		}
	}
}

func (this *Cluster) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SrcClusterId = this.Id()
	this.Send(head, rpc.Marshal(&head, &funcName, params...))
}

func (this *Cluster) Send(head rpc.RpcHead, packet rpc.Packet) {
	switch head.SendType {
	//case rpc.SEND_BALANCE:
	//	_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	//	this.m_Conn.Publish(getRpcChannel(head), packet.Buff)
	case rpc.SEND_POINT:
		if head.ClusterId == 0 && head.DestServerType == rpc.SERVICE_GAME {
			pMailBox := this.MailBox.Get(head.Id)
			if pMailBox != nil {
				head.ClusterId = pMailBox.ClusterId
			}
		} else if head.ClusterId == 0 {
			stubType, bEx := this.Stub.StubRoute[head.ActorName]
			if bEx {
				index := head.Id % int64(this.Stub.StubCount[stubType])
				pStub := this.StubMailBox.Get(stubType, index)
				if pStub != nil {
					head.ClusterId = pStub.ClusterId
				}
			}
		}
		this.m_Conn.Publish(getRpcChannel(head), packet.Buff)
	default:
		this.m_Conn.Publish(getRpcTopicChannel(head), packet.Buff)
	}
}

func (this *Cluster) CallMsg(cb interface{}, head rpc.RpcHead, funcName string, params ...interface{}) error {
	head.SrcClusterId = this.Id()
	packet := rpc.Marshal(&head, &funcName, params...)

	switch head.SendType {
	//case rpc.SEND_BALANCE:
	//	_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	//	this.m_Conn.Publish(getRpcChannel(head), packet.Buff)
	case rpc.SEND_POINT:
		if head.ClusterId == 0 && head.DestServerType == rpc.SERVICE_GAME {
			pMailBox := this.MailBox.Get(head.Id)
			if pMailBox != nil {
				head.ClusterId = pMailBox.ClusterId
			}
		} else if head.ClusterId == 0 {
			stubType, bEx := this.Stub.StubRoute[head.ActorName]
			if bEx {
				index := head.Id % int64(this.Stub.StubCount[stubType])
				pStub := this.StubMailBox.Get(stubType, index)
				if pStub != nil {
					head.ClusterId = pStub.ClusterId
				}
			}
		}
	default:
		base.LOG.Printf("CALL MSG [%s] CAN NOT BOARDCAST", funcName)
		//_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	}

	reply, err := this.m_Conn.Request(getRpcCallChannel(head), packet.Buff, CALL_TIME_OUT)
	if err == nil {
		rpcPacket, _ := rpc.Unmarshal(reply.Data)
		cf := &CallFunc{Func: cb, FuncVal: reflect.ValueOf(cb), FuncType: reflect.TypeOf(cb), FuncParams: reflect.TypeOf(cb).String()}
		f := cf.FuncVal
		k := cf.FuncType
		err, params := rpc.UnmarshalBodyCall(rpcPacket, k)
		if err != nil {
			return err
		}
		iLen := len(params)
		if iLen >= 1 {
			in := make([]reflect.Value, iLen)
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			f.Call(in)
		} else {
			base.LOG.Printf("CallMsg [%s] params at least one context", funcName)
			return errors.New("callmsg params at least one context")
		}
	}
	return err
}

func (this *Cluster) RandomCluster(head rpc.RpcHead) rpc.RpcHead {
	if head.Id == 0 {
		head.Id = int64(uint32(base.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}

func (this *Cluster) IsEnoughStub(stub rpc.STUB) bool {
	return this.StubMailBox.Count(stub) == this.Stub.StubCount[stub]
}

//集群新加member
func (this *Cluster) Cluster_Add(ctx context.Context, info *common.ClusterInfo) {
	_, bEx := this.m_ClusterInfoMap[info.Id()]
	if !bEx {
		this.AddCluster(info)
		this.m_ClusterInfoMap[info.Id()] = info
	}
}

//集群删除member
func (this *Cluster) Cluster_Del(ctx context.Context, info *common.ClusterInfo) {
	delete(this.m_ClusterInfoMap, info.Id())
	this.DelCluster(info)
}

var (
	MGR Cluster
)

//链接断开
/*func (this *Cluster) DISCONNECT(ctx context.Context, ClusterId uint32) {
	pInfo, bEx := this.m_ClusterInfoMap[ClusterId]
	if bEx {
		this.DelCluster(pInfo)
	}
	delete(this.m_ClusterInfoMap, ClusterId)
}*/
