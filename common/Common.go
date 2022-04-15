package common

import (
	"fmt"
	"gonet/base"
	"gonet/rpc"
	"strings"
)

type (
	//集群信息
	ClusterInfo rpc.ClusterInfo

	IClusterInfo interface {
		Id() uint32
		String() string
		ServiceType() rpc.SERVICE
		IpString() string
	}

	StubMailBox struct {
		rpc.StubMailBox
	}
)

func (this *ClusterInfo) IpString() string {
	return fmt.Sprintf("%s:%d", this.Ip, this.Port)
}

func (this *ClusterInfo) String() string {
	return strings.ToLower(this.Type.String())
}

func (this *ClusterInfo) Id() uint32 {
	return base.ToHash(this.IpString())
}

func (this *ClusterInfo) ServiceType() rpc.SERVICE {
	return this.Type
}

func (this *StubMailBox) StubName() string {
	return this.StubType.String()
}

func (this *StubMailBox) Key() string {
	return fmt.Sprintf("%s/%d", this.StubType.String(), this.Id)
}
