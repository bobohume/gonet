message 基础类库
     clientsocket
            单点链接
     
     servesocket
            管理ServerSocketClient,ServerSocketClient为具体每个链接对象
            
      websocket
            websocekt管理WebSocketClient,WebSocketClient为具体每个链接对象

     websocket具体修改net包下的
          websocket模式下，要在net/websocket/server.go 
              func (s Server) serveWebSocket(w http.ResponseWriter, req *http.Request) {
              	rwc, buf, err := w.(http.Hijacker).Hijack()
              	if err != nil {
              		panic("Hijack failed: " + err.Error())
              	}
              	// The server should abort the WebSocket connection if it finds
              	// the client did not send a handshake that matches with protocol
              	// specification.
              	//defer rwc.Close() //注释掉  
              	conn, err := newServerConn(rwc, buf, req, &s.Config, s.Handshake)
              	if err != nil {
              		return
              	}
              	if conn == nil {
              		panic("unexpected nil conn")
              	}
              	s.Handler(conn)
              } 
              
              
          在net/websocket/hybi.go 
              // newHybiConn creates a new WebSocket connection speaking hybi draft protocol.
              func newHybiConn(config *Config, buf *bufio.ReadWriter, rwc io.ReadWriteCloser, request *http.Request) *Conn {
                if buf == nil {
                    br := bufio.NewReader(rwc)
                    bw := bufio.NewWriter(rwc)
                    buf = bufio.NewReadWriter(br, bw)
                }
                ws := &Conn{config: config, request: request, buf: buf, rwc: rwc,
                    frameReaderFactory: hybiFrameReaderFactory{buf.Reader},
                    frameWriterFactory: hybiFrameWriterFactory{
                        buf.Writer, request == nil},
                    PayloadType:        BinaryFrame,//改成二进制
                    defaultCloseStatus: closeStatusNormal}
                ws.frameHandler = &hybiFrameHandler{conn: ws}
                return ws
              }
              
[websocket注视掉如下代码](https://studygolang.com/articles/14842),在netgateserver里面注释回//websocket这段
