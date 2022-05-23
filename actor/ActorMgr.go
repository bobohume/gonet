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
	ACTOR_TYPE_STUB      ACTOR_TYPE = iota //stub
) //ACTOR_TYPE

//一些全局的actor,不可删除的,不用锁考虑性能
//不是全局的actor,请使用actor pool
type (
	Op struct {
		name      string //name
		actorType ACTOR_TYPE
		pool      IActorPool //ACTOR_TYPE_VIRTUAL ACTOR_TYPE_POOL
	}

	OpOption func(*Op)

	ActorMgr struct {
		actorTypeMap map[reflect.Type]IActor
		actorMap     map[string]IActor
		isStart      bool
	}

	IActorMgr interface {
		Init()
		RegisterActor(ac IActor, params ...OpOption) //注册回调
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
	return op.actorType == actorType
}

func WithType(actor_type ACTOR_TYPE) OpOption {
	return func(op *Op) {
		op.actorType = actor_type
	}
}

func withPool(pPool IActorPool) OpOption { //ACTOR_TYPE_VIRTUAL ACTOR_TYPE_POOL
	return func(op *Op) {
		op.pool = pPool
	}
}

func (a *ActorMgr) Init() {
	a.actorTypeMap = make(map[reflect.Type]IActor)
	a.actorMap = make(map[string]IActor)
}

func (a *ActorMgr) Start() {
	a.isStart = true
}

func (a *ActorMgr) RegisterActor(ac IActor, params ...OpOption) {
	op := Op{}
	op.applyOpts(params)
	rType := reflect.TypeOf(ac)
	name := base.GetClassName(rType)
	_, bEx := a.actorTypeMap[rType]
	if bEx {
		log.Panicf("InitActor actor[%s] must  global variable", name)
		return
	}

	op.name = name
	ac.register(ac, op)
	a.actorTypeMap[rType] = ac
	a.actorMap[name] = ac
	if op.pool != nil {
		ac.bindPool(op.pool)
	}
}

func (a *ActorMgr) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SocketId = 0
	a.SendActor(funcName, head, rpc.Marshal(&head, &funcName, params...))
}

func (a *ActorMgr) SendActor(funcName string, head rpc.RpcHead, packet rpc.Packet) bool {
	var ac IActor
	bEx := false
	ac, bEx = a.actorMap[head.ActorName]
	if bEx && ac != nil {
		if ac.HasRpc(funcName) {
			switch ac.GetActorType() {
			case ACTOR_TYPE_SINGLETON:
				ac.Acotr().Send(head, packet)
				return true
			case ACTOR_TYPE_VIRTUAL:
				return ac.getPool().SendAcotr(head, packet)
			case ACTOR_TYPE_POOL:
				return ac.getPool().SendAcotr(head, packet)
			}
		}
	}
	return false
}

func (a *ActorMgr) PacketFunc(packet rpc.Packet) bool {
	rpcPacket, head := rpc.Unmarshal(packet.Buff)
	packet.RpcPacket = rpcPacket
	head.SocketId = packet.Id
	head.Reply = packet.Reply
	return a.SendActor(rpcPacket.FuncName, head, packet)
}

var (
	MGR *ActorMgr
)

func init() {
	MGR = &ActorMgr{}
	MGR.Init()
}
