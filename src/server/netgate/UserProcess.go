package netgate

import (
	"actor"
	"base"
	"github.com/golang/protobuf/proto"
	"message"
)

type(
	UserPrcoess struct{
		actor.Actor
	}

	IUserPrcoess interface {
		actor.IActor

		CheckClient(int, string, interface{})bool
		SwtichSendToWorld(int, string, interface{}, []byte)
		SwtichSendToAccount(int, string, interface{}, []byte)
	}
)

func (this *UserPrcoess)CheckClient(sockId int, packetName string, packet interface{}) bool{
	packetHead := packet.(*message.Ipacket)
	if packetHead != nil{
		if IsCheckClient(packetName){
			return  true
		}

		accountId := SERVER.GetPlayerMgr().GetSocketAccount(sockId)
		if accountId <= 0 || accountId != int(*packetHead.Id) {
			SERVER.GetLog().Fatalf("Old socket communication or viciousness[%d].", sockId)
			return false
		}
		return  true
	}
	return  false
}

func (this *UserPrcoess)SwtichSendToWorld(socketId int, packetName string, packet interface{}, buff []byte){
	if this.CheckClient(socketId, packetName, packet) == true{
		buff = base.SetTcpEnd(buff)
		SERVER.GetWorldSocket().Send(buff)
	}
}

func (this *UserPrcoess)SwtichSendToAccount(socketId int, packetName string, packet interface{}, buff []byte){
	if this.CheckClient(socketId, packetName, packet) == true{
		buff = base.SetTcpEnd(buff)
		SERVER.GetAccountSocket().Send(buff)
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
		SERVER.GetLog().Printf("包解析错误1  socket=%d", socketid)
		return true
	}

	err := proto.Unmarshal(data, packet)
	if err != nil{
		SERVER.GetLog().Printf("包解析错误2  socket=%d", socketid)
		return true
	}

	packetHead := message.GetPakcetHead(packet)
	if packetHead == nil || *packetHead.Ckx != message.Default_Ipacket_Ckx || *packetHead.Stx != message.Default_Ipacket_Stx {
		SERVER.GetLog().Printf("(A)致命的越界包,已经被忽略 socket=%d", socketid)
		return true
	}

	packetName := message.GetMessageName(packet)
	if packetName  == base.ToLower("C_A_LoginRequest") {
		packet.(*message.C_A_LoginRequest).SocketId = proto.Int32(int32(socketid))
	}else if packetName  == base.ToLower("C_A_RegisterRequest") {
		packet.(*message.C_A_RegisterRequest).SocketId = proto.Int32(int32(socketid))
	}
	/*if *packetHead.Message  == base.GetMessageCode1("C_A_LoginRequest"){
		packet.(*message.C_A_LoginRequest).SocketId = proto.Int32(int32(socketid))
	}else if(*packetHead.Message  == base.GetMessageCode1("C_A_RegisterRequest")){
		packet.(*message.C_A_RegisterRequest).SocketId = proto.Int32(int32(socketid))
	}*/

	/*if !IsValidClientMsg(*packetHead.Message){
		SERVER.GetLog().Printf("收到未注册[%s]消息", *packetHead.Message)
		return true
	}*/

	//解析整个包
	bitstream := base.NewBitStream(make([]byte, 1024), 1024)
	if !message.GetProtoBufPacket(packet, bitstream) {
		SERVER.GetLog().Printf("收到[%s]消息,格式有问题", packetName)
		return true
	}

	if *packetHead.DestServerType == int32(message.SERVICE_WORLDSERVER){
		this.SwtichSendToWorld(socketid, packetName, packetHead, bitstream.GetBuffer())
	}else if *packetHead.DestServerType == int32(message.SERVICE_ACCOUNTSERVER){
		this.SwtichSendToAccount(socketid, packetName, packetHead, bitstream.GetBuffer())
	}else{
		this.Actor.PacketFunc(socketid,bitstream.GetBuffer())
	}

	return true
}

/*func (this *UserEventPrcoess) PacketFunc(socketid int, buff []byte) bool{
	defer func() {
		if err := recover(); err != nil {
			SERVER.GetLog().Println("UserEventPrcoess PacketFunc", err)
		}
	}()

	//解析报头
	packetHead := &message.Ipacket{}
	err := proto.Unmarshal(buff, packetHead)
	if err != nil{
		SERVER.GetLog().Printf("(A)包头解析错误  socket=%d", socketid)
		return true
	}

	if *packetHead.Ckx != message.Default_Ipacket_Ckx || *packetHead.Stx != message.Default_Ipacket_Stx {
		SERVER.GetLog().Printf("(A)致命的越界包,已经被忽略 socket=%d", socketid)
		return true
	}

	if !IsValidClientMsg(*packetHead.Message){
		SERVER.GetLog().Printf("收到未注册[%s]消息", *packetHead.Message)
		return true
	}

	//解析整个包
	packet := s_clientMsgFilters[*packetHead.Message]()
	proto.Unmarshal(buff, packet)
	bitstream := base.NewBitStream(make([]byte, 1024), 1024)
	bitstream.WriteString(*packetHead.Message)
	if !base.ProtoToBitStream(packet, bitstream) {
		SERVER.GetLog().Printf("收到[%s]消息,格式有问题", *packetHead.Message)
		return true
	}

	if *packetHead.DestServerType == int32(message.SERVICE_WORLDSERVER){
		this.SwtichSendToWorld(socketid, packetHead, bitstream.GetBuffer())
	}else if *packetHead.DestServerType == int32(message.SERVICE_ACCOUNTSERVER){
		this.SwtichSendToAccount(socketid, packetHead, bitstream.GetBuffer())
	}else{
		this.Actor.PacketFunc(socketid,bitstream.GetBuffer())
	}

	return true
}*/

func (this *UserPrcoess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("C_G_LoginRequest", func(accountId int, UID int) {
		//SERVER.GetPlayerMgr().SendMsg("ADD_ACCOUNT", caller.SocketId, accountId, UID)
	})

	this.RegisterCall("C_G_LogoutRequest", func(accountId int, UID int){
		SERVER.GetLog().Printf("logout Socket:%d Account:%d UID:%d ",this.GetSocketId(), accountId,UID )
		SERVER.GetPlayerMgr().SendMsg("DEL_ACCOUNT", this.GetSocketId())
		SendToClient(this.GetSocketId(), &message.C_G_LogoutResponse{PacketHead:message.BuildPacketHead( 0, 0)})
	})

	this.Actor.Start()
}


