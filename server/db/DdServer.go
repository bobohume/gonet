package db

import (
	"gonet/actor"
	"gonet/base"
	"gonet/base/ini"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/network"
	"gonet/orm"
	"gonet/rpc"
	"log"
)

type (
	ServerMgr struct {
		m_pService *network.ServerSocket
		m_config   ini.Config
	}

	IServerMgr interface {
		Init(string) bool
		InitDB() bool
		GetServer() *network.ServerSocket
	}

	Config struct {
		common.Server `yaml:"db"`
		common.Db     `yaml:"DB"`
		common.Etcd   `yaml:"etcd"`
		common.Nats   `yaml:"nats"`
		common.Raft   `yaml:"raft"`
		common.Stub   `yaml:"stub"`
	}
)

var (
	CONF   Config
	SERVER ServerMgr
)

func (this *ServerMgr) Init() bool {
	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func() {
		base.LOG.Println("**********************************************************")
		base.LOG.Printf("\tDB Version:\t%s", base.BUILD_NO)
		base.LOG.Printf("\tDb IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		base.LOG.Printf("\tActorDBServer(LAN):\t%s", CONF.Db.Ip)
		base.LOG.Printf("\tActorDBName:\t\t%s", CONF.Db.Name)
		base.LOG.Println("**********************************************************")
	}
	ShowMessage()

	base.LOG.Println("正在初始化数据库连接...")
	if this.InitDB() {
		base.LOG.Printf("[%s]数据库连接是失败...", CONF.Db.Name)
		log.Fatalf("[%s]数据库连接是失败...", CONF.Db.Name)
		return false
	}
	base.LOG.Printf("[%s]数据库初始化成功!", CONF.Db.Name)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	this.m_pService.Init(CONF.Server.Ip, CONF.Server.Port)
	this.m_pService.Start()

	//本身db集群管理
	cluster.MGR.InitCluster(&common.ClusterInfo{Type: rpc.SERVICE_DB, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)}, CONF.Etcd.Endpoints, CONF.Nats.Endpoints,
		cluster.WithStubMailBoxEtcd(CONF.Raft.Endpoints, &CONF.Stub))
	cluster.MGR.BindPacketFunc(actor.MGR.PacketFunc)

	PLAYERSAVEMGR.Init()

	return false
}

func (this *ServerMgr) InitDB() bool {
	return orm.OpenDB(CONF.Db) != nil
}

func (this *ServerMgr) GetServer() *network.ServerSocket {
	return this.m_pService
}
