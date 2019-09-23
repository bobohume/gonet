message 基础类库
     clientsocket
            单点链接
     
     servesocket
            管理ServerSocketClient,ServerSocketClient为具体每个链接对象
            
      websocket
            websocekt管理WebSocketClient,WebSocketClient为具体每个链接对象

     websocket具体修改netgate目录的netgateServer.go
     修改如下
     将m_pService	*network.ServerSocket
     修改为m_pService	*network.WebSocket
     
     将GetServer() *network.ServerSocket
     修改为GetServer() *network.WebSocket
     
     将func (this *ServerMgr) GetServer() *network.ServerSocket{
     修改为func (this *ServerMgr) GetServer() *network.WebSocket{
     
     将this.m_pService = new(network.ServerSocket)
     修改为this.m_pService = new(network.WebSocket)
         
              
在netgateserver里面注释回//websocket这段
