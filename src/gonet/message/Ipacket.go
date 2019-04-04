package message

import (
	"fmt"
	"reflect"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"strings"
	"log"
)

var(
	Packet_CreateFactorStringMap map[string] func()proto.Message
	Packet_CreateFactorMap map[uint32] func()proto.Message
	Packet_CreateFactorInit bool
)

func parseTypeElem(val reflect.Value, packetHead **Ipacket) {
	sType := strings.ToLower(val.Type().String())
	index := strings.Index(sType, ".")
	if index!= -1{
		sType = sType[:index]
	}

	switch sType {
	case "*message":
		if !val.IsNil(){
			value := val.Elem().Interface()
			parseTypeStruct(value, packetHead)
		}
	}
}

func setPacketHead(packetHead **Ipacket, TypeName string, protoVal reflect.Value) bool{
	if TypeName == "PacketHead"{
		*packetHead = protoVal.Interface().(*Ipacket)
	}else{
		return false
	}

	return true
}

func parseTypeStruct(message interface{}, packetHead **Ipacket) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GetPakcetHead", err)
		}
	}()

	protoType := reflect.TypeOf(message)
	protoVal := reflect.ValueOf(message)
	if protoType.Kind() == reflect.Ptr {
		protoType = reflect.TypeOf(message).Elem()
		protoVal = reflect.ValueOf(message).Elem()
	}

	for i := 0; i < protoType.NumField(); i++{
		if !setPacketHead(packetHead, protoType.Field(i).Name, protoVal.Field(i)){
			parseTypeElem(protoVal.Field(i), packetHead)
		}else{
			break
		}
	}
}

func GetPakcetHead(message interface{}) *Ipacket{
	packetHead := BuildPacketHead( 0, 0)
	parseTypeStruct(message, &packetHead)
	return packetHead
}

func BuildPacketHead(id int64, destservertype int) *Ipacket{
	ipacket := &Ipacket{
		Stx:	proto.Int32(Default_Ipacket_Stx),
		DestServerType:	proto.Int32(int32(destservertype)),
		Ckx:	proto.Int32(Default_Ipacket_Ckx),
		Id:	proto.Int64(id),
	}
	return ipacket
}

func GetMessageName(packet proto.Message) string{
	sType := strings.ToLower(proto.MessageName(packet))
	index := strings.Index(sType, ".")
	if index!= -1{
		sType = sType[index+1:]
	}
	return sType
}

func Encode(packet proto.Message) []byte{
	packetId := base.GetMessageCode1(GetMessageName(packet))
	buff,_ := proto.Marshal(packet)
	data := append(base.IntToBytes(int(packetId)), buff...)
	return data
}

func Decode(buff []byte) (uint32, []byte){
	packetId := uint32(base.BytesToInt(buff[0:4]))
	return packetId, buff[4:]
}

func GetMessagePacket(packet proto.Message, bitstream *base.BitStream) bool {
	bitstream.WriteString(GetMessageName(packet))
	bitstream.WriteInt(1, 8)
	{
		sType := strings.ToLower(reflect.ValueOf(packet).Type().String())
		index := strings.Index(sType, ".")
		if index!= -1{
			sType = sType[:index]
		}
		switch sType {
		case "*message":
			bitstream.WriteInt(base.RPC_PB, 8)
			buf, _ :=proto.Marshal(packet)
			nLen := len(buf)
			bitstream.WriteInt(nLen, base.Bit32)
			bitstream.WriteBits(nLen << 3, buf)
		default:
			log.Printf("packet params type not supported", packet, sType)
			return false
		}
	}
	return true
}

func RegisterPacket(packet proto.Message) {
	packetName := GetMessageName(packet)
	packetFunc := func() proto.Message{
		packet := reflect.New(reflect.ValueOf(packet).Elem().Type()).Interface().(proto.Message)
		return packet
	}

	Packet_CreateFactorStringMap[packetName] = packetFunc
	Packet_CreateFactorMap[base.GetMessageCode1(packetName)] = packetFunc
}

func GetPakcet(packetId uint32) proto.Message{
	if !Packet_CreateFactorInit{
		Packet_CreateFactorStringMap = make(map[string] func()proto.Message)
		Packet_CreateFactorMap 		 = make(map[uint32] func()proto.Message)

		//注册消息
		RegisterPacket(&C_A_LoginRequest{})
		RegisterPacket(&C_A_RegisterRequest{})
		RegisterPacket(&C_G_LogoutResponse{})
		RegisterPacket(&C_W_CreatePlayerRequest{})
		RegisterPacket(&C_W_Game_LoginRequset{})
		RegisterPacket(&C_W_LoginCopyMap{})
		RegisterPacket(&C_W_Move{})
		RegisterPacket(&C_W_ChatMessage{})
		// test for client
		RegisterPacket(&W_C_SelectPlayerResponse{})
		RegisterPacket(&W_C_CreatePlayerResponse{})
		RegisterPacket(&W_C_LoginMap{})
		RegisterPacket(&W_C_Move{})
		RegisterPacket(&W_C_ADD_SIMOBJ{})
		RegisterPacket(&A_C_LoginRequest{})
		RegisterPacket(&A_C_RegisterResponse{})
		RegisterPacket(&W_C_ChatMessage{})
		Packet_CreateFactorInit = true
	}

	packetFunc,exist := Packet_CreateFactorMap[packetId]
	if exist{
		return packetFunc()
	}

	return nil;
}

func GetPakcetByName(packetName string) proto.Message{
	return GetPakcet(base.GetMessageCode1(packetName))
}

func UnmarshalText(packet proto.Message, packetBuf []byte) error{
	return proto.Unmarshal(packetBuf, packet)
}