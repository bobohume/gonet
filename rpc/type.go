package rpc

import (
	"fmt"
	"gonet/base"
	"reflect"
	"strings"
)

type (
	ICluster interface {
		SendMsg(head RpcHead, funcName string, params ...interface{})
		Call(parmas ...interface{})
		Id() uint32
	}
)

var MGR ICluster

func Route(head *RpcHead, funcName string) string {
	serverArgs := strings.Split(funcName, "<-")
	if len(serverArgs) == 2 {
		switch strings.ToLower(serverArgs[0]) {
		case "client":
			head.DestServerType = SERVICE_CLIENT
		case "gate":
			head.DestServerType = SERVICE_GATE
		case "gm":
			head.DestServerType = SERVICE_GM
		case "game":
			head.DestServerType = SERVICE_GAME
		case "zone":
			head.DestServerType = SERVICE_ZONE
		case "db":
			head.DestServerType = SERVICE_DB
		}
		funcName = serverArgs[1]
	}

	actorArgs := strings.Split(funcName, ".")
	if len(actorArgs) == 2 {
		head.ActorName = actorArgs[0]
		funcName = actorArgs[1]
	}

	return funcName
}

type (
	//集群信息

	IClusterInfo interface {
		Id() uint32
		ServiceName() string
		ServiceType() SERVICE
		IpString() string
	}
)

func (c *ClusterInfo) IpString() string {
	return fmt.Sprintf("%s:%d", c.Ip, c.Port)
}

func (c *ClusterInfo) ServiceName() string {
	return strings.ToLower(c.Type.String())
}

func (c *ClusterInfo) Id() uint32 {
	return base.ToHash(c.IpString())
}

func (c *ClusterInfo) ServiceType() SERVICE {
	return c.Type
}

func (s *StubMailBox) StubName() string {
	return s.StubType.String()
}

func (s *StubMailBox) Key() string {
	return fmt.Sprintf("%s/%d", s.StubType.String(), s.Id)
}

var GCall = reflect.Value{}

func InitMGR(mgr ICluster) {
	MGR = mgr
	GCall = reflect.ValueOf(mgr.Call)
}
