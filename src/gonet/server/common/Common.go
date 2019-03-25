package common

type(
	ServerInfo struct {
		Type int//服务类型编号
		Ip string//服务IP
		Port int//服务端口
		SocketId int//连接句柄
	}
)