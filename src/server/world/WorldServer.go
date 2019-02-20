 package world

 import (
	 "base"
	 "database/sql"
	 "db"
	 "github.com/golang/protobuf/proto"
	 "log"
	 "message"
	 "network"
	 "rd"
	 "server/common"
 )

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pAccountClient *network.ClientSocket
		m_pMonitorClient *common.MonitorClient
		m_pServerMgr *ServerSocketManager
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
		GetAccountSocket() *network.ClientSocket
	}

	BitStream base.BitStream
)

var(
	UserNetIP string
	UserNetPort string
	AccountServerIp string
	AccountServerPort string
	DB_Server string
	DB_Name string
	DB_UserId string
	DB_Password string
	Web_Url string
	SERVER ServerMgr
	RdID int
	OpenRedis bool
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
	UserNetIP, UserNetPort 	= this.m_config.Get2("World_LANAddress", ":")
	AccountServerIp, AccountServerPort 	= this.m_config.Get2("Account_LANAddress", ":")
	DB_Server 	= this.m_config.Get("ActorDB_LANIP")
	DB_Name		= this.m_config.Get("ActorDB_Name");
	DB_UserId	= this.m_config.Get("ActorDB_UserId");
	DB_Password	= this.m_config.Get("ActorDB_Password")
	RdID 		= 0//this.m_config.Int("WorkID") / 10
	OpenRedis	= this.m_config.Bool("Redis_Open")
	Web_Url		= this.m_config.Get("World_Url")

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tWorldServerr Version:\t%s",base.BUILD_NO)
		this.m_Log.Printf("\tWorldServerIP(LAN):\t%s:%s", UserNetIP, UserNetPort)
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

	if OpenRedis{
		rd.OpenRedisPool(this.m_config.Get("Redis_Host"), this.m_config.Get("Redis_Pwd"))
	}

	//链接monitor
	this.m_pMonitorClient = new(common.MonitorClient)
	monitorIp, monitroPort := this.m_config.Get2("Monitor_LANAddress", ":")
	this.m_pMonitorClient.Connect(int(message.SERVICE_WORLDSERVER), monitorIp, monitroPort, UserNetIP, UserNetPort)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetMaxReceiveBufferSize(1024)
	this.m_pService.SetMaxSendBufferSize(1024)
	this.m_pService.Start()


	//连接account
	this.m_pAccountClient = new(network.ClientSocket)
	port = base.Int(AccountServerPort)
	this.m_pAccountClient.Init(AccountServerIp, port)
	packet3 := new(AccountProcess)
	packet3.Init(1000)
	this.m_pAccountClient.BindPacketFunc(packet3.PacketFunc)
	this.m_pAccountClient.Start()

	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)

	this.m_pServerMgr = new(ServerSocketManager)
	this.m_pServerMgr.Init(1000)

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

 func (this *ServerMgr) GetLog() *base.CLog{
	 return &this.m_Log
 }

 func (this *ServerMgr) GetServer() *network.ServerSocket{
	 return this.m_pService
 }

 func (this *ServerMgr) GetServerMgr() *ServerSocketManager{
	 return this.m_pServerMgr
 }

 func (this *ServerMgr) GetAccountSocket() *network.ClientSocket{
 	return this.m_pAccountClient
 }

 func SendToClient(socketId int, packet proto.Message){
 	 buff := message.Encode(packet)
 	 nLen := len(buff) + 128
	 pakcetHead := message.GetPakcetHead(packet)
	 if pakcetHead != nil{
		 bitstream := base.NewBitStream(make([]byte, nLen), nLen)
		 bitstream.WriteString(message.GetMessageName(packet))
		 bitstream.WriteInt64(*pakcetHead.Id, base.Bit64)
		 bitstream.WriteBits(len(buff) << 3, buff)
		 SERVER.GetServer().SendByID(socketId, bitstream.GetBuffer())
	 }
 }