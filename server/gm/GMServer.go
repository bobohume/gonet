package gm

import (
	"gonet/actor"
	"gonet/base"
	"gonet/base/ini"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/network"
	"gonet/orm"
	"gonet/rpc"
	"gonet/server/message"
	"log"

	"github.com/golang/protobuf/proto"
)

type (
	ServerMgr struct {
		m_pService  *network.ServerSocket
		m_Inited    bool
		m_config    ini.Config
		m_SnowFlake *cluster.Snowflake
	}

	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetServer() *network.ServerSocket
	}

	Config struct {
		common.Server    `yaml:"gm"`
		common.Db        `yaml:"DB"`
		common.Etcd      `yaml:"etcd"`
		common.SnowFlake `yaml:"snowflake"`
		common.Raft      `yaml:"raft"`
		common.Nats      `yaml:"nats"`
		common.Stub      `yaml:"stub"`
	}
)

var (
	CONF   Config
	SERVER ServerMgr
)

func (this *ServerMgr) Init() bool {
	if this.m_Inited {
		return true
	}

	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func() {
		base.LOG.Println("**********************************************************")
		base.LOG.Printf("\tGM Version:\t%s", base.BUILD_NO)
		base.LOG.Printf("\tGM IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		base.LOG.Printf("\tDBServer(LAN):\t%s", CONF.Db.Ip)
		base.LOG.Printf("\tDBName:\t\t%s", CONF.Db.Name)
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

	//本身账号集群管理
	cluster.MGR.InitCluster(&common.ClusterInfo{Type: rpc.SERVICE_GM, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)},
		CONF.Etcd.Endpoints, CONF.Nats.Endpoints, cluster.WithMailBoxEtcd(CONF.Raft.Endpoints), cluster.WithStubMailBoxEtcd(CONF.Raft.Endpoints, &CONF.Stub))
	cluster.MGR.BindPacketFunc(actor.MGR.PacketFunc)

	//snowflake
	this.m_SnowFlake = cluster.NewSnowflake(CONF.SnowFlake.Endpoints)

	SIMPLEMGR.Init()
	return false
}

func (this *ServerMgr) InitDB() bool {
	return orm.OpenDB(CONF.Db) != nil
}

func (this *ServerMgr) GetServer() *network.ServerSocket {
	return this.m_pService
}

//发送game
func SendToGame(ClusterId uint32, funcName string, params ...interface{}) {
	head := rpc.RpcHead{ClusterId: ClusterId, DestServerType: rpc.SERVICE_GAME, SrcClusterId: cluster.MGR.Id()}
	cluster.MGR.SendMsg(head, funcName, params...)
}

//广播game
func BoardCastToGame(funcName string, params ...interface{}) {
	head := rpc.RpcHead{DestServerType: rpc.SERVICE_GAME, SendType: rpc.SEND_BOARD_CAST, SrcClusterId: cluster.MGR.Id()}
	cluster.MGR.SendMsg(head, funcName, params...)
}

//发送到客户端
func SendToClient(head rpc.RpcHead, packet proto.Message) {
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		head.DestServerType = rpc.SERVICE_GATE
		head.Id = pakcetHead.Id
	}
	cluster.MGR.SendMsg(head, "", proto.MessageName(packet), packet)
}
