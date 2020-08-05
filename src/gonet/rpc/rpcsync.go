package rpc

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	RpcSyncSeq     int64
	RpcSyncSeqMap 	sync.Map
)

type(
	RpcSync struct{
		RpcChan chan RetInfo
		Seq int64
	}

	RetInfo struct {
		Err error
		Ret []interface{}
	}
)

const(
	MAX_RPC_TIMEOUT = 3*time.Second
)

func CrateRpcSync() *RpcSync{
	req := RpcSync{}
	req.Seq = atomic.AddInt64(&RpcSyncSeq, 1)
	req.RpcChan = make(chan RetInfo)
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

func Sync(seq int64, ret RetInfo) bool{
	req := GetRpcSync(seq)
	if req != nil{
		req.RpcChan <- ret
		return true
	}
	return false
}

//同步Call
/*
	resut := pPlayer.SyncMsg(rpc.RpcHead{}, "rpctest", gateClusterId, zoneClusterId)

	this.RegisterCall("rpctest", func(ctx context.Context, gateClusterId uint32, zoneClusterId uint32) rpc.RetInfo{
		return rpc.RetInfo{nil, []interface{}{gateClusterId, zoneClusterId}}
	})
*/