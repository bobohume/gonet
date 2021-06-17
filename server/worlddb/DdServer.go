 package worlddb

 import (
	 "database/sql"
	 "gonet/base"
	 "gonet/base/ini"
	 "gonet/common"
	 "gonet/db"
	 "gonet/network"
	 "log"
 )

type(
	ServerMgr struct{
		m_pService     *network.ServerSocket
		m_pWorldClient *network.ClientSocket
		m_pActorDB     *sql.DB
		m_Inited       bool
		m_config       ini.Config
		m_Log          base.CLog
	}

	IServerMgr interface{
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
	}

	Config struct {
		common.Server	`yaml:"db"`
		common.Db	`yaml:"dbDB"`
		common.Etcd		`yaml:"etcd"`
		common.Nats		`yaml:"nats"`
	}
)

var(
	CONF Config
	SERVER ServerMgr
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//初始化log文件
	this.m_Log.Init("world")
	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tDbServerIP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
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

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	this.m_pService.Init(CONF.Server.Ip, CONF.Server.Port)
	this.m_pService.Start()

	var packet EventProcess
	packet.Init()
	this.m_pService.BindPacketFunc(packet.PacketFunc)

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

 func (this *ServerMgr) GetServer() *network.ServerSocket{
	 return this.m_pService
 }

func (this *ServerMgr) GetLog() *base.CLog{
	return &this.m_Log
}