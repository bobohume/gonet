package common

import (
	"fmt"
	"gonet/base"
	"gonet/message"
	"strings"
)

type(
	ServerInfo struct {
		Type int//服务类型编号
		Ip string//服务IP
		Port int//服务端口
		SocketId int//连接句柄
	}

	//集群信息
	ClusterInfo struct {
		Type 		int  		`json:"type"`//服务类型编号
		Ip 			string 		`json:"ip"`//服务IP
		Port 		int  		`json:"port"`//服务端口
		Weight 		int 		`json:"weight"`//权重
	}
)

func (this *ClusterInfo) IpString() string{
	return this.Ip + fmt.Sprintf(":%d", this.Port)
}

func (this *ClusterInfo) String() string{
	return ToServiceString(this.Type)
}

func (this *ClusterInfo) Id() uint32{
	return base.ToHash(this.IpString())
}

func ToServiceString(nType int)string{
	sType := message.SERVICE(nType)
	return strings.ToLower(sType.String())
}