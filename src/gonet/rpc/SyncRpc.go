package rpc

import (
	"errors"
	"gonet/message"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var (
	RpcSyncSeq     int64
	RpcSyncSeqMap 	sync.Map
	ActorMgr IActor//actor mgr
	ClusterMgr ICluster//cluster mgr
)

type(
	RpcSync struct{
		m_RpcChan chan []byte
		m_Seq int64
	}

	IActor interface {
		SendMsg(head RpcHead, funcName string, params ...interface{})
	}

	ICluster interface {
		IActor
		Id() uint32
		ServiceType() message.SERVICE
	}
)

func Init(clusterMgr ICluster, actorMgr IActor){
	ActorMgr = actorMgr
	ClusterMgr = clusterMgr
}

func SyncMsg(head RpcHead, funcName string, params ...interface{}){
	if head.ActorName == ""{
		head.SrcClusterId = ClusterMgr.Id()
		head.SrcServerType = ClusterMgr.ServiceType()
		ClusterMgr.SendMsg(head, funcName, params...)
	}else{
		ActorMgr.SendMsg(head, funcName, params...)
	}
}

func SyncCall(call interface{}, head RpcHead, funcName string, params ...interface{}) error{
	req := crateRpcSync()
	head.SeqId = req.m_Seq
	SyncMsg(head, funcName, params...)
	// 等待同步回复
	select {
	case v := <-req.m_RpcChan:
		rpcPacket, _ := UnmarshalHead(v)
		funcName := rpcPacket.FuncName
		f := reflect.ValueOf(call)
		k := reflect.TypeOf(call)
		strParams := k.String()
		params := UnmarshalBody(rpcPacket, k)
		if k.NumIn()  != len(params) {
			log.Printf("func [%s] can not call, func params [%s], params [%v]", funcName, strParams, params)
			return errors.New("params no fit")
		}

		if len(params) >= 1{
			bParmasFit := true
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
				//params no fit
				if k.In(i).Kind() != in[i].Kind(){
					bParmasFit = false
				}
			}

			if bParmasFit{
				f.Call(in)
			}else{
				log.Printf("func [%s] params no fit, func params [%s], params [func(%v)]", funcName, strParams, in)
				return errors.New("params no fit")
			}
		}else{
			f.Call(nil)
		}
	case <-time.After(3*time.Second):
		// 清理请求
		getRpcSync(req.m_Seq)
		return errors.New("time out")
	}
	return nil
}

func crateRpcSync() *RpcSync{
	req := RpcSync{}
	req.m_Seq = atomic.AddInt64(&RpcSyncSeq, 1)
	req.m_RpcChan = make(chan []byte)
	RpcSyncSeqMap.Store(req.m_Seq, &req)
	return &req
}

func getRpcSync(seq int64) *RpcSync{
	req, bOk := RpcSyncSeqMap.Load(seq)
	if bOk{
		RpcSyncSeqMap.Delete(seq)
		return req.(*RpcSync)
	}
	return nil
}

func Sync(seq int64, data []byte) bool{
	req := getRpcSync(seq)
	if req != nil{
		req.m_RpcChan <- data
		return true
	}
	return false
}

//同步Call
/*
	this.RegisterCall("test", func(kk int)(int, int) {
		return 1, 2
	})

	rpc.SyncCall(func(kk, jj int){
		fmt.Println(kk, jj)
	}, rpc.RpcHead{DestServerType:message.SERVICE_ACCOUNTSERVER, SendType:message.SEND_BALANCE}, "test", 2)

	rpc.SyncCall(func(kk, jj int){
		fmt.Println(kk, jj)
	}, rpc.RpcHead{ActorName:"chatmgr"}, "test", 2)
*/