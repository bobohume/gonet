package gate

import (
	"context"
	"gonet/actor"
	"gonet/base/cluster"
	"gonet/network"
	"gonet/rpc"
	"time"
)

const (
	PACKET_STATS_INTERVAL = 60
)

type (
	PacketStat struct {
		Times [PACKET_STATS_INTERVAL]uint16
	}

	PacketStats struct {
		HeartTime int64
		IsLimit   bool
		Packets   map[string]*PacketStat
	}

	StatsPrcoess struct {
		actor.Actor
		statsMap map[uint32]*PacketStats
		statTime int
	}
)

func (p *PacketStats) Init() {
	p.Packets = map[string]*PacketStat{}
}

func (p *PacketStat) Count() uint32 {
	times := uint32(0)
	for i := 0; i < PACKET_STATS_INTERVAL; i++ {
		times += uint32(p.Times[i])
	}
	return times
}

func (p *PacketStats) Clear(index int) {
	for _, v := range p.Packets {
		for i := 0; i < PACKET_STATS_INTERVAL; i++ {
			v.Times[index] = 0
		}
	}
}

func (s *StatsPrcoess) Init() {
	s.Actor.Init()
	s.statsMap = map[uint32]*PacketStats{}
	actor.MGR.RegisterActor(s)
	s.RegisterTimer(time.Second, s.Update)
	s.Actor.Start()
}

func (s *StatsPrcoess) Update() {
	s.statTime++
	s.statTime = s.statTime % PACKET_STATS_INTERVAL
	for _, v := range s.statsMap {
		v.Clear(s.statTime)
	}
	deleteMaps := map[uint32]struct{}{}
	for k, v := range s.statsMap {
		if v.HeartTime < time.Now().Unix() {
			deleteMaps[k] = struct{}{}
		}
	}
	for k, _ := range deleteMaps {
		cluster.MGR.SendMsg(rpc.RpcHead{SendType: rpc.SEND_LOCAL}, "PlayerMgr.DEL_ACCOUNT", k)
		cluster.MGR.SendMsg(rpc.RpcHead{SendType: rpc.SEND_LOCAL}, "UserPrcoess.DISCONNECT", k)
		delete(s.statsMap, k)
	}
}

func (s *StatsPrcoess) Stats(ctx context.Context, socketid uint32, packetName string, limit uint32) {
	v, isOK := s.statsMap[socketid]
	if !isOK {
		ps := &PacketStats{HeartTime: time.Now().Unix() + network.HEART_TIME_OUT}
		ps.Init()
		s.statsMap[socketid] = ps
		v = ps
	}

	v1, isEx := v.Packets[packetName]
	if !isEx {
		ps := &PacketStat{}
		v.Packets[packetName] = ps
		v1 = ps
	}
	v1.Times[s.statTime]++
	v.HeartTime = time.Now().Unix() + network.HEART_TIME_OUT
	if limit > 0 && v1.Count() > limit {
		cluster.MGR.SendMsg(rpc.RpcHead{SendType: rpc.SEND_LOCAL}, "PlayerMgr.DEL_ACCOUNT", socketid)
		cluster.MGR.SendMsg(rpc.RpcHead{SendType: rpc.SEND_LOCAL}, "UserPrcoess.DISCONNECT", socketid)
		delete(s.statsMap, socketid)
		return
	}
}
