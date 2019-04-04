package json3

import (
	"encoding/json"
	"gonet/base"
	"log"
	"reflect"
	"strings"
)

type SERVICE int32

const (
	SERVICE_NONE          SERVICE = 0
	SERVICE_CLIENT        SERVICE = 1
	SERVICE_GATESERVER    SERVICE = 2
	SERVICE_ACCOUNTSERVER SERVICE = 3
	SERVICE_WORLDSERVER   SERVICE = 4
	SERVICE_MONITORSERVER SERVICE = 5
)

type CHAT int32
const (
	CHAT_MSG_TYPE_WORLD   CHAT = 0
	CHAT_MSG_TYPE_PRIVATE CHAT = 1
	CHAT_MSG_TYPE_ORG     CHAT = 2
	CHAT_MSG_TYPE_COUNT   CHAT = 3
)

const Default_Ipacket_Stx int32 = 39
const Default_Ipacket_Ckx int32 = 114

type(
	Ipacket struct {
		Stx              int32 `json:"Stx,omitempty"`
		DestServerType   int32 `json:"DestServerType,omitempty"`
		Ckx              int32 `json:"Ckx,omitempty"`
		Id               int64 `json:"Id,omitempty"`
	}

	MessageBase struct {
		PacketHead       Ipacket `json:"PacketHead"`
		MessageName 	 string  `json:"PacketName"`
	}

	Message interface {
		Name() string
		//Reset()
		//String() string
		//ProtoMessage()
		SetName(string)
		Ipacket() *Ipacket
	}
)

func (this *MessageBase) Name() string{
	return this.MessageName
}

/*func (this *MessageBase) Reset(){
}

func (this *MessageBase) ProtoMessage(){
}

func (this *MessageBase) String()string{
	return "json"
}*/

func (this *MessageBase) SetName(name string){
	this.MessageName = name
}

func (this *MessageBase) Ipacket() *Ipacket{
	return &this.PacketHead
}

var(
	Packet_CreateFactorStringMap map[string] func()Message
	Packet_CreateFactorMap map[uint32] func()Message
	Packet_CreateFactorInit bool
)

func GetPakcetHead(packet Message) *Ipacket{
	return packet.Ipacket()
}

func BuildMessageBase(id int64, destservertype int, packetName string) *MessageBase{
	packetName = strings.ToLower(packetName)
	return &MessageBase{*BuildPacketHead(id, destservertype), packetName}
}

func BuildPacketHead(id int64, destservertype int) *Ipacket{
	ipacket := &Ipacket{
		Stx:	Default_Ipacket_Stx,
		DestServerType:	int32(destservertype),
		Ckx:	Default_Ipacket_Ckx,
		Id:	  id,
	}
	return ipacket
}

func GetMessageName(packet Message) string{
	return 	strings.ToLower(packet.Name())
}

func Encode(packet Message) []byte{
	packetId := base.GetMessageCode1(GetMessageName(packet))
	buff,_ := json.Marshal(packet)
	data := append(base.IntToBytes(int(packetId)), buff...)
	return data
}

func Decode(buff []byte) (uint32, []byte){
	packetId := uint32(base.BytesToInt(buff[0:4]))
	return packetId, buff[4:]
}

func GetMessagePacket(packet Message, bitstream *base.BitStream) bool {
	bitstream.WriteString(GetMessageName(packet))
	bitstream.WriteInt(1, 8)
	{
		sType := strings.ToLower(reflect.ValueOf(packet).Type().String())
		index := strings.Index(sType, ".")
		if index!= -1{
			sType = sType[:index]
		}
		switch sType {
		case "*json3":
			bitstream.WriteInt(base.RPC_JSON, 8)
			buf, _ := json.Marshal(packet)
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

func RegisterPacket(packet Message) {
	packetName := GetMessageName(packet)
	packetFunc := func() Message{
		packet := reflect.New(reflect.ValueOf(packet).Elem().Type()).Interface().(Message)
		packet.SetName(packetName)
		return packet
	}

	Packet_CreateFactorStringMap[packetName] = packetFunc
	Packet_CreateFactorMap[base.GetMessageCode1(packetName)] = packetFunc
}

func GetPakcet(packetId uint32) Message{
	if !Packet_CreateFactorInit{
		Packet_CreateFactorStringMap = make(map[string] func()Message)
		Packet_CreateFactorMap 		 = make(map[uint32] func()Message)

		//注册消息
		RegisterPacket(&C_A_LoginRequest{MessageBase:MessageBase{Ipacket{}, "C_A_LoginRequest_json"},})
		// test for client
		Packet_CreateFactorInit = true
	}

	packetFunc,exist := Packet_CreateFactorMap[packetId]
	if exist{
		return packetFunc()
	}

	return nil;
}

func GetPakcetByName(packetName string) Message{
	return GetPakcet(base.GetMessageCode1(packetName))
}

func UnmarshalText(packet Message, packetBuf []byte) error{
	return json.Unmarshal(packetBuf, packet)
}
