package game

import (
	"gonet/actor"
	"gonet/base"
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
		service   *network.ServerSocket
		isInited  bool
		snowFlake *cluster.Snowflake
	}

	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetServer() *network.ServerSocket
	}

	Config struct {
		common.Server    `yaml:"game"`
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
	RdID   int
)

func (s *ServerMgr) Init() bool {
	if s.isInited {
		return true
	}
	//test reload file
	/*file := &common.FileMonitor{}
	file.Init()
	file.AddFile("GONET_SERVER.CFG", func() {base.ReadConf("gonet.yaml", &CONF)})
	file.AddFile(data.SKILL_DATA_NAME, func() {
		data.SKILLDATA.Read()
	})*/

	//初始配置文件
	base.ReadConf("gonet.yaml", &CONF)

	ShowMessage := func() {
		base.LOG.Println("**********************************************************")
		base.LOG.Printf("\tGAME Version:\t%s", base.BUILD_NO)
		base.LOG.Printf("\tGAME IP(LAN):\t%s:%d", CONF.Server.Ip, CONF.Server.Port)
		base.LOG.Printf("\tDBServer(LAN):\t%s", CONF.Db.Ip)
		base.LOG.Printf("\tDBName:\t\t%s", CONF.Db.Name)
		base.LOG.Println("**********************************************************")
	}
	ShowMessage()

	base.LOG.Println("正在初始化数据库连接...")
	if s.InitDB() {
		base.LOG.Printf("[%s]数据库连接是失败...", CONF.Db.Name)
		log.Fatalf("[%s]数据库连接是失败...", CONF.Db.Name)
		return false
	}
	base.LOG.Printf("[%s]数据库初始化成功!", CONF.Db.Name)

	//初始化socket
	s.service = new(network.ServerSocket)
	s.service.Init(CONF.Server.Ip, CONF.Server.Port)
	s.service.Start()

	//snowflake
	s.snowFlake = cluster.NewSnowflake(CONF.SnowFlake.Endpoints)

	//本身game集群管理
	cluster.MGR.InitCluster(&common.ClusterInfo{Type: rpc.SERVICE_GAME, Ip: CONF.Server.Ip, Port: int32(CONF.Server.Port)},
		CONF.Etcd.Endpoints, CONF.Nats.Endpoints, cluster.WithMailBoxEtcd(CONF.Raft.Endpoints), cluster.WithStubMailBoxEtcd(CONF.Raft.Endpoints, &CONF.Stub))

	var packet EventProcess
	packet.Init()
	cluster.MGR.BindPacketFunc(actor.MGR.PacketFunc)
	return false
}

func (s *ServerMgr) InitDB() bool {
	return orm.OpenDB(CONF.Db) != nil
}

func (s *ServerMgr) GetServer() *network.ServerSocket {
	return s.service
}

//发送gamemgr
func SendToGM(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.DestServerType = rpc.SERVICE_GM
	cluster.MGR.SendMsg(head, funcName, params...)
}

//发送给客户端
func SendToClient(clusterId uint32, packet proto.Message) {
	pakcetHead := packet.(message.Packet).GetPacketHead()
	if pakcetHead != nil {
		cluster.MGR.SendMsg(rpc.RpcHead{DestServerType: rpc.SERVICE_GATE, ClusterId: clusterId, Id: pakcetHead.Id}, "", proto.MessageName(packet), packet)
	}
}

//--------------发送给地图----------------------//
func SendToZone(Id int64, ClusterId uint32, funcName string, params ...interface{}) {
	head := rpc.RpcHead{Id: Id, ClusterId: ClusterId, DestServerType: rpc.SERVICE_ZONE, SrcClusterId: cluster.MGR.Id()}
	cluster.MGR.SendMsg(head, funcName, params...)
}
