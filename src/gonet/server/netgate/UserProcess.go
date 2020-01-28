package netgate

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
	"gonet/rpc"
)

type(
	UserPrcoess struct{
		actor.Actor
	}

	IUserPrcoess interface {
		actor.IActor

		CheckClient(int, string, *message.Ipacket)bool
		CheckClientEx(int, string, *message.Ipacket) *AccountInfo
		SwtichSendToWorld(int, string, *message.Ipacket, []byte)
		SwtichSendToAccount(int, string, *message.Ipacket, []byte)
	}
)

func (this *UserPrcoess) CheckClient(sockId int, packetName string, packetHead *message.Ipacket) bool{
	if packetHead != nil{
		if IsCheckClient(packetName){
			return  true
		}

		accountId := SERVER.GetPlayerMgr().GetAccount(sockId)
		if accountId <= 0 || accountId != packetHead.Id {
			SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
			return false
		}
		return  true
	}
	return  false
}

func (this *UserPrcoess) CheckClientEx(sockId int, packetName string, packetHead *message.Ipacket) *AccountInfo{
	if packetHead != nil{
		if IsCheckClient(packetName){
			return  nil
		}

		pAccountInfo := SERVER.GetPlayerMgr().GetAccountInfo(sockId)
		if pAccountInfo != nil && (pAccountInfo.AccountId <= 0 || pAccountInfo.AccountId != packetHead.Id){
			SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
			return nil
		}
		return pAccountInfo
	}
	return nil
}

func (this *UserPrcoess)SwtichSendToWorld(socketId int, packetName string, packetHead *message.Ipacket, buff []byte){
	pAccountInfo := this.CheckClientEx(socketId, packetName, packetHead)
	if pAccountInfo != nil{
		buff = base.SetTcpEnd(buff)
		SERVER.GetWorldCluster().Send(pAccountInfo.WSocketId, buff)
	}
}

func (this *UserPrcoess) SwtichSendToAccount(socketId int, packetName string, packetHead *message.Ipacket, buff []byte){
	if this.CheckClient(socketId, packetName, packetHead) == true {
		buff = base.SetTcpEnd(buff)
		SERVER.GetAccountCluster().BalanceSend(buff)
	}
}

func (this *UserPrcoess) PacketFunc(socketid int, buff []byte) bool{
	packetId, data := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet == nil{
		//客户端主动断开
		if packetId == network.DISCONNECTINT{
			stream := base.NewBitStream(buff, len(buff))
			stream.ReadInt(32)
			SERVER.GetPlayerMgr().SendMsg("DEL_ACCOUNT", stream.ReadInt(32))
		}else{
			SERVER.GetLog().Printf("包解析错误1  socket=%d", socketid)
		}
		return true
	}

	//获取配置的路由地址
	destServerType := packet.(message.Packet).GetPacketHead().DestServerType
	err := message.UnmarshalText(packet, data)
	if err != nil{
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
	if packetName  == base.ToLower("C_A_LoginRequest") {
		packet.(*message.C_A_LoginRequest).SocketId = int32(socketid)
	}else if packetName  == base.ToLower("C_A_RegisterRequest") {
		packet.(*message.C_A_RegisterRequest).SocketId = int32(socketid)
	}

	//解析整个包
	if packetHead.DestServerType == message.SERVICE_WORLDSERVER{
		this.SwtichSendToWorld(socketid, packetName, packetHead, rpc.Marshal(packetName, packet))
	}else if packetHead.DestServerType == message.SERVICE_ACCOUNTSERVER{
		this.SwtichSendToAccount(socketid, packetName, packetHead, rpc.Marshal(packetName, packet))
	}else{
		this.Actor.PacketFunc(socketid, rpc.Marshal(packetName, packet))
	}

	return true
}

func (this *UserPrcoess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("C_G_LogoutRequest", func(accountId int, UID int){
		SERVER.GetLog().Printf("logout Socket:%d Account:%d UID:%d ",this.GetSocketId(), accountId,UID )
		SERVER.GetPlayerMgr().SendMsg("DEL_ACCOUNT", this.GetSocketId())
		SendToClient(this.GetSocketId(), &message.C_G_LogoutResponse{PacketHead:message.BuildPacketHead( 0, 0)})
	})

	this.Actor.Start()
}