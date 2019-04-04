package json3

type C_A_LoginRequest struct {
	MessageBase
	AccountName      string  `protobuf:"bytes,2,req,name=AccountName" json:"AccountName,omitempty"`
	BuildNo          string  `protobuf:"bytes,3,req,name=BuildNo" json:"BuildNo,omitempty"`
	SocketId         int32   `protobuf:"varint,4,req,name=SocketId" json:"SocketId,omitempty"`
}
