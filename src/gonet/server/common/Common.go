package common

import (
	"fmt"
	"gonet/base"
	"gonet/message"
	"strings"
)

type(
	//集群信息
	ClusterInfo message.ClusterInfo
)

func (this *ClusterInfo) IpString() string{
	return fmt.Sprintf("%s:%d", this.Ip, this.Port)
}

func (this *ClusterInfo) String() string{
	return  strings.ToLower(this.Type.String())
}

func (this *ClusterInfo) Id() uint32{
	return base.ToHash(this.IpString())
}

func (this *ClusterInfo) ServiceType() message.SERVICE{
	return this.Type
}