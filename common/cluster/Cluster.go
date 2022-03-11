package cluster

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"gonet/actor"
	"gonet/base"
	"gonet/base/vector"
	"gonet/common"
	"gonet/network"
	"gonet/rpc"
	"log"
	"reflect"
	"sync"
	"time"
)

const(
	MAX_CLUSTER_NUM = int(rpc.SERVICE_DB) + 1
	CALL_TIME_OUT = 50 * time.Millisecond
)

type(
	HashClusterMap map[uint32] *common.ClusterInfo
	HashClusterSocketMap map[uint32] *common.ClusterInfo

	//集群服务器
	Cluster struct{
		actor.Actor
		*Service //集群注册
		m_ClusterMap [MAX_CLUSTER_NUM]HashClusterMap
		m_ClusterLocker [MAX_CLUSTER_NUM]*sync.RWMutex
		m_HashRing	[MAX_CLUSTER_NUM]*base.HashRing//hash一致性
		m_Conn      *nats.Conn
		m_DieChan	chan bool
		m_Master    *Master
		m_ClusterInfoMap map[uint32] *common.ClusterInfo
		m_PacketFuncList	*vector.Vector//call back
		//share_rpc	int `rpc:"GetRpcHead;UpdateTimer"`
	}

	ICluster interface{
		actor.IActor
		InitCluster(info *common.ClusterInfo, Endpoints []string, natsUrl string)
		RegisterClusterCall()//注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *common.ClusterInfo

		BindPacketFunc(packetFunc network.PacketFunc)
		CallMsg(interface{}, rpc.RpcHead, string, ...interface{}) error//同步给集群特定服务器

		RandomCluster(head rpc.RpcHead)	rpc.RpcHead//随机分配
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

func (this *EmptyClusterInfo) String() string{
	return ""
}

func (this *Cluster) InitCluster(info *common.ClusterInfo, Endpoints []string, natsUrl string) {
	this.Actor.Init()
	for  i := 0; i < MAX_CLUSTER_NUM; i++{
		this.m_ClusterLocker[i] = &sync.RWMutex{}
		this.m_ClusterMap[i] = make(HashClusterMap)
		this.m_HashRing[i] = base.NewHashRing()
	}
	//注册服务器
	this.Service = NewService(info, Endpoints)
	this.m_Master = NewMaster(&EmptyClusterInfo{}, Endpoints, &this.Actor)
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)
	this.m_PacketFuncList = vector.NewVector()

	conn, err := setupNatsConn(
		natsUrl,
		this.m_DieChan,
	)
	if err != nil {
		base.GLOG.Fatalln("nats connect error!!!!")
	}
	this.m_Conn = conn

	this.m_Conn.Subscribe(getChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff:msg.Data})
	})

	this.m_Conn.Subscribe(getTopicChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff:msg.Data})
	})

	this.m_Conn.Subscribe(getCallChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff:msg.Data, Reply:msg.Reply})
	})

	rpc.GCall = reflect.ValueOf(this.call)
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

//params[0]:rpc.RpcHead
//params[1]:error
func (this *Cluster) call(parmas ...interface{}) {
	head := *parmas[0].(*rpc.RpcHead)
	reply := head.Reply
	head.Reply = ""
	head.ClusterId = head.SrcClusterId
	if parmas[1] == nil{
		parmas[1] = ""
	}else{
		parmas[1] = parmas[1].(error).Error()
	}
	packet := rpc.Marshal(head, "", parmas[1:]...)
	this.m_Conn.Publish(reply, packet.Buff)
}

func (this *Cluster) AddCluster(info *common.ClusterInfo){
	this.m_ClusterLocker[info.Type].Lock()
	this.m_ClusterMap[info.Type][info.Id()] = info
	this.m_ClusterLocker[info.Type].Unlock()
	this.m_HashRing[info.Type].Add(info.IpString())
	base.GLOG.Printf("服务器[%s:%s:%d]建立连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) DelCluster(info *common.ClusterInfo){
	this.m_ClusterLocker[info.Type].RLock()
	_, bEx := this.m_ClusterMap[info.Type][info.Id()]
	this.m_ClusterLocker[info.Type].RUnlock()
	if bEx{
		this.m_ClusterLocker[info.Type].Lock()
		delete(this.m_ClusterMap[info.Type], info.Id())
		this.m_ClusterLocker[info.Type].Unlock()
	}

	this.m_HashRing[info.Type].Remove(info.IpString())
	base.GLOG.Printf("服务器[%s:%s:%d]断开连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) GetCluster(head rpc.RpcHead) *common.ClusterInfo {
	this.m_ClusterLocker[head.DestServerType].RLock()
	defer this.m_ClusterLocker[head.DestServerType].RUnlock()
	pClient, bEx := this.m_ClusterMap[head.DestServerType][head.ClusterId]
	if bEx{
		return pClient
	}
	return nil
}


func (this *Cluster) BindPacketFunc(callfunc network.PacketFunc){
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Cluster) HandlePacket(packet rpc.Packet){
	for _,v := range this.m_PacketFuncList.Values() {
		if (v.(network.PacketFunc)(packet)){
			break
		}
	}
}

func (this *Cluster) SendMsg(head rpc.RpcHead, funcName string, params  ...interface{}){
	head.SrcClusterId = this.Id()
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *Cluster) Send(head rpc.RpcHead, packet rpc.Packet){
	switch head.SendType{
	case rpc.SEND_BALANCE:
		_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
		this.m_Conn.Publish(getRpcChannel(head), packet.Buff)
	case rpc.SEND_POINT:
		this.m_Conn.Publish(getRpcChannel(head), packet.Buff)
	default:
		this.m_Conn.Publish(getRpcTopicChannel(head), packet.Buff)
	}
}

func (this *Cluster) CallMsg(cb interface{}, head rpc.RpcHead, funcName string, params  ...interface{})error{
	head.SrcClusterId = this.Id()
	packet := rpc.Marshal(head, funcName, params...)

	switch head.SendType{
	case rpc.SEND_POINT:
	default:
		_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	}

	reply, err := this.m_Conn.Request(getRpcCallChannel(head) ,packet.Buff, CALL_TIME_OUT)
	if err == nil{
		rpcPacket, _ := rpc.Unmarshal(reply.Data)
		cf := &CallFunc{Func:cb, FuncVal:reflect.ValueOf(cb), FuncType:reflect.TypeOf(cb), FuncParams:reflect.TypeOf(cb).String()}
		f := cf.FuncVal
		k := cf.FuncType
		err, params := rpc.UnmarshalBodyCall(rpcPacket, k)
		if err != nil{
			return  err
		}
		iLen := len(params)
		if iLen >= 1{
			in := make([]reflect.Value, iLen)
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			f.Call(in)
		}else{
			log.Printf("CallMsg [%s] params at least one context", funcName)
			return errors.New("callmsg params at least one context")
		}
	}
	return err
}

func (this *Cluster) RandomCluster(head rpc.RpcHead) rpc.RpcHead{
	if head.Id == 0{
		head.Id = int64(uint32(base.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.m_HashRing[head.DestServerType].Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil{
		head.SocketId = pCluster.SocketId
	}
	return head
}


//集群新加member
func (this *Cluster) Cluster_Add(ctx context.Context, info *common.ClusterInfo){
	_, bEx := this.m_ClusterInfoMap[info.Id()]
	if !bEx {
		this.AddCluster(info)
		this.m_ClusterInfoMap[info.Id()] = info
	}
}

//集群删除member
func (this *Cluster) Cluster_Del(ctx context.Context, info *common.ClusterInfo){
	delete(this.m_ClusterInfoMap, info.Id())
	this.DelCluster(info)
}

//链接断开
/*func (this *Cluster) DISCONNECT(ctx context.Context, ClusterId uint32) {
	pInfo, bEx := this.m_ClusterInfoMap[ClusterId]
	if bEx {
		this.DelCluster(pInfo)
	}
	delete(this.m_ClusterInfoMap, ClusterId)
}*/
