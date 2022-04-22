package actor

import (
	"gonet/base"
	"gonet/network"
	"gonet/rpc"
	"log"
	"reflect"
)

type ACTOR_TYPE uint32

const (
	ACTOR_TYPE_SINGLETON ACTOR_TYPE = iota //单列
	ACTOR_TYPE_VIRTUAL   ACTOR_TYPE = iota //玩家 必须初始一个全局的actor 作为类型判断
	ACTOR_TYPE_POOL      ACTOR_TYPE = iota //固定数量actor池
) //ACTOR_TYPE

//一些全局的actor,不可删除的,不用锁考虑性能
//不是全局的actor,请使用actor pool
type (
	Op struct {
		m_name string //name
		m_type ACTOR_TYPE
		m_Pool IActorPool //ACTOR_TYPE_VIRTUAL ACTOR_TYPE_POOL
	}

	OpOption func(*Op)

	ActorMgr struct {
		m_ActorTypeMap map[reflect.Type]IActor
		m_ActorMap     map[string]IActor
		m_bStart       bool
	}

	IActorMgr interface {
		Init()
		RegisterActor(pActor IActor, params ...OpOption) //注册回调
		PacketFunc(rpc.Packet) bool                      //回调函数
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

func WithType(actor_type ACTOR_TYPE) OpOption {
	return func(op *Op) {
		op.m_type = actor_type
	}
}

func withPool(pPool IActorPool) OpOption { //ACTOR_TYPE_VIRTUAL ACTOR_TYPE_POOL
	return func(op *Op) {
		op.m_Pool = pPool
	}
}

func (this *ActorMgr) Init() {
	this.m_ActorTypeMap = make(map[reflect.Type]IActor)
	this.m_ActorMap = make(map[string]IActor)
}

func (this *ActorMgr) Start() {
	this.m_bStart = true
}

func (this *ActorMgr) RegisterActor(pActor IActor, params ...OpOption) {
	op := Op{}
	op.applyOpts(params)
	rType := reflect.TypeOf(pActor)
	name := base.GetClassName(rType)
	_, bEx := this.m_ActorTypeMap[rType]
	if bEx {
		log.Panicf("InitActor actor[%s] must  global variable", name)
		return
	}

	op.m_name = name
	pActor.register(pActor, op)
	this.m_ActorTypeMap[rType] = pActor
	this.m_ActorMap[name] = pActor
	if op.m_Pool != nil {
		pActor.bindPool(op.m_Pool)
	}
}

func (this *ActorMgr) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	this.SendActor(funcName, head, rpc.Marshal(&head, &funcName, params...))
}

func (this *ActorMgr) SendActor(funcName string, head rpc.RpcHead, packet rpc.Packet) bool {
	var pActor IActor
	bEx := false
	pActor, bEx = this.m_ActorMap[head.ActorName]
	if bEx && pActor != nil {
		if pActor.HasRpc(funcName) {
			switch pActor.GetActorType() {
			case ACTOR_TYPE_SINGLETON:
				pActor.Acotr().Send(head, packet)
				return true
			case ACTOR_TYPE_VIRTUAL:
				return pActor.getPool().SendAcotr(head, packet)
			case ACTOR_TYPE_POOL:
				return pActor.getPool().SendAcotr(head, packet)
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
