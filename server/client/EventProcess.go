package main

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/game/lmath"
	"gonet/server/message"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

type (
	EventProcess struct {
		actor.Actor

		Client      *network.ClientSocket
		AccountId   int64
		PlayerId    int64
		AccountName string
		PassWd      string
		SimId       int64
		Pos         lmath.Point3F
		Rot         lmath.Point3F
		dh          base.Dh
	}

	IEventProcess interface {
		actor.IActor
		LoginGame()
		LoginAccount()
		SendPacket(proto.Message)
	}
)

func ToSlat(accountName string, pwd string) string {
	return fmt.Sprintf("%s__%s", accountName, pwd)
}

func ToCrc(accountName string, pwd string, buildNo string, nKey int64) uint32 {
	return base.GetMessageCode1(fmt.Sprintf("%s_%s_%s_%d", accountName, pwd, buildNo, nKey))
}

func SendPacket(packet proto.Message) {
	CLIENT.Send(rpc.RpcHead{}, rpc.Packet{Buff: message.Encode(packet)})
}

func (e *EventProcess) SendPacket(packet proto.Message) {
	e.Client.Send(rpc.RpcHead{}, rpc.Packet{Buff: message.Encode(packet)})
}

func (e *EventProcess) PacketFunc(packet1 rpc.Packet) bool {
	packetId, data := message.Decode(packet1.Buff)
	packetRoute := message.GetPakcetRoute(packetId)
	if packetRoute == nil {
		return true
	}
	packet := packetRoute.Func()
	err := message.UnmarshalText(packet, data)
	if err == nil {
		head := rpc.RpcHead{}
		e.Send(head, rpc.Marshal(&head, &packetRoute.FuncName, packet))
		return true
	}

	return true
}

func (e *EventProcess) Init() {
	e.Actor.Init()
	e.Pos = lmath.Point3F{1, 1, 1}
	e.dh.Init()
	e.RegisterTimer((network.HEART_TIME_OUT/3)*1000*1000*1000, e.Update) //定时器
	actor.MGR.RegisterActor(e)
	e.Actor.Start()
}

func (e *EventProcess) LoginGame() {
	packet1 := &message.LoginPlayerRequset{PacketHead: message.BuildPacketHead(e.AccountId, rpc.SERVICE_GATE),
		PlayerId: e.PlayerId,
		Key:      e.dh.ShareKey(),
	}
	e.SendPacket(packet1)
}

var (
	id int32
)

func (e *EventProcess) LoginAccount() {
	id := atomic.AddInt32(&id, 1)
	e.AccountName = fmt.Sprintf("test3211%d", id)
	e.PassWd = base.MD5(ToSlat(e.AccountName, "123456"))
	//e.AccountName = fmt.Sprintf("test%d", base.RandI(0, 7000))
	packet1 := &message.LoginAccountRequest{PacketHead: message.BuildPacketHead(0, rpc.SERVICE_GATE),
		AccountName: e.AccountName, Password: e.PassWd, BuildNo: base.BUILD_NO, Key: e.dh.PubKey()}
	e.SendPacket(packet1)
}

var (
	PACKET *EventProcess
)

func (e *EventProcess) Move(yaw float32, time float32) {
	packet1 := &message.C_Z_Move{PacketHead: message.BuildPacketHead(e.PlayerId, rpc.SERVICE_GATE),
		Move: &message.C_Z_Move_Move{Mode: 0, Normal: &message.C_Z_Move_Move_Normal{Pos: &message.Point3F{X: e.Pos.X, Y: e.Pos.Y, Z: e.Pos.Z}, Yaw: yaw, Duration: time}}}
	e.SendPacket(packet1)
}

func (e *EventProcess) Update() {
	packet1 := &message.HeardPacket{}
	e.SendPacket(packet1)
}

func (e *EventProcess) LoginAccountResponse(ctx context.Context, packet *message.LoginAccountResponse) {
	if packet.GetError() == base.ACCOUNT_NOEXIST {
	} else if packet.GetError() == base.PASSWORD_ERROR {
		fmt.Println("账号【", packet.GetAccountName(), "】密码错误")
	}
}

func (e *EventProcess) SelectPlayerResponse(ctx context.Context, packet *message.SelectPlayerResponse) {
	e.AccountId = packet.GetAccountId()
	e.dh.ExchangePubk(packet.GetKey())
	nLen := len(packet.GetPlayerData())
	//fmt.Println(len(packet.PlayerData), e.AccountId, packet.PlayerData)
	if nLen == 0 {
		packet1 := &message.CreatePlayerRequest{PacketHead: message.BuildPacketHead(e.AccountId, rpc.SERVICE_GATE),
			PlayerName: "我是大坏蛋",
			Sex:        int32(0)}
		e.SendPacket(packet1)
	} else {
		e.PlayerId = packet.GetPlayerData()[0].GetPlayerID()
		e.LoginGame()
	}
}

func (e *EventProcess) ChatMessageResponse(ctx context.Context, packet *message.ChatMessageResponse) {
	fmt.Println("收到【", packet.GetSenderName(), "】发送的消息[", packet.GetMessage()+"]")
}

// map
func (e *EventProcess) Z_C_LoginMap(ctx context.Context, packet *message.Z_C_LoginMap) {
	e.SimId = packet.GetId()
	e.Pos = lmath.Point3F{packet.GetPos().GetX(), packet.GetPos().GetY(), packet.GetPos().GetZ()}
	e.Rot = lmath.Point3F{0, 0, packet.GetRotation()}
	//fmt.Println("login map")
}

func (e *EventProcess) Z_C_ENTITY(ctx context.Context, packet *message.Z_C_ENTITY) {
	for _, v := range packet.EntityInfo {
		if v.Data != nil {
			if v.Data.RemoveFlag {
				fmt.Printf("Z_C_ENTITY_DATA  destory:[%d], [%d], [%t]\n", v.GetId(), v.Data.Type, v.Data.RemoveFlag)
				continue
			}
			fmt.Printf("Z_C_ENTITY_DATA :[%d], [%d], [%t]\n", v.GetId(), v.Data.Type, v.Data.RemoveFlag)
		}
		if v.Move != nil {
			if v.Id == e.SimId {
				e.Pos = lmath.Point3F{v.Move.GetPos().GetX(), v.Move.GetPos().GetY(), v.Move.GetPos().GetZ()}
				e.Rot = lmath.Point3F{0, 0, v.Move.GetRotation()}
			}
			fmt.Printf("Z_C_ENTITY_MOVE :[%d], Pos:[x:%f, y:%f, z:%f], Rot[%f]\n", v.GetId(), v.Move.GetPos().GetX(), v.Move.GetPos().GetY(), v.Move.GetPos().GetZ(), v.Move.GetRotation())
		}
	}
}

// 链接断开
func (e *EventProcess) DISCONNECT(ctx context.Context, socketId uint32) {
	e.Stop()
}
