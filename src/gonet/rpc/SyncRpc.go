package rpc

import (
	"sync"
)

var (
	RpcSyncSeq     int64
	RpcSyncSeqMap 	sync.Map
)

type(
	RpcSync struct{
		m_RpcChan chan []byte
		m_Seq int64
	}
)

/*func SyncCall(call interface{}, head RpcHead, funcName string, params ...interface{}) error{
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
			in := make([]reflect.Value, len(params))
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			f.Call(in)
		}else{
			log.Printf("func [%s] params at least one context", funcName)
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
}*/

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