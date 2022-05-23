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

func (c *ClusterInfo) IpString() string {
	return fmt.Sprintf("%s:%d", c.Ip, c.Port)
}

func (c *ClusterInfo) String() string {
	return strings.ToLower(c.Type.String())
}

func (c *ClusterInfo) Id() uint32 {
	return base.ToHash(c.IpString())
}

func (c *ClusterInfo) ServiceType() rpc.SERVICE {
	return c.Type
}

func (s *StubMailBox) StubName() string {
	return s.StubType.String()
}

func (s *StubMailBox) Key() string {
	return fmt.Sprintf("%s/%d", s.StubType.String(), s.Id)
}
