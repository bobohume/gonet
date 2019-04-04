package world

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"sort"
)

type(
	CoverList struct {
		base.Vector
	}

	CoverInfo struct {
		srcV int
		destV int
		pFunc CoverBlobFunc
	}

	//blob转换函数
	CoverBlobFunc func(v1 int, src proto.Message, v2 int, dest proto.Message) (int, proto.Message)
	//blob管理
	BlobMgr struct {
		m_BlobMap map[int] func()proto.Message
		m_BlobCoverList CoverList
	}

	IBlobMgr interface {
		Init()//初始化
		RegisterBlob(int, proto.Message)//注册blob，初始信息等
		RegisterCoverBlob(v1 int, v2 int, pFunc CoverBlobFunc)//这里是blob转化
		GetBlob(version int, blob proto.Message) (int, proto.Message)//获取blob
	}
)

func (this *BlobMgr) Init(){
	this.m_BlobMap = make(map[int] func()proto.Message)
	//注册消息
	//RegisterPacket(&C_A_LoginRequest{})
	//CoverBlob(1001, 1002, func(v1 int, src proto.Message, v2 int, dest proto.Message){
	//
	//})
	sort.Sort(&this.m_BlobCoverList)
}

//注册blob，初始化
func (this *BlobMgr) RegisterBlob(version int, packet proto.Message) {
	packetFunc := func() proto.Message{
		packet := proto.Clone(packet)
		return packet
	}

	this.m_BlobMap[version] = packetFunc
}

//转化blob
func (this *BlobMgr) RegisterCoverBlob(v1 int, v2 int, pFunc CoverBlobFunc){
	if v1 >= v2{
		panic(fmt.Sprintf("CoverBlob verison [%d] to [%d]", v1, v2))
	}
	this.m_BlobCoverList.Push_back(&CoverInfo{v1, v2, pFunc})
}

//blob转换
func (this *BlobMgr) GetBlob(version int, blob proto.Message) (int, proto.Message){
	for _, v := range this.m_BlobCoverList.Array(){
		pCover := v.(*CoverInfo)
		if pCover.srcV <= version &&  pCover.destV > version{
			version, blob = pCover.pFunc(pCover.srcV, blob, pCover.destV, this.m_BlobMap[pCover.destV]())
		}
	}
	return version, blob
}

//sort interface
func (t CoverList) Less(i, j int) bool{
	return t.Get(i).(*CoverInfo).srcV < t.Get(i).(*CoverInfo).srcV
}

//test
/*db2 := &db2.BlobMgr{}
db2.Init()
db2.RegisterBlob(1, &message.PlayerData_1{PlayerGold:proto.Int(1000)})
db2.RegisterBlob(2, &message.PlayerData_2{PlayerGold:proto.Int64(1000)})
db2.RegisterBlob(5, &message.PlayerData_5{PlayerGold:proto.Int64(1000)})
db2.RegisterCoverBlob(1, 2, func(v1 int, src proto.Message, v2 int, dest proto.Message) (int,proto.Message) {
	src1  := src.(*message.PlayerData_1)
	dest1 := dest.(*message.PlayerData_2)
	dest1.PlayerGold = proto.Int64(int64(src1.GetPlayerGold()))
	dest1.PlayerGold1 = proto.Int32(100)
	base.Copy(dest1, src1)
	return v2, dest1
})
db2.RegisterCoverBlob(2, 5, func(v1 int, src proto.Message, v2 int, dest proto.Message) (int,proto.Message) {
	src1  := src.(*message.PlayerData_2)
	dest1 := dest.(*message.PlayerData_5)
	dest1.PlayerGold = proto.Int64(int64(src1.GetPlayerGold()))
	dest1.PlayerGold1 = proto.Int32(100)
	dest1.PlayerGold2 = proto.Int32(200)
	base.Copy(dest1, src1)
	return v2, dest1
})

db2.GetBlob(1, &message.PlayerData_1{PlayerGold:proto.Int(1000), PlayerName:proto.String("tet"), PlayerID:proto.Int64(base.UUID.UUID())})
*/