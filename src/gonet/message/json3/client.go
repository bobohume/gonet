package json3

type C_A_LoginRequest struct {
	MessageBase
	AccountName      string  `protobuf:"bytes,2,req,name=AccountName" json:"AccountName,omitempty"`
	BuildNo          string  `protobuf:"bytes,3,req,name=BuildNo" json:"BuildNo,omitempty"`
	SocketId         int32   `protobuf:"varint,4,req,name=SocketId" json:"SocketId,omitempty"`
}

type C_A_RegisterRequest struct {
	MessageBase
	AccountName      string  `protobuf:"bytes,2,req,name=AccountName" json:"AccountName,omitempty"`
	SocketId         int32   `protobuf:"varint,3,req,name=SocketId" json:"SocketId,omitempty"`
}

type C_W_CreatePlayerRequest struct {
	MessageBase
	PlayerName       string  `protobuf:"bytes,2,req,name=PlayerName" json:"PlayerName,omitempty"`
	Sex              int32   `protobuf:"varint,3,req,name=Sex" json:"Sex,omitempty"`
}

type C_W_Game_LoginRequset struct {
	MessageBase
	PlayerId         int64   `protobuf:"varint,2,req,name=PlayerId" json:"PlayerId,omitempty"`
}

type C_G_LogoutResponse struct {
	MessageBase
}

type C_W_ChatMessage struct {
	MessageBase
	Sender           int64   `protobuf:"varint,2,req,name=Sender" json:"Sender,omitempty"`
	Recver           int64   `protobuf:"varint,3,req,name=Recver" json:"Recver,omitempty"`
	MessageType      int32   `protobuf:"varint,4,req,name=MessageType" json:"MessageType,omitempty"`
	Message          string  `protobuf:"bytes,5,req,name=Message" json:"Message,omitempty"`
}

type W_C_ChatMessage struct {
	MessageBase
	Sender           int64   `protobuf:"varint,2,req,name=Sender" json:"Sender,omitempty"`
	SenderName       string  `protobuf:"bytes,3,req,name=SenderName" json:"SenderName,omitempty"`
	Recver           int64   `protobuf:"varint,4,req,name=Recver" json:"Recver,omitempty"`
	RecverName       string  `protobuf:"bytes,5,req,name=RecverName" json:"RecverName,omitempty"`
	MessageType      int32   `protobuf:"varint,6,req,name=MessageType" json:"MessageType,omitempty"`
	Message          string  `protobuf:"bytes,7,req,name=Message" json:"Message,omitempty"`
}

type PlayerData struct {
	PlayerID         int64  `protobuf:"varint,1,req,name=PlayerID" json:"PlayerID,omitempty"`
	PlayerName       string `protobuf:"bytes,2,req,name=PlayerName" json:"PlayerName,omitempty"`
	PlayerGold       int32  `protobuf:"varint,3,req,name=PlayerGold" json:"PlayerGold,omitempty"`
}

type W_C_CreatePlayerResponse struct {
	MessageBase
	Error            int32   `protobuf:"varint,2,req,name=Error" json:"Error,omitempty"`
	PlayerId         int64   `protobuf:"varint,3,req,name=PlayerId" json:"PlayerId,omitempty"`
}

type W_C_SelectPlayerResponse struct {
	MessageBase
	AccountId        int64        `protobuf:"varint,2,req,name=AccountId" json:"AccountId,omitempty"`
	PlayerData       []*PlayerData `protobuf:"bytes,3,rep,name=PlayerData" json:"PlayerData,omitempty"`
}

type A_C_LoginRequest struct {
	MessageBase
	Error            int32   `protobuf:"varint,2,req,name=Error" json:"Error,omitempty"`
	SocketId         int32   `protobuf:"varint,3,req,name=SocketId" json:"SocketId,omitempty"`
	AccountName      string  `protobuf:"bytes,4,req,name=AccountName" json:"AccountName,omitempty"`
}