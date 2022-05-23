package gate

import (
	"context"
	"gonet/actor"
	"gonet/base"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/message"
)

type (
	UserPrcoess struct {
		actor.Actor
		keyMap map[uint32]*base.Dh
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

func (u *UserPrcoess) CheckClientEx(sockId uint32, packetName string, head rpc.RpcHead) bool {
	if IsCheckClient(packetName) {
		return true
	}

	playerId := SERVER.GetPlayerMgr().GetPlayerId(sockId)
	if playerId <= 0 || playerId != head.Id {
		base.LOG.Fatalf("Old socket communication or viciousness[%d].", sockId)
		return false
	}
	return true
}

func (u *UserPrcoess) CheckClient(sockId uint32, packetName string, head rpc.RpcHead) *Player {
	pPlayer := SERVER.GetPlayerMgr().GetPlayer(sockId)
	if pPlayer != nil && (pPlayer.PlayerID <= 0 || pPlayer.PlayerID != head.Id) {
		base.LOG.Fatalf("Old socket communication or viciousness[%d].", sockId)
		return nil
	}
	return pPlayer
}

func (u *UserPrcoess) SwtichSendToGame(socketId uint32, packetName string, head rpc.RpcHead, packet rpc.Packet) {
	pPlayer := u.CheckClient(socketId, packetName, head)
	if pPlayer != nil {
		head.ClusterId = pPlayer.GClusterId
		head.DestServerType = rpc.SERVICE_GAME
		SERVER.GetCluster().Send(head, packet)
	}
}

func (u *UserPrcoess) SwtichSendToGM(socketId uint32, packetName string, head rpc.RpcHead, packet rpc.Packet) {
	if u.CheckClientEx(socketId, packetName, head) == true {
		head.DestServerType = rpc.SERVICE_GM
		SERVER.GetCluster().Send(head, packet)
	}
}

func (u *UserPrcoess) SwtichSendToZone(socketId uint32, packetName string, head rpc.RpcHead, packet rpc.Packet) {
	pPlayer := u.CheckClient(socketId, packetName, head)
	if pPlayer != nil {
		head.ClusterId = pPlayer.ZClusterId
		head.DestServerType = rpc.SERVICE_ZONE
		SERVER.GetCluster().Send(head, packet)
	}
}

func (u *UserPrcoess) PacketFunc(packet1 rpc.Packet) bool {
	buff := packet1.Buff
	socketid := packet1.Id
	packetId, data := message.Decode(buff)
	packetRoute := message.GetPakcetRoute(packetId)
	if packetRoute == nil {
		//客户端主动断开
		if packetId == network.DISCONNECTINT {
			stream := base.NewBitStream(buff, len(buff))
			stream.ReadInt(32)
			SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{}, "DEL_ACCOUNT", uint32(stream.ReadInt(32)))
			u.SendMsg(rpc.RpcHead{}, "DISCONNECT", socketid)
		} else if packetId == network.HEART_PACKET { //心跳netsocket做处理，这里不处理
		} else {
			base.LOG.Printf("包解析错误1  socket=%d", socketid)
		}
		return true
	}

	//获取配置的路由地址
	packet := packetRoute.Func()
	err := message.UnmarshalText(packet, data)
	if err != nil {
		base.LOG.Printf("包解析错误2  socket=%d", socketid)
		return true
	}

	packetHead := packet.(message.Packet).GetPacketHead()
	if packetHead == nil || packetHead.Ckx != message.Default_Ipacket_Ckx || packetHead.Stx != message.Default_Ipacket_Stx {
		base.LOG.Printf("(A)致命的越界包,已经被忽略 socket=%d", socketid)
		return true
	}

	packetName := packetRoute.FuncName
	head := rpc.RpcHead{Id: packetHead.Id, SrcClusterId: SERVER.GetCluster().Id()}
	rpcPacket := rpc.Marshal(&head, &packetName, packet)
	//解析整个包
	if head.DestServerType == rpc.SERVICE_GAME {
		u.SwtichSendToGame(socketid, packetName, head, rpcPacket)
	} else if head.DestServerType == rpc.SERVICE_GM {
		u.SwtichSendToGM(socketid, packetName, head, rpcPacket)
	} else if head.DestServerType == rpc.SERVICE_ZONE {
		u.SwtichSendToZone(socketid, packetName, head, rpcPacket)
	} else {
		actor.MGR.PacketFunc(rpc.Packet{Id: socketid, Buff: rpcPacket.Buff})
	}

	return true
}

func (u *UserPrcoess) addKey(SocketId uint32, dh *base.Dh) {
	u.keyMap[SocketId] = dh
}

func (u *UserPrcoess) delKey(SocketId uint32) {
	delete(u.keyMap, SocketId)
}

func (u *UserPrcoess) Init() {
	u.Actor.Init()
	u.keyMap = map[uint32]*base.Dh{}
	actor.MGR.RegisterActor(u)
	u.Actor.Start()
}

func (u *UserPrcoess) C_G_LogoutRequest(ctx context.Context, playerid int, UID int) {
	base.LOG.Printf("logout Socket:%d Account:%d UID:%d ", u.GetRpcHead(ctx).SocketId, playerid, UID)
	SERVER.GetPlayerMgr().SendMsg(rpc.RpcHead{}, "DEL_ACCOUNT", u.GetRpcHead(ctx).SocketId)
}

func (u *UserPrcoess) LoginAccountRequest(ctx context.Context, packet *message.LoginAccountRequest) {
	head := u.GetRpcHead(ctx)
	dh := base.Dh{}
	dh.Init()
	dh.ExchangePubk(packet.GetKey())
	u.addKey(head.SocketId, &dh)
	head.Id = int64(base.GetMessageCode1(packet.AccountName))
	packet.Key = dh.PubKey()
	funcName := "AccountMgr.LoginAccountRequest"
	u.SwtichSendToGM(head.SocketId, funcName, head, rpc.Marshal(&head, &funcName, packet, head.SocketId))
}

func (u *UserPrcoess) LoginPlayerRequset(ctx context.Context, packet *message.LoginPlayerRequset) {
	head := u.GetRpcHead(ctx)
	dh, bEx := u.keyMap[head.SocketId]
	if bEx {
		if dh.ShareKey() == packet.GetKey() {
			u.delKey(head.SocketId)
			funcName := "AccountMgr.LoginPlayerRequset"
			u.SwtichSendToGM(head.SocketId, funcName, head, rpc.Marshal(&head, &funcName, packet))
		} else {
			base.LOG.Println("client key cheat", dh.ShareKey(), packet.GetKey())
		}
	}
}

func (u *UserPrcoess) DISCONNECT(ctx context.Context, socketid uint32) {
	u.delKey(socketid)
}
