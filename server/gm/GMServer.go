package gm

import (
	"database/sql"
	"gonet/actor"
	"gonet/base"
	"gonet/base/ini"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/orm"
	"gonet/network"
	"gonet/rpc"
	"gonet/server/message"
	"log"

	"github.com/golang/protobuf/proto"
)

type (
	ServerMgr struct {
		m_pService   *network.ServerSocket
		m_pCluster   *cluster.Cluster
		m_pActorDB   *sql.DB
		m_Inited     bool
		m_config     ini.Config
		m_Log        base.CLog
		m_SnowFlake  *cluster.Snowflake
		m_PlayerRaft *cluster.PlayerRaft
	}

	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetCluster() *cluster.Cluster
		GetPlayerRaft() *cluster.PlayerRaft
	}

	Config struct {
		common.Server    `yaml:"gm"`
		common.Db        `yaml:"DB"`
		common.Etcd      `yaml:"etcd"`
		common.SnowFlake `yaml:"snowflake"`
		common.Raft      `yaml:"raft"`
		common.Nats      `yaml:"nats"`
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

	//初始化log文件
	this.m_Log.Init("gamemgr")
	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func() {
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tGM Version:\t%s", base.BUILD_NO)
		this.m_Log.Printf("\tGM IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		this.m_Log.Printf("\tDBServer(LAN):\t%s", CONF.Db.Ip)
		this.m_Log.Printf("\tDBName:\t\t%s", CONF.Db.Name)
		this.m_Log.Println("**********************************************************")
	}
	ShowMessage()

	this.m_Log.Println("正在初始化数据库连接...")
	if this.InitDB() {
		this.m_Log.Printf("[%s]数据库连接是失败...", CONF.Db.Name)
		log.Fatalf("[%s]数据库连接是失败...", CONF.Db.Name)
		return false
	}
	this.m_Log.Printf("[%s]数据库初始化成功!", CONF.Db.Name)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	this.m_pService.Init(CONF.Server.Ip, CONF.Server.Port)
	this.m_pService.Start()

	//本身账号集群管理
	this.m_pCluster = new(cluster.Cluster)
	this.m_pCluster.InitCluster(&common.ClusterInfo{Type: rpc.SERVICE_GM, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)}, CONF.Etcd.Endpoints, CONF.Nats.Endpoints)

	this.m_pCluster.BindPacketFunc(actor.MGR.PacketFunc)

	//snowflake
	this.m_SnowFlake = cluster.NewSnowflake(CONF.SnowFlake.Endpoints)

	//playerraft
	this.m_PlayerRaft = cluster.NewPlayerRaft(CONF.Raft.Endpoints)

	SIMPLEMGR.Init()
	return false
}

func (this *ServerMgr) InitDB() bool {
	this.m_pActorDB = orm.OpenDB(CONF.Db)
	err := this.m_pActorDB.Ping()
	return err != nil
}

func (this *ServerMgr) GetDB() *sql.DB {
	return this.m_pActorDB
}

func (this *ServerMgr) GetLog() *base.CLog {
	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket {
	return this.m_pService
}

func (this *ServerMgr) GetCluster() *cluster.Cluster {
	return this.m_pCluster
}

func (this *ServerMgr) GetPlayerRaft() *cluster.PlayerRaft {
	return this.m_PlayerRaft
}

//发送game
func SendToGame(ClusterId uint32, funcName string, params ...interface{}) {
	head := rpc.RpcHead{ClusterId: ClusterId, DestServerType: rpc.SERVICE_GAME, SrcClusterId: SERVER.GetCluster().Id()}
	SERVER.GetCluster().SendMsg(head, funcName, params...)
}

//广播game
func BoardCastToGame(funcName string, params ...interface{}) {
	head := rpc.RpcHead{DestServerType: rpc.SERVICE_GAME, SendType: rpc.SEND_BOARD_CAST, SrcClusterId: SERVER.GetCluster().Id()}
	SERVER.GetCluster().SendMsg(head, funcName, params...)
}

//发送到客户端
func SendToClient(head rpc.RpcHead, packet proto.Message) {
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		head.DestServerType = rpc.SERVICE_GATE
		head.Id = pakcetHead.Id
	}
	SERVER.GetCluster().SendMsg(head, "", proto.MessageName(packet), packet)
}
