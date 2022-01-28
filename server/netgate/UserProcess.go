package netgate

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/message"
	"strings"
)

var (
	C_A_LoginRequest    = strings.ToLower("C_A_LoginRequest")
	C_A_RegisterRequest = strings.ToLower("C_A_RegisterRequest")
)

type (
	UserPrcoess struct {
		actor.Actor
		m_KeyMap map[uint32]*base.Dh
	}

	IUserPrcoess interface {
		actor.IActor

		CheckClientEx(uint32, string, rpc.RpcHead) bool
		CheckClient(uint32, string, rpc.RpcHead) *AccountInfo
		SwtichSendToWorld(uint32, string, rpc.RpcHead, []byte)
		SwtichSendToAccount(uint32, string, rpc.RpcHead, []byte)
		SwtichSendToZone(uint32, string, rpc.RpcHead, []byte)

		addKey(uint32, *base.Dh)
		delKey(uint32)
	}
)

func (this *UserPrcoess) CheckClientEx(sockId uint32, packetName string, head rpc.RpcHead) bool {
	if IsCheckClient(packetName) {
		return true
	}

	accountId := SERVER.GetPlayerMgr().GetAccount(sockId)
	if accountId <= 0 || accountId != head.Id {
		SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
		return false
	}
	return true
}

func (this *UserPrcoess) CheckClient(sockId uint32, packetName string, head rpc.RpcHead) *AccountInfo {
	pAccountInfo := SERVER.GetPlayerMgr().GetAccountInfo(sockId)
	if pAccountInfo != nil && (pAccountInfo.AccountId <= 0 || pAccountInfo.AccountId != head.Id) {
		SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
		return nil
	}
	return pAccountInfo
}

func (this *UserPrcoess) SwtichSendToWorld(socketId uint32, packetName string, head rpc.RpcHead, buff []byte) {
	pAccountInfo := this.CheckClient(socketId, packetName, head)
	if pAccountInfo != nil {
		head.ClusterId = pAccountInfo.WClusterId
		head.DestServerType = rpc.SERVICE_WORLDSERVER
		SERVER.GetCluster().Send(head, buff)
	}
}

func (this *UserPrcoess) SwtichSendToAccount(socketId uint32, packetName string, head rpc.RpcHead, buff []byte) {
	if this.CheckClientEx(socketId, packetName, head) == true {
		head.SendType = rpc.SEND_BALANCE
		head.DestServerType = rpc.SERVICE_ACCOUNTSERVER
		SERVER.GetCluster().Send(head, buff)
	}
}

func (this *UserPrcoess) SwtichSendToZone(socketId uint32, packetName string, head rpc.RpcHead, buff []byte) {
	pAccountInfo := this.CheckClient(socketId, packetName, head)
	if pAccountInfo != nil {
		head.ClusterId = pAccountInfo.ZClusterId
		head.DestServerType = rpc.SERVICE_ZONESERVER
		SERVER.GetCluster().Send(head, buff)
	}
}

func (this *UserPrcoess) PacketFunc(packet1 rpc.Packet) bool {
	buff := packet1.Buff
	socketid := packet1.Id
	packetId, data := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet == nil {
		//客户端主动断开
		if packetId == network.DISCONNECTINT {
			stream := base.NewBitStream(buff, len(buff))
			stream.ReadInt(32)
			SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{}, "DEL_ACCOUNT", uint32(stream.ReadInt(32)))
			this.SendMsg(rpc.RpcHead{}, "DISCONNECT", socketid)
		} else if packetId == network.HEART_PACKET { //心跳netsocket做处理，这里不处理
		} else {
			SERVER.GetLog().Printf("包解析错误1  socket=%d", socketid)
		}
		return true
	}

	//获取配置的路由地址
	destServerType := packet.(message.Packet).GetPacketHead().DestServerType
	err := message.UnmarshalText(packet, data)
	if err != nil {
		SERVER.GetLog().Printf("包解析错误2  socket=%d", socketid)
		return true
	}

	packetHead := packet.(message.Packet).GetPacketHead()
	packetHead.DestServerType = destServerType
	if packetHead == nil || packetHead.Ckx != message.Default_Ipacket_Ckx || packetHead.Stx != message.Default_Ipacket_Stx {
		SERVER.GetLog().Printf("(A)致命的越界包,已经被忽略 socket=%d", socketid)
		return true
	}

	packetName := message.GetMessageName(packet)
	head := rpc.RpcHead{Id: packetHead.Id, SrcClusterId: SERVER.GetCluster().Id()}
	if packetName == C_A_LoginRequest {
		head.ClusterId = socketid
	} else if packetName == C_A_RegisterRequest {
		head.ClusterId = socketid
	}

	//解析整个包
	if packetHead.DestServerType == message.SERVICE_WORLDSERVER {
		this.SwtichSendToWorld(socketid, packetName, head, rpc.Marshal(head, packetName, packet))
	} else if packetHead.DestServerType == message.SERVICE_ACCOUNTSERVER {
		this.SwtichSendToAccount(socketid, packetName, head, rpc.Marshal(head, packetName, packet))
	} else if packetHead.DestServerType == message.SERVICE_ZONESERVER {
		this.SwtichSendToZone(socketid, packetName, head, rpc.Marshal(head, packetName, packet))
	} else {
		this.Actor.PacketFunc(rpc.Packet{Id: socketid, Buff: rpc.Marshal(head, packetName, packet)})
	}

	return true
}

func (this *UserPrcoess) addKey(SocketId uint32, pDh *base.Dh) {
	this.m_KeyMap[SocketId] = pDh
}

func (this *UserPrcoess) delKey(SocketId uint32) {
	delete(this.m_KeyMap, SocketId)
}

func (this *UserPrcoess) Init() {
	this.Actor.Init()
	this.m_KeyMap = map[uint32]*base.Dh{}
	this.RegisterCall("C_G_LogoutRequest", func(ctx context.Context, accountId int, UID int) {
		SERVER.GetLog().Printf("logout Socket:%d Account:%d UID:%d ", this.GetRpcHead(ctx).SocketId, accountId, UID)
		SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{}, "DEL_ACCOUNT", this.GetRpcHead(ctx).SocketId)
		SendToClient(this.GetRpcHead(ctx).SocketId, &message.C_G_LogoutResponse{PacketHead: message.BuildPacketHead(0, 0)})
	})

	this.RegisterCall("C_G_LoginResquest", func(ctx context.Context, packet *message.C_G_LoginResquest) {
		head := this.GetRpcHead(ctx)
		dh := base.Dh{}
		dh.Init()
		dh.ExchangePubk(packet.GetKey())
		this.addKey(head.SocketId, &dh)
		SendToClient(head.SocketId, &message.G_C_LoginResponse{PacketHead: message.BuildPacketHead(0, 0), Key: dh.PubKey()})
	})

	this.RegisterCall("C_A_LoginRequest", func(ctx context.Context, packet *message.C_A_LoginRequest) {
		head := this.GetRpcHead(ctx)
		dh, bEx := this.m_KeyMap[head.SocketId]
		if bEx {
			if dh.ShareKey() == packet.GetKey() {
				this.delKey(head.SocketId)
				head.Id = int64(base.ToHash(packet.AccountName))
				this.SwtichSendToAccount(head.SocketId, base.ToLower("C_A_LoginRequest"), head, rpc.Marshal(head, base.ToLower("C_A_LoginRequest"), packet))
			} else {
				SERVER.GetLog().Println("client key cheat", dh.ShareKey(), packet.GetKey())
			}
		}
	})

	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketid uint32) {
		this.delKey(socketid)
	})

	this.Actor.Start()
}
