package message

import (
	"gonet/base"
	"gonet/rpc"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	Packet_CreateFactorStringMap map[string]*PacketRoute
	Packet_CreateFactorMap       map[uint32]*PacketRoute
	Packet_CrcNamesMap           map[uint32]string
)

const (
	Default_Ipacket_Stx int32 = 0x27
	Default_Ipacket_Ckx int32 = 0x72
)

type (
	//获取包头
	Packet interface {
		GetPacketHead() *Ipacket
	}

	PacketRoute struct {
		Func func() proto.Message
		FuncName string
	}
)

func BuildPacketHead(id int64, destservertype rpc.SERVICE) *Ipacket {
	ipacket := &Ipacket{
		Stx:            Default_Ipacket_Stx,
		DestServerType: SERVICE(destservertype),
		Ckx:            Default_Ipacket_Ckx,
		Id:             id,
	}
	return ipacket
}

func GetMessageName(packet proto.Message) string {
	sType := strings.ToLower(proto.MessageName(packet))
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}

func Encode(packet proto.Message) []byte {
	packetId := base.GetMessageCode1(GetMessageName(packet))
	buff, _ := proto.Marshal(packet)
	data := append(base.IntToBytes(int(packetId)), buff...)
	return data
}

func Decode(buff []byte) (uint32, []byte) {
	packetId := uint32(base.BytesToInt(buff[0:4]))
	return packetId, buff[4:]
}

func RegisterPacket(packet proto.Message, funcName string) {
	packetName := GetMessageName(packet)
	val := reflect.ValueOf(packet).Elem()
	packetFunc := func() proto.Message {
		packet := reflect.New(val.Type())
		packet.Elem().Field(3).Set(val.Field(3))
		//packet.Elem().Set(val)
		return packet.Interface().(proto.Message)
	}

	packetRoute := &PacketRoute{packetFunc, funcName}
	Packet_CreateFactorStringMap[packetName] = packetRoute
	Packet_CreateFactorMap[base.GetMessageCode1(packetName)] = packetRoute
}

func GetPakcetRoute(packetId uint32) *PacketRoute{
	packetRoute, exist := Packet_CreateFactorMap[packetId]
	if exist {
		return packetRoute
	}

	return nil
}

func GetPakcetName(packetId uint32) string {
	return Packet_CrcNamesMap[packetId]
}

func UnmarshalText(packet proto.Message, packetBuf []byte) error {
	return proto.Unmarshal(packetBuf, packet)
}

func init() {
	Packet_CreateFactorStringMap = make(map[string]*PacketRoute)
	Packet_CreateFactorMap = make(map[uint32]*PacketRoute)
	Packet_CrcNamesMap = make(map[uint32]string)
}

//统计crc对应string
func initCrcNames() {
	protoFiles := []protoreflect.MessageDescriptors{}
	protoFiles = append(protoFiles, File_message_proto.Messages())
	protoFiles = append(protoFiles, File_client_proto.Messages())
	protoFiles = append(protoFiles, File_game_proto.Messages())
	for _, v := range protoFiles {
		for i := 0; i < v.Len(); i++ {
			packetName := strings.ToLower(string(v.Get(i).Name()))
			crcVal := base.GetMessageCode1(packetName)
			Packet_CrcNamesMap[crcVal] = packetName
		}
	}
}

//网关防火墙
func Init() {
	initCrcNames()
	//注册消息
	//PacketHead 中的 DestServerType 决定转发到那个服务器
	RegisterPacket(&LoginAccountRequest{}, "gate<-UserPrcoess.LoginAccountRequest")
	RegisterPacket(&LoginPlayerRequset{}, "gate<-UserPrcoess.LoginPlayerRequset")
	RegisterPacket(&CreatePlayerRequest{}, "gm<-AccountMgr.CreatePlayerRequest")
	RegisterPacket(&ChatMessageRequest{}, "gm<-ChatMgr.ChatMessageRequest")

	RegisterPacket(&C_Z_LoginCopyMap{}, "zone<-MapMgr.C_Z_LoginCopyMap")
	RegisterPacket(&C_Z_Move{}, "zone<-Map.C_Z_Move")
	RegisterPacket(&C_Z_Skill{}, "zone<-Map.C_Z_Skill")
}

//client消息回调
func InitClient() {
	initCrcNames()
	//注册消息
	RegisterPacket(&Z_C_LoginMap{}, "client<-EventProcess.Z_C_LoginMap")
	RegisterPacket(&Z_C_ENTITY{}, "client<-EventProcess.Z_C_ENTITY")
	RegisterPacket(&ChatMessageResponse{}, "client<-EventProcess.ChatMessageResponse")
	RegisterPacket(&LoginAccountResponse{}, "client<-EventProcess.LoginAccountResponse")
	RegisterPacket(&SelectPlayerResponse{}, "client<-EventProcess.SelectPlayerResponse")
}
