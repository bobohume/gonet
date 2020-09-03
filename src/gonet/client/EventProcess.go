package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
	"gonet/rpc"
	"sync/atomic"
)

type (
	EventProcess struct {
		actor.Actor

		Client *network.ClientSocket
		AccountId int64
		PlayerId int64
		AccountName string
		SimId int64
	}

	IEventProcess interface {
		actor.IActor
		LoginGame()
		LoginAccount()
		SendPacket(proto.Message)
	}
)

func SendPacket(packet proto.Message){
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	CLIENT.Send(rpc.RpcHead{}, buff)
}

func (this *EventProcess) SendPacket(packet proto.Message){
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	this.Client.Send(rpc.RpcHead{},buff)
}

func (this *EventProcess) PacketFunc(socketid uint32, buff []byte) bool {
	packetId, data := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet == nil{
		return true
	}
	err := message.UnmarshalText(packet, data)
	if err == nil{
		this.Send(rpc.RpcHead{}, rpc.Marshal(rpc.RpcHead{}, message.GetMessageName(packet), packet))
		return true
	}

	return true
}

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("W_C_SelectPlayerResponse", func(ctx context.Context, packet *message.W_C_SelectPlayerResponse) {
		this.AccountId = packet.GetAccountId()
		nLen := len(packet.GetPlayerData())
		//fmt.Println(len(packet.PlayerData), this.AccountId, packet.PlayerData)
		if nLen == 0{
			packet1 := &message.C_W_CreatePlayerRequest{PacketHead:message.BuildPacketHead( this.AccountId, message.SERVICE_GATESERVER),
				PlayerName:"我是大坏蛋",
				Sex:int32(0),}
			this.SendPacket(packet1)
		}else{
			this.PlayerId = packet.GetPlayerData()[0].GetPlayerID()
			this.LoginGame()
		}
	})

	this.RegisterCall("W_C_CreatePlayerResponse", func(ctx context.Context, packet *message.W_C_CreatePlayerResponse) {
		if packet.GetError() == 0 {
			this.PlayerId = packet.GetPlayerId()
			this.LoginGame()
		}else{//创建失败

		}
	})

	this.RegisterCall("A_C_LoginResponse", func(ctx context.Context, packet *message.A_C_LoginResponse) {
		if packet.GetError() == base.ACCOUNT_NOEXIST {
			packet1 := &message.C_A_RegisterRequest{PacketHead:message.BuildPacketHead( 0, message.SERVICE_GATESERVER),
				AccountName: packet.AccountName}
			this.SendPacket(packet1)
		}
	})

	this.RegisterCall("A_C_RegisterResponse", func(ctx context.Context, packet *message.A_C_RegisterResponse) {
		//注册失败
		if packet.GetError() != 0 {
		}
	})

	this.RegisterCall("W_C_ChatMessage", func(ctx context.Context, packet *message.W_C_ChatMessage) {
		fmt.Println("收到【", packet.GetSenderName(), "】发送的消息[", packet.GetMessage()+"]")
	})

	//map
	this.RegisterCall("Z_C_LoginMap", func(ctx context.Context, packet *message.Z_C_LoginMap) {
		this.SimId = packet.GetId()
		//fmt.Println("login map")
	})

	this.RegisterCall("Z_C_ENTITY", func(ctx context.Context, packet *message.Z_C_ENTITY) {
		for _, v := range packet.EntityInfo{
			if v.Data != nil{
				if v.Data.RemoveFlag{
					fmt.Printf("Z_C_ENTITY_DATA  destory:[%d], [%d], [%t]\n", v.GetId(), v.Data.Type, v.Data.RemoveFlag )
					continue
				}
				fmt.Printf("Z_C_ENTITY_DATA :[%d], [%d], [%t]\n",v.GetId(), v.Data.Type, v.Data.RemoveFlag )
			}
			if v.Move != nil{
				if v.Id == this.SimId{
				}
				fmt.Printf("Z_C_ENTITY_MOVE :[%d], Pos:[x:%f, y:%f, z:%f], Rot[%f]\n", v.GetId(), v.Move.GetPos().GetX(),  v.Move.GetPos().GetY(), v.Move.GetPos().GetZ(), v.Move.GetRotation())
			}
		}
	})
	this.Actor.Start()
}

func (this *EventProcess)  LoginGame(){
	packet1 := &message.C_W_Game_LoginRequset{PacketHead:message.BuildPacketHead( this.AccountId, message.SERVICE_GATESERVER),
		PlayerId:this.PlayerId,}
	this.SendPacket(packet1)
}

var(
	id int32
)

func (this *EventProcess)  LoginAccount() {
	id := atomic.AddInt32(&id, 1)
	this.AccountName = fmt.Sprintf("test%d", id)
	//this.AccountName = fmt.Sprintf("test%d", base.RAND.RandI(0, 7000))
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, message.SERVICE_GATESERVER),
		AccountName: this.AccountName, BuildNo: base.BUILD_NO}
	this.SendPacket(packet1)
}

var(
	PACKET *EventProcess
)