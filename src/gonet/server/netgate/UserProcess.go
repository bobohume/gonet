package netgate

import (
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"
)

type(
	UserPrcoess struct{
		actor.Actor
	}

	IUserPrcoess interface {
		actor.IActor

		CheckClient(int, string, interface{})bool
		CheckClientEx(int, string, interface{}) *AccountInfo
		SwtichSendToWorld(int, string, interface{}, []byte)
		SwtichSendToAccount(int, string, interface{}, []byte)
	}
)

func (this *UserPrcoess) CheckClient(sockId int, packetName string, packet interface{}) bool{
	packetHead := packet.(*message.Ipacket)
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

func (this *UserPrcoess) CheckClientEx(sockId int, packetName string, packet interface{}) *AccountInfo{
	packetHead := packet.(*message.Ipacket)
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

func (this *UserPrcoess)SwtichSendToWorld(socketId int, packetName string, packet interface{}, buff []byte){
	pAccountInfo := this.CheckClientEx(socketId, packetName, packet)
	if pAccountInfo != nil{
		buff = base.SetTcpEnd(buff)
		SERVER.GetWorldCluster().Send(pAccountInfo.WSocketId, buff)
	}
}

func (this *UserPrcoess) SwtichSendToAccount(socketId int, packetName string, packet interface{}, buff []byte){
	if this.CheckClient(socketId, packetName, packet) == true {
		buff = base.SetTcpEnd(buff)
		SERVER.GetAccountCluster().BalanceSend(buff)
	}
}

func (this *UserPrcoess) PacketFunc(socketid int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			SERVER.GetLog().Println("UserPrcoess PacketFunc", err)
		}
	}()

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

	err := message.UnmarshalText(packet, data)
	if err != nil{
		SERVER.GetLog().Printf("包解析错误2  socket=%d", socketid)
		return true
	}

	packetHead := packet.(message.Packet).GetPacketHead()
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
	bitstream := base.NewBitStream(make([]byte, 1024), 1024)
	if !message.GetMessagePacket(packet, bitstream) {
		SERVER.GetLog().Printf("收到[%s]消息,格式有问题", packetName)
		return true
	}

	if packetHead.DestServerType == int32(message.SERVICE_WORLDSERVER){
		this.SwtichSendToWorld(socketid, packetName, packetHead, bitstream.GetBuffer())
	}else if packetHead.DestServerType == int32(message.SERVICE_ACCOUNTSERVER){
		this.SwtichSendToAccount(socketid, packetName, packetHead, bitstream.GetBuffer())
	}else{
		this.Actor.PacketFunc(socketid,bitstream.GetBuffer())
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