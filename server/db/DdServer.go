 package db

 import (
	 "database/sql"
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

type(
	ServerMgr struct{
		m_pService     *network.ServerSocket
		m_pActorDB     *sql.DB
		m_Inited       bool
		m_config       ini.Config
		m_Log          base.CLog
		m_pCluster   *cluster.Cluster
		m_PlayerRaft *cluster.PlayerRaft
	}

	IServerMgr interface{
		Init(string) bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetCluster() *cluster.Cluster
		GetPlayerRaft() *cluster.PlayerRaft
	}

	Config struct {
		common.Server	`yaml:"db"`
		common.Db	`yaml:"DB"`
		common.Etcd		`yaml:"etcd"`
		common.Nats		`yaml:"nats"`
		common.Raft      `yaml:"raft"`
	}
)

var(
	CONF Config
	SERVER ServerMgr
)

func (this *ServerMgr) Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("db")
	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tDB Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tDb IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		this.m_Log.Printf("\tActorDBServer(LAN):\t%s", CONF.Db.Ip)
		this.m_Log.Printf("\tActorDBName:\t\t%s", CONF.Db.Name)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	this.m_Log.Println("正在初始化数据库连接...")
	if this.InitDB(){
		this.m_Log.Printf("[%s]数据库连接是失败...", CONF.Db.Name)
		log.Fatalf("[%s]数据库连接是失败...", CONF.Db.Name)
		return false
	}
	this.m_Log.Printf("[%s]数据库初始化成功!", CONF.Db.Name)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	this.m_pService.Init(CONF.Server.Ip, CONF.Server.Port)
	this.m_pService.Start()

	//playerraft
	this.m_PlayerRaft = cluster.NewPlayerRaft(CONF.Raft.Endpoints)

	PLAYERMGR.Init()

	//本身db集群管理
	this.m_pCluster = new(cluster.Cluster)
	this.m_pCluster.InitCluster(&common.ClusterInfo{Type: rpc.SERVICE_DB, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)}, CONF.Etcd.Endpoints, CONF.Nats.Endpoints)
	this.m_pCluster.BindPacketFunc(actor.MGR.PacketFunc)

	return false
}

func (this *ServerMgr)InitDB() bool{
	this.m_pActorDB = orm.OpenDB(CONF.Db)
	err := this.m_pActorDB.Ping()
	return  err != nil
}

func (this *ServerMgr) GetDB() *sql.DB{
	return this.m_pActorDB
}

 func (this *ServerMgr) GetServer() *network.ServerSocket{
	 return this.m_pService
 }

func (this *ServerMgr) GetLog() *base.CLog{
	return &this.m_Log
}

func (this *ServerMgr) GetCluster() *cluster.Cluster{
	return this.m_pCluster
}

 func (this *ServerMgr) GetPlayerRaft() *cluster.PlayerRaft {
	 return this.m_PlayerRaft
 }
