package world

import (
	"database/sql"
	"github.com/golang/protobuf/proto"
	"gonet/base"
	"gonet/db"
	"gonet/message"
	"gonet/network"
	"gonet/rd"
	"gonet/server/common/cluster"
	"log"
)

type(
	ServerMgr struct{
		m_pService	*network.ServerSocket
		m_pServerMgr *ServerSocketManager
		m_pActorDB *sql.DB
		m_Inited bool
		m_config base.Config
		m_Log	base.CLog
		m_Cluster *cluster.Service
		m_AccountCluster *cluster.Cluster
		m_SnowFlake *cluster.Snowflake
	}

	IServerMgr interface{
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetAccountCluster() *cluster.Cluster
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
	//Web_Url string
	SERVER ServerMgr
	RdID int
	OpenRedis bool
	EtcdEndpoints []string
)

func (this *ServerMgr)Init() bool{
	if(this.m_Inited){
		return true
	}

	//test reload file
	/*file := &common.FileMonitor{}
	file.Init(1000)
	file.AddFile("GONET_SERVER.CFG", func() {this.m_config.Read("GONET_SERVER.CFG")})
	file.AddFile(data.SKILL_DATA_NAME, func() {
		data.SKILLDATA.Read()
	})*/

	//初始化log文件
	this.m_Log.Init("world")
	//初始ini配置文件
	this.m_config.Read("GONET_SERVER.CFG")
	EtcdEndpoints = this.m_config.Get5("Etcd_Cluster", ",")
	UserNetIP, UserNetPort 	= this.m_config.Get2("World_LANAddress", ":")
	AccountServerIp, AccountServerPort 	= this.m_config.Get2("Account_LANAddress", ":")
	DB_Server 	= this.m_config.Get3("WorldDB", "DB_LANIP")
	DB_Name		= this.m_config.Get3("WorldDB","DB_Name");
	DB_UserId	= this.m_config.Get3("WorldDB", "DB_UserId");
	DB_Password	= this.m_config.Get3("WorldDB", "DB_Password")
	RdID 		= 0//this.m_config.Int("WorkID") / 10
	OpenRedis	= this.m_config.Bool("Redis_Open")
	//Web_Url		= this.m_config.Get("World_Url")

	ShowMessage := func(){
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tWorldServer Version:\t%s",base.BUILD_NO)
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

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(UserNetPort)
	this.m_pService.Init(UserNetIP, port)
	this.m_pService.SetMaxReceiveBufferSize(1024)
	this.m_pService.SetMaxSendBufferSize(1024)
	this.m_pService.Start()

	//snowflake
	this.m_SnowFlake = cluster.NewSnowflake(UserNetIP, base.Int(UserNetPort), this.m_config.Get5("Etcd_SnowFlake_Cluster", ","))

	//注册到集群
	this.m_Cluster = cluster.NewService(int(message.SERVICE_WORLDSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)

	//账号服务器集群
	this.m_AccountCluster = new(cluster.Cluster)
	this.m_AccountCluster.Init(1000, int(message.SERVICE_ACCOUNTSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	this.m_AccountCluster.BindPacket(&AccountProcess{})

	this.m_pServerMgr = new(ServerSocketManager)
	this.m_pServerMgr.Init(1000)

	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.BindPacketFunc(this.m_pServerMgr.PacketFunc)
	return  false
}

func (this *ServerMgr)InitDB() bool{
	this.m_pActorDB = db.OpenDB(DB_Server, DB_UserId, DB_Password, DB_Name)
	err := this.m_pActorDB.Ping()
	this.m_pActorDB.SetMaxOpenConns(base.Int(this.m_config.Get3("WorldDB", "DB_MaxOpenConns")))
	this.m_pActorDB.SetMaxIdleConns(base.Int(this.m_config.Get3("WorldDB", "DB_MaxIdleConns")))
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

func (this *ServerMgr) GetAccountCluster() *cluster.Cluster{
 	return this.m_AccountCluster
}

//发送给客户端
func SendToClient(socketId int, packet proto.Message){
	buff := message.Encode(packet)
	nLen := len(buff) + 128
	pakcetHead := message.GetPakcetHead(packet)
	if pakcetHead != nil {
		bitstream := base.NewBitStream(make([]byte, nLen), nLen)
		bitstream.WriteString(message.GetMessageName(packet))
		bitstream.WriteInt64(pakcetHead.Id, base.Bit64)
		bitstream.WriteBits(len(buff)<<3, buff)
		SERVER.GetServer().SendById(socketId, bitstream.GetBuffer())
	}
}