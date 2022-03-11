package gate

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
	LoginAccountRequest    = strings.ToLower("LoginAccountRequest")
)

type (
	UserPrcoess struct {
		actor.Actor
		m_KeyMap map[uint32]*base.Dh
	}

	IUserPrcoess interface {
		actor.IActor

		CheckClientEx(uint32, string, rpc.RpcHead) bool
		CheckClient(uint32, string, rpc.RpcHead) *Player
		SwtichSendToGame(uint32, string, rpc.RpcHead, []byte)
		SwtichSendToGM(uint32, string, rpc.RpcHead, []byte)
		SwtichSendToZone(uint32, string, rpc.RpcHead, []byte)

		addKey(uint32, *base.Dh)
		delKey(uint32)
	}
)

func (this *UserPrcoess) CheckClientEx(sockId uint32, packetName string, head rpc.RpcHead) bool {
	if IsCheckClient(packetName) {
		return true
	}

	playerId := SERVER.GetPlayerMgr().GetPlayerId(sockId)
	if playerId <= 0 || playerId != head.Id {
		SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
		return false
	}
	return true
}

func (this *UserPrcoess) CheckClient(sockId uint32, packetName string, head rpc.RpcHead) *Player {
	pPlayer := SERVER.GetPlayerMgr().GetPlayer(sockId)
	if pPlayer != nil && (pPlayer.PlayerID <= 0 || pPlayer.PlayerID != head.Id) {
		SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
		return nil
	}
	return pPlayer
}

func (this *UserPrcoess) SwtichSendToGame(socketId uint32, packetName string, head rpc.RpcHead, packet rpc.Packet) {
	pPlayer := this.CheckClient(socketId, packetName, head)
	if pPlayer != nil {
		head.ClusterId = pPlayer.GClusterId
		head.DestServerType = rpc.SERVICE_GAME
		SERVER.GetCluster().Send(head, packet)
	}
}

func (this *UserPrcoess) SwtichSendToGM(socketId uint32, packetName string, head rpc.RpcHead, packet rpc.Packet) {
	if this.CheckClientEx(socketId, packetName, head) == true {
		head.SendType = rpc.SEND_BALANCE
		head.DestServerType = rpc.SERVICE_GM
		SERVER.GetCluster().Send(head, packet)
	}
}

func (this *UserPrcoess) SwtichSendToZone(socketId uint32, packetName string, head rpc.RpcHead, packet rpc.Packet) {
	pPlayer := this.CheckClient(socketId, packetName, head)
	if pPlayer != nil {
		head.ClusterId = pPlayer.ZClusterId
		head.DestServerType = rpc.SERVICE_ZONE
		SERVER.GetCluster().Send(head, packet)
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
	if packetName == LoginAccountRequest {
		head.ClusterId = socketid
	}

	//解析整个包
	if packetHead.DestServerType == message.SERVICE_GAME {
		this.SwtichSendToGame(socketid, packetName, head, rpc.Marshal(head, packetName, packet))
	} else if packetHead.DestServerType == message.SERVICE_GM{
		this.SwtichSendToGM(socketid, packetName, head, rpc.Marshal(head, packetName, packet))
	} else if packetHead.DestServerType == message.SERVICE_ZONE{
		this.SwtichSendToZone(socketid, packetName, head, rpc.Marshal(head, packetName, packet))
	} else {
		actor.MGR.PacketFunc(rpc.Packet{Id: socketid, Buff: rpc.Marshal(head, packetName, packet).Buff})
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
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

func (this *UserPrcoess) C_G_LogoutRequest(ctx context.Context, playerid int, UID int) {
	SERVER.GetLog().Printf("logout Socket:%d Account:%d UID:%d ", this.GetRpcHead(ctx).SocketId, playerid, UID)
	SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{}, "DEL_ACCOUNT", this.GetRpcHead(ctx).SocketId)
}

func (this *UserPrcoess) LoginAccountRequest(ctx context.Context, packet *message.LoginAccountRequest) {
	head := this.GetRpcHead(ctx)
	dh := base.Dh{}
	dh.Init()
	dh.ExchangePubk(packet.GetKey())
	this.addKey(head.SocketId, &dh)
	head.Id = int64(base.GetMessageCode1(packet.AccountName))
	packet.Key = dh.PubKey()
	this.SwtichSendToGM(head.SocketId, base.ToLower("LoginAccountRequest"), head, rpc.Marshal(head, base.ToLower("LoginAccountRequest"), packet))
}

func (this *UserPrcoess) LoginPlayerRequset(ctx context.Context, packet *message.LoginPlayerRequset) {
	head := this.GetRpcHead(ctx)
	dh, bEx := this.m_KeyMap[head.SocketId]
	if bEx {
		if dh.ShareKey() == packet.GetKey() {
			this.delKey(head.SocketId)
			this.SwtichSendToGM(head.SocketId, base.ToLower("LoginPlayerRequset"), head, rpc.Marshal(head, base.ToLower("LoginPlayerRequset"), packet))
		} else {
			SERVER.GetLog().Println("client key cheat", dh.ShareKey(), packet.GetKey())
		}
	}
}

func (this *UserPrcoess) DISCONNECT(ctx context.Context, socketid uint32) {
	this.delKey(socketid)
}
