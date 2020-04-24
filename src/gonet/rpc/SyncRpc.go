package rpc

/*import (
	"gonet/message"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var (
	RpcSyncSeq     uint64
	RpcSyncSeqMap sync.Map
	CLUSTER_ACTOR IActor
)

type(
	RpcSync struct{
		m_RpcChan chan []byte
		m_Seq uint64
	}

	IActor interface {
		SendMsg(head RpcHead, funcName string, params ...interface{})
		Id() uint32
		ServiceType() message.SERVICE
	}
)

func SyncCall(call interface{}, head RpcHead, funcName string, params ...interface{}) {
	req := crateRpcSync()
	head.SeqId = req.m_Seq
	head.SrcClusterId = CLUSTER_ACTOR.Id()
	head.SrcServerType = CLUSTER_ACTOR.ServiceType()
	CLUSTER_ACTOR.SendMsg(head, funcName, params...)
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
			return
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
			}
		}else{
			f.Call(nil)
		}
	case <-time.After(3*time.Second):
		// 清理请求
		getRpcSync(req.m_Seq)
		return
	}
}

func crateRpcSync() *RpcSync{
	req := RpcSync{}
	req.m_Seq = atomic.AddUint64(&RpcSyncSeq, 1)
	req.m_RpcChan = make(chan []byte)
	RpcSyncSeqMap.Store(req.m_Seq, &req)
	return &req
}

func getRpcSync(seq uint64) *RpcSync{
	req, bOk := RpcSyncSeqMap.Load(seq)
	if bOk{
		RpcSyncSeqMap.Delete(seq)
		return req.(*RpcSync)
	}
	return nil
}

func Sync(seq uint64, data []byte) bool{
	req := getRpcSync(seq)
	if req != nil{
		req.m_RpcChan <- data
		return true
	}
	return false
}*/