 package worlddb

 import (
	 "database/sql"
	 "gonet/base"
	 "gonet/db"
	 "gonet/network"
	 "log"
 )

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pWorldClient *network.ClientSocket
		m_pActorDB *sql.DB
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
	}

	IServerMgr interface{
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetWorldSocket() *network.ClientSocket
	}

	BitStream base.BitStream
)

var(
	UserNetIP string
	UserNetPort string
	WorldServerIp string
	WorldServerPort string
	DB_Server string
	DB_Name string
	DB_UserId string
	DB_Password string
	SERVER ServerMgr
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//test reload file
	/*file := &common.FileMonitor{}
	file.Init(1000)
	file.AddFile("SXZ_SERVER.CFG", func() {this.m_config.Read("SXZ_SERVER.CFG")})
	file.AddFile(data.SKILL_DATA_NAME, func() {
		data.SKILLDATA.Read()
	})*/

	//初始化log文件
	this.m_Log.Init("world")
	//初始ini配置文件
	this.m_config.Read("SXZ_SERVER.CFG")
	UserNetIP, UserNetPort 	= this.m_config.Get2("Dd_LANAddress", ":")
	WorldServerIp, WorldServerPort 	= this.m_config.Get2("World_LANAddress", ":")
	DB_Server 	= this.m_config.Get("ActorDB_LANIP")
	DB_Name		= this.m_config.Get("ActorDB_Name");
	DB_UserId	= this.m_config.Get("ActorDB_UserId");
	DB_Password	= this.m_config.Get("ActorDB_Password")

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tServer Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tDbServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
		this.m_Log.Printf("\tActorDBServer(LAN):\t%s", DB_Server)
		this.m_Log.Printf("\tActorDBName:\t\t%s", DB_Name)
		this.m_Log.Println("**********************************************************");
	}
	ShowMessage()

	this.m_Log.Println("正在初始化数据库连接...")
	if (this.InitDB()){
		this.m_Log.Printf("[%s]数据库连接是失败...", DB_Name)
		log.Fatalf("[%s]数据库连接是失败...", DB_Name)
		return false
	}
	this.m_Log.Printf("[%s]数据库初始化成功!", DB_Name)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetMaxReceiveBufferSize(1024)
	this.m_pService.SetMaxSendBufferSize(1024)
	this.m_pService.Start()

	//连接world
	this.m_pWorldClient = new(network.ClientSocket)
	port = base.Int(WorldServerPort)
	this.m_pWorldClient.Init(WorldServerIp, port)
	packet3 := new(WorldProcess)
	packet3.Init(1000)
	this.m_pWorldClient.BindPacketFunc(packet3.PacketFunc)
	this.m_pWorldClient.Start()

	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)

	return  false
}

func (this *ServerMgr)InitDB() bool{
	this.m_pActorDB = db.OpenDB(DB_Server, DB_UserId, DB_Password, DB_Name)
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

func (this *ServerMgr) GetWorldSocket() *network.ClientSocket{
	return this.m_pWorldClient
}