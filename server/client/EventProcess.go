package main

import (
	"context"
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/message"
	"gonet/server/game/lmath"
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
		m_Dh        base.Dh
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

func (this *EventProcess) SendPacket(packet proto.Message) {
	this.Client.Send(rpc.RpcHead{}, rpc.Packet{Buff: message.Encode(packet)})
}

func (this *EventProcess) PacketFunc(packet1 rpc.Packet) bool {
	packetId, data := message.Decode(packet1.Buff)
	packet := message.GetPakcet(packetId)
	if packet == nil {
		return true
	}
	err := message.UnmarshalText(packet, data)
	if err == nil {
		this.Send(rpc.RpcHead{}, rpc.Marshal(rpc.RpcHead{}, message.GetMessageName(packet), packet))
		return true
	}

	return true
}

func (this *EventProcess) Init() {
	this.Actor.Init()
	this.Pos = lmath.Point3F{1, 1, 1}
	this.m_Dh.Init()
	this.RegisterTimer((network.HEART_TIME_OUT/3)*1000*1000*1000, this.Update) //定时器
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *EventProcess) LoginGame() {
	packet1 := &message.C_W_Game_LoginRequset{PacketHead: message.BuildPacketHead(this.AccountId, rpc.SERVICE_GATESERVER),
		PlayerId: this.PlayerId}
	this.SendPacket(packet1)
}

var (
	id int32
)

func (this *EventProcess) LoginAccount() {
	id := atomic.AddInt32(&id, 1)
	this.AccountName = fmt.Sprintf("test321%d", id)
	this.PassWd = base.MD5(ToSlat(this.AccountName, "123456"))
	//this.AccountName = fmt.Sprintf("test%d", base.RAND.RandI(0, 7000))
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, rpc.SERVICE_GATESERVER),
		AccountName: this.AccountName, Password: this.PassWd, BuildNo: base.BUILD_NO, Key: this.m_Dh.ShareKey()}
	this.SendPacket(packet1)
}

func (this *EventProcess) LoginGate() {
	packet1 := &message.C_G_LoginResquest{PacketHead: message.BuildPacketHead(0, rpc.SERVICE_GATESERVER),
		Key: this.m_Dh.PubKey()}
	this.SendPacket(packet1)
}

var (
	PACKET *EventProcess
)

func (this *EventProcess) Move(yaw float32, time float32) {
	packet1 := &message.C_Z_Move{PacketHead: message.BuildPacketHead(this.AccountId, rpc.SERVICE_GATESERVER),
		Move: &message.C_Z_Move_Move{Mode: 0, Normal: &message.C_Z_Move_Move_Normal{Pos: &message.Point3F{X: this.Pos.X, Y: this.Pos.Y, Z: this.Pos.Z}, Yaw: yaw, Duration: time}}}
	this.SendPacket(packet1)
}

func (this *EventProcess) Update() {
	packet1 := &message.HeardPacket{}
	this.SendPacket(packet1)
}

func (this *EventProcess) W_C_SelectPlayerResponse(ctx context.Context, packet *message.W_C_SelectPlayerResponse) {
	this.AccountId = packet.GetAccountId()
	nLen := len(packet.GetPlayerData())
	//fmt.Println(len(packet.PlayerData), this.AccountId, packet.PlayerData)
	if nLen == 0 {
		packet1 := &message.C_W_CreatePlayerRequest{PacketHead: message.BuildPacketHead(this.AccountId, rpc.SERVICE_GATESERVER),
			PlayerName: "我是大坏蛋",
			Sex:        int32(0)}
		this.SendPacket(packet1)
	} else {
		this.PlayerId = packet.GetPlayerData()[0].GetPlayerID()
		this.LoginGame()
	}
}

func (this *EventProcess) W_C_CreatePlayerResponse(ctx context.Context, packet *message.W_C_CreatePlayerResponse) {
	if packet.GetError() == 0 {
		this.PlayerId = packet.GetPlayerId()
		this.LoginGame()
	} else { //创建失败

	}
}

func (this *EventProcess) G_C_LoginResponse(ctx context.Context, packet *message.G_C_LoginResponse) {
	this.m_Dh.ExchangePubk(packet.GetKey())
	this.LoginAccount()
}

func (this *EventProcess) A_C_LoginResponse(ctx context.Context, packet *message.A_C_LoginResponse) {
	if packet.GetError() == base.ACCOUNT_NOEXIST {
		packet1 := &message.C_A_RegisterRequest{PacketHead: message.BuildPacketHead(0, rpc.SERVICE_GATESERVER),
			AccountName: packet.AccountName, Password: this.PassWd}
		this.SendPacket(packet1)
	} else if packet.GetError() == base.PASSWORD_ERROR {
		fmt.Println("账号【", packet.GetAccountName(), "】密码错误")
	}
}

func (this *EventProcess) A_C_RegisterResponse(ctx context.Context, packet *message.A_C_RegisterResponse) {
	//注册失败
	if packet.GetError() != 0 {
	}
}

func (this *EventProcess) W_C_ChatMessage(ctx context.Context, packet *message.W_C_ChatMessage) {
	fmt.Println("收到【", packet.GetSenderName(), "】发送的消息[", packet.GetMessage()+"]")
}

//map
func (this *EventProcess) Z_C_LoginMap(ctx context.Context, packet *message.Z_C_LoginMap) {
	this.SimId = packet.GetId()
	this.Pos = lmath.Point3F{packet.GetPos().GetX(), packet.GetPos().GetY(), packet.GetPos().GetZ()}
	this.Rot = lmath.Point3F{0, 0, packet.GetRotation()}
	//fmt.Println("login map")
}

func (this *EventProcess) Z_C_ENTITY(ctx context.Context, packet *message.Z_C_ENTITY) {
	for _, v := range packet.EntityInfo {
		if v.Data != nil {
			if v.Data.RemoveFlag {
				fmt.Printf("Z_C_ENTITY_DATA  destory:[%d], [%d], [%t]\n", v.GetId(), v.Data.Type, v.Data.RemoveFlag)
				continue
			}
			fmt.Printf("Z_C_ENTITY_DATA :[%d], [%d], [%t]\n", v.GetId(), v.Data.Type, v.Data.RemoveFlag)
		}
		if v.Move != nil {
			if v.Id == this.SimId {
				this.Pos = lmath.Point3F{v.Move.GetPos().GetX(), v.Move.GetPos().GetY(), v.Move.GetPos().GetZ()}
				this.Rot = lmath.Point3F{0, 0, v.Move.GetRotation()}
			}
			fmt.Printf("Z_C_ENTITY_MOVE :[%d], Pos:[x:%f, y:%f, z:%f], Rot[%f]\n", v.GetId(), v.Move.GetPos().GetX(), v.Move.GetPos().GetY(), v.Move.GetPos().GetZ(), v.Move.GetRotation())
		}
	}
}

//链接断开
func (this *EventProcess) DISCONNECT(ctx context.Context, socketId uint32) {
	this.Stop()
}
