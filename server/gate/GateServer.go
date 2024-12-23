package gate

import (
	"gonet/actor"
	"gonet/base"
	"gonet/base/cluster"
	"gonet/base/conf"
	"gonet/network"
	"gonet/rpc"
)

type (
	ServerMgr struct {
		service   *network.ServerSocket
		isInited  bool
		playerMgr *PlayerMgr
		cluster   *cluster.Cluster
		stats     *StatsPrcoess
	}

	IServerMgr interface {
		Init() bool
		GetServer() *network.ServerSocket
		GetCluster() *cluster.Service
		GetPlayerMgr() *PlayerMgr
		OnServerStart()
	}

	Config struct {
		conf.Server `yaml:"gate"`
		conf.Etcd   `yaml:"etcd"`
		conf.Nats   `yaml:"nats"`
		conf.Raft   `yaml:"raft"`
		conf.Stub   `yaml:"stub"`
	}
)

var (
	CONF   Config
	SERVER ServerMgr
)

func (s *ServerMgr) GetServer() *network.ServerSocket {
	return s.service
}

func (s *ServerMgr) GetCluster() *cluster.Cluster {
	return s.cluster
}

func (s *ServerMgr) GetPlayerMgr() *PlayerMgr {
	return s.playerMgr
}

func (s *ServerMgr) GetStats() *StatsPrcoess {
	return s.stats
}

func (s *ServerMgr) Init() bool {
	if s.isInited {
		return true
	}

	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func() {
		base.LOG.Println("**********************************************************")
		base.LOG.Printf("\tGATE Version:\t%s", base.BUILD_NO)
		base.LOG.Printf("\tGATE IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		base.LOG.Println("**********************************************************")
	}
	ShowMessage()

	//初始化socket
	s.service = new(network.ServerSocket)
	s.service.Init(CONF.Server.Ip, CONF.Server.Port)
	s.service.SetMaxPacketLen(base.MAX_CLIENT_PACKET)
	s.service.SetConnectType(network.CLIENT_CONNECT)
	//s.service.Start()
	packet := new(UserPrcoess)
	packet.Init()
	s.service.BindPacketFunc(packet.PacketFunc)
	s.service.Start()

	//websocket
	/*s.service = new(network.WebSocket)
	s.service.Init(CONF.Server.Ip, CONF.Server.Port)
	s.service.SetConnectType(network.CLIENT_CONNECT)
	//s.service.Start()
	packet := new(UserPrcoess)
	packet.Init()
	s.service.BindPacketFunc(packet.PacketFunc)
	s.service.Start()*/
	//注册到集群服务器

	var packet1 EventProcess
	packet1.Init()

	stats := new(StatsPrcoess)
	stats.Init()

	s.cluster = new(cluster.Cluster)
	s.cluster.InitCluster(&rpc.ClusterInfo{Type: rpc.SERVICE_GATE, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)},
		CONF.Etcd.Endpoints, CONF.Nats.Endpoints, cluster.WithStubMailBoxEtcd(CONF.Raft.Endpoints, &CONF.Stub))
	s.cluster.BindPacketFunc(actor.MGR.PacketFunc)
	s.cluster.BindPacketFunc(DispatchPacket)

	//初始玩家管理
	s.playerMgr = new(PlayerMgr)
	s.playerMgr.Init()
	return false
}

func (s *ServerMgr) OnServerStart() {
	s.service.Start()
}
