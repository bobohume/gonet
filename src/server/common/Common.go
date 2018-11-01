package common

import "base"

type(
	ServerInfo struct {
		Type int//服务类型编号
		Ip string//服务IP
		Port int//服务端口
		SocketId int//连接句柄
	}
)


func (this *ServerInfo) ReadData(b *base.BitStream){
	base.ReadData(this, b)
}

func (this *ServerInfo) WriteData(b *base.BitStream){
	base.WriteData(this, b)
}