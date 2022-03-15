package cluster

import (
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/common/cluster/et"
	"gonet/common/cluster/etv3"
	"gonet/rpc"
	"strings"

	"github.com/nats-io/nats.go"
)

const (
	ETCD_DIR     = "server/"
	OFFLINE_TIME = etv3.OFFLINE_TIME
)

type (
	Service    et.Service
	Master     et.Master
	Snowflake  et.Snowflake
	PlayerRaft et.PlayerRaft
	//etv2 like redis, save in memroy and store in qucikphoto
	//etv3 like b+ tree, key is key + version so must auto compact delte the old verison
	//if key is lease time out, the old version already in db
)

//注册服务器
func NewService(info *common.ClusterInfo, Endpoints []string) *Service {
	service := &et.Service{}
	service.Init(info, Endpoints)
	return (*Service)(service)
}

//监控服务器
func NewMaster(info common.IClusterInfo, Endpoints []string, pActor actor.IActor) *Master {
	master := &et.Master{}
	master.Init(info, Endpoints, pActor)
	return (*Master)(master)
}

//uuid生成器
func NewSnowflake(Endpoints []string) *Snowflake {
	uuid := &et.Snowflake{}
	uuid.Init(Endpoints)
	return (*Snowflake)(uuid)
}

//注册playerraft
func NewPlayerRaft(Endpoints []string) *PlayerRaft {
	playerRaft := &et.PlayerRaft{}
	playerRaft.Init(Endpoints)
	return (*PlayerRaft)(playerRaft)
}

func (this *PlayerRaft) GetPlayer(Id int64) *rpc.PlayerClusterInfo {
	return (*et.PlayerRaft)(this).GetPlayer(Id)
}

func (this *PlayerRaft) Publish(info *rpc.PlayerClusterInfo) bool {
	return (*et.PlayerRaft)(this).Publish(info)
}

func (this *PlayerRaft) Lease(leaseId int64) error {
	return (*et.PlayerRaft)(this).Lease(leaseId)
}

func getChannel(clusterInfo common.ClusterInfo) string {
	return fmt.Sprintf("%s/%s/%d", ETCD_DIR, clusterInfo.String(), clusterInfo.Id())
}

func getTopicChannel(clusterInfo common.ClusterInfo) string {
	return fmt.Sprintf("%s/%s", ETCD_DIR, clusterInfo.String())
}

func getCallChannel(clusterInfo common.ClusterInfo) string {
	return fmt.Sprintf("%s/%s/call/%d", ETCD_DIR, clusterInfo.String(), clusterInfo.Id())
}

func getRpcChannel(head rpc.RpcHead) string {
	return fmt.Sprintf("%s/%s/%d", ETCD_DIR, strings.ToLower(head.DestServerType.String()), head.ClusterId)
}

func getRpcTopicChannel(head rpc.RpcHead) string {
	return fmt.Sprintf("%s/%s", ETCD_DIR, strings.ToLower(head.DestServerType.String()))
}

func getRpcCallChannel(head rpc.RpcHead) string {
	return fmt.Sprintf("%s/%s/call/%d", ETCD_DIR, strings.ToLower(head.DestServerType.String()), head.ClusterId)
}

func setupNatsConn(connectString string, appDieChan chan bool, options ...nats.Option) (*nats.Conn, error) {
	natsOptions := append(
		options,
		nats.DisconnectHandler(func(_ *nats.Conn) {
			base.GLOG.Println("disconnected from nats!")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			base.GLOG.Printf("reconnected to nats server %s with address %s in cluster %s!", nc.ConnectedServerId(), nc.ConnectedAddr(), nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			err := nc.LastError()
			if err == nil {
				base.GLOG.Println("nats connection closed with no error.")
				return
			}

			base.GLOG.Println("nats connection closed. reason: %q", nc.LastError())
			if appDieChan != nil {
				appDieChan <- true
			}
		}),
	)

	nc, err := nats.Connect(connectString, natsOptions...)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
