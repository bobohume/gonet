package rpc

import (
	"sync"
	"sync/atomic"
)

var (
	RpcSyncSeq     int64
	RpcSyncSeqMap 	sync.Map
)

type(
	RpcSync struct{
		RpcChan chan []interface{}
		Seq int64
	}
)

func CrateRpcSync() *RpcSync{
	req := RpcSync{}
	req.Seq = atomic.AddInt64(&RpcSyncSeq, 1)
	req.RpcChan = make(chan []interface{})
	RpcSyncSeqMap.Store(req.Seq, &req)
	return &req
}

func GetRpcSync(seq int64) *RpcSync{
	req, bOk := RpcSyncSeqMap.Load(seq)
	if bOk{
		RpcSyncSeqMap.Delete(seq)
		return req.(*RpcSync)
	}
	return nil
}

func Sync(seq int64, data []interface{}) bool{
	req := GetRpcSync(seq)
	if req != nil{
		req.RpcChan <- data
		return true
	}
	return false
}

//同步Call
/*
	resut := pPlayer.SyncMsg(rpc.RpcHead{}, "rpctest", gateClusterId, zoneClusterId)

	this.RegisterCall("rpctest", func(ctx context.Context, gateClusterId uint32, zoneClusterId uint32) (error, uint32, uint32){
		return nil ,gateClusterId, zoneClusterId
	})
*/