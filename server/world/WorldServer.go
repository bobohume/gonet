package world

import (
	"database/sql"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/common"
	"gonet/common/cluster"
	"gonet/db"
	"gonet/network"
	"gonet/rd"
	"gonet/rpc"
	"gonet/server/message"
	"log"
)

type(
	ServerMgr struct{
		m_pService  *network.ServerSocket
		m_pCluster  *cluster.Cluster
		m_pActorDB  *sql.DB
		m_Inited    bool
		m_Log       base.CLog
		m_SnowFlake *cluster.Snowflake
	}

	IServerMgr interface{
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetCluster() *cluster.Cluster
	}

	Config struct {
		common.Server	`yaml:"world"`
		common.Db	`yaml:"worldDB"`
		common.Redis	`yaml:"redis"`
		common.Etcd		`yaml:"etcd"`
		common.SnowFlake `yaml:"snowflake"`
		common.Nats		`yaml:"nats"`
	}
)

var(
	CONF Config
	SERVER ServerMgr
	RdID int
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//test reload file
	/*file := &common.FileMonitor{}
	file.Init()
	file.AddFile("GONET_SERVER.CFG", func() {base.ReadConf("gonet.yaml", &CONF)})
	file.AddFile(data.SKILL_DATA_NAME, func() {
		data.SKILLDATA.Read()
	})*/

	//初始化log文件
	this.m_Log.Init("world")
	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tWorldServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tWorldServerIP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		this.m_Log.Printf("\tActorDBServer(LAN):\t%s", CONF.Db.Ip)
		this.m_Log.Printf("\tActorDBName:\t\t%s", CONF.Db.Name)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	this.m_Log.Println("正在初始化数据库连接...")
	if (this.InitDB()){
		this.m_Log.Printf("[%s]数据库连接是失败...", CONF.Db.Name)
		log.Fatalf("[%s]数据库连接是失败...", CONF.Db.Name)
		return false
	}
	this.m_Log.Printf("[%s]数据库初始化成功!", CONF.Db.Name)

	if CONF.Redis.OpenFlag{
		rd.OpenRedisPool(CONF.Redis.Ip, CONF.Redis.Password)
	}

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	this.m_pService.Init(CONF.Server.Ip, CONF.Server.Port)
	this.m_pService.Start()

	//snowflake
	this.m_SnowFlake = cluster.NewSnowflake(CONF.SnowFlake.Endpoints)

	//本身world集群管理
	this.m_pCluster = new(cluster.Cluster)
	this.m_pCluster.Init(&common.ClusterInfo{Type: rpc.SERVICE_WORLDSERVER, Ip:CONF.Server.Ip, Port:int32(CONF.Server.Port)}, CONF.Etcd.Endpoints, CONF.Nats.Endpoints)

	var packet EventProcess
	packet.Init()
	this.m_pCluster.BindPacketFunc(packet.PacketFunc)
	return  false
}

func (this *ServerMgr)InitDB() bool{
	this.m_pActorDB = db.OpenDB(CONF.Db)
	err := this.m_pActorDB.Ping()
	return  err != nil
}

func (this *ServerMgr) GetDB() *sql.DB{
	return this.m_pActorDB
}

func (this *ServerMgr) GetLog() *base.CLog{
	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket{
 	return this.m_pService
}

func (this *ServerMgr) GetCluster() *cluster.Cluster{
	return this.m_pCluster
}

//发送account
func SendToAccount(funcName string, params  ...interface{}){
	head := rpc.RpcHead{DestServerType:rpc.SERVICE_ACCOUNTSERVER, SendType:rpc.SEND_BALANCE, SrcClusterId:SERVER.GetCluster().Id()}
	SERVER.GetCluster().SendMsg(head,  funcName, params...)
}

//发送给客户端
func SendToClient(clusterId uint32, packet proto.Message){
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		SERVER.GetCluster().SendMsg(rpc.RpcHead{DestServerType:rpc.SERVICE_GATESERVER, ClusterId:clusterId, Id:pakcetHead.Id}, "", proto.MessageName(packet), packet)
	}
}

//--------------发送给地图----------------------//
func SendToZone(Id int64, ClusterId uint32, funcName string, params  ...interface{}){
	head := rpc.RpcHead{Id:Id, ClusterId:ClusterId, DestServerType:rpc.SERVICE_ZONESERVER, SrcClusterId:SERVER.GetCluster().Id()}
	SERVER.GetCluster().SendMsg(head, funcName, params...)
}