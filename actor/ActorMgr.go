package actor

import (
	"gonet/base"
	"gonet/network"
	"gonet/rpc"
	"log"
	"reflect"
	"strings"
	"sync"
)

type ACTOR_TYPE uint32

const (
	ACTOR_TYPE_SINGLETON ACTOR_TYPE = iota //单列
	ACTOR_TYPE_PLAYER    ACTOR_TYPE = iota //玩家 必须初始一个全局的actor 作为类型判断
) //ACTOR_TYPE

const (
    MAX_RPC_TAG = 10
)

//一些全局的actor,不可删除的,不用锁考虑性能
//不是全局的actor,请使用actor pool
type (
	Op struct {
		m_name string //name
		m_type ACTOR_TYPE
		m_RpcMethodMap  map[string] string
	}

	OpOption func(*Op)

	ActorMgr struct {
		m_ActorMap     map[reflect.Type]IActor
		m_ActorNameMap map[string]IActor
	    m_MsgMap       map[string]IActor
		m_RpcMethodMap map[reflect.Type] map[string] string
		m_PlayerMap    map[int64]IActor
		m_PlayerLock   *sync.RWMutex
		m_bStart       bool
	}

	IActorMgr interface {
		Init()
		RegisterActor(pActor IActor, params ...OpOption) //注册回调
		PacketFunc(rpc.Packet) bool //回调函数
		SendMsg(rpc.RpcHead, string, ...interface{})
	}

	ICluster interface {
		BindPacketFunc(packetFunc network.PacketFunc)
	}
)

func (op *Op) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func (op *Op) IsActorType(actorType ACTOR_TYPE) bool {
	return op.m_type == actorType
}

func WithName(name string) OpOption {
	return func(op *Op) {
		op.m_name = name
	}
}

func WithRpcMethodMap(rpcMethodMap map[string] string) OpOption {
	return func(op *Op) {
		op.m_RpcMethodMap = rpcMethodMap
	}
}

func WithType(actor_type ACTOR_TYPE) OpOption {
	return func(op *Op) {
		op.m_type = actor_type
	}
}

func (this *ActorMgr) Init() {
	this.m_ActorMap = make(map[reflect.Type]IActor)
	this.m_ActorNameMap = make(map[string]IActor)
	this.m_MsgMap = make(map[string]IActor)
	this.m_RpcMethodMap = map[reflect.Type] map[string] string{}
	this.m_PlayerMap = make(map[int64]IActor)
	this.m_PlayerLock = &sync.RWMutex{}
}

func (this *ActorMgr) Start() {
	this.m_bStart = true
}

func (this *ActorMgr) RegisterActor(pActor IActor, params ...OpOption) {
	op := Op{}
	op.applyOpts(params)
	if len(op.m_name) == 0 {
		op.m_name = base.GetClassName(pActor)
	}
	rType := reflect.TypeOf(pActor)
	_, bEx := this.m_ActorMap[rType]
	if bEx {
		log.Panicf("InitActor actor[%s] must  global variable", op.m_name)
		return
	}

	//rpc
	shareRpcMethodMap := GetRpcMethodMap(rType, "share_rpc")
	methodNum := rType.NumMethod()
	this.m_RpcMethodMap[rType] = map[string]string{}
	for i := 0; i < methodNum; i++{
		m := rType.Method(i)
		if m.Type.NumIn() >= 2{
			if m.Type.In(1).String() == "context.Context" {
				funcName := strings.ToLower(m.Name)
				methodName := m.Name
				_, bInShare := shareRpcMethodMap[funcName]
				if !bInShare{
					pMsgHandle, bEx := this.m_MsgMap[funcName]
					if bEx && pMsgHandle != nil{
						log.Panicf("RegisterFuncName [%s] exist_actor [%s] actor [%s]", methodName, pMsgHandle.GetName(), op.m_name)
						return
					}
					this.m_MsgMap[funcName] = pActor
				}
				this.m_RpcMethodMap[rType][funcName] = methodName
			}
		}
	}

	op.m_RpcMethodMap = this.m_RpcMethodMap[rType]
	pActor.Register(pActor, op)
	this.m_ActorMap[rType] = pActor
	this.m_ActorNameMap[op.m_name] = pActor
}

func (this *ActorMgr) AddPlayer(pActor IActor) {
	rType := reflect.TypeOf(pActor)
	op := Op{m_type:ACTOR_TYPE_PLAYER, m_name: this.m_ActorMap[rType].GetName(), m_RpcMethodMap: this.m_RpcMethodMap[rType]}
	pActor.Register(pActor, op)
    this.m_PlayerLock.Lock()
    this.m_PlayerMap[pActor.GetId()] = pActor
    this.m_PlayerLock.Unlock()
}

func (this *ActorMgr) DelPlayer(Id int64) {
	this.m_PlayerLock.Lock()
	delete(this.m_PlayerMap, Id)
	this.m_PlayerLock.Unlock()
}

func (this *ActorMgr) GetPlayer(Id int64) IActor{
	this.m_PlayerLock.RLock()
	pActor, bEx := this.m_PlayerMap[Id]
	this.m_PlayerLock.RUnlock()
	if bEx{
		return pActor
	}
	return nil
}

func (this *ActorMgr) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	this.SendActor(funcName, head, rpc.Marshal(head, funcName, params...))
}

func (this *ActorMgr) SendActor(funcName string, head rpc.RpcHead, packet rpc.Packet) bool{
    var pActor IActor
	funcName = strings.ToLower(funcName)
    bEx := false
    if head.ActorName != ""{
		pActor, bEx = this.m_ActorNameMap[head.ActorName]
    }else{
		pActor, bEx = this.m_MsgMap[funcName]
    }

    if bEx && pActor != nil{
        if pActor.HasRpc(funcName){
			switch pActor.GetActorType(){
			case ACTOR_TYPE_SINGLETON:
				pActor.GetAcotr().Send(head, packet)
				return true
			case ACTOR_TYPE_PLAYER:
				if head.Id != 0{
					pActor := this.GetPlayer(head.Id)
					if pActor != nil{
						pActor.GetAcotr().Send(head, packet)
						return true
					}
				}
			}
		}
    }
	return false
}

func (this *ActorMgr) PacketFunc(packet rpc.Packet) bool {
	rpcPacket, head := rpc.Unmarshal(packet.Buff)
	packet.RpcPacket = rpcPacket
	head.SocketId = packet.Id
	head.Reply = packet.Reply
	return this.SendActor(rpcPacket.FuncName, head, packet)
}

var (
	MGR *ActorMgr
)

func init() {
	MGR = &ActorMgr{}
	MGR.Init()
}
