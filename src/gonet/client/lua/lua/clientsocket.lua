local socket = require("socket")
require("isocket")
ClientSocket = Socket:new()

function ClientSocket:new(o, ip, port)
    o = o or Socket:new(o, ip, port)
    setmetatable(o, self)
    self.__index = self
    return o
end

function ClientSocket:Connect()
    if self.m_nState == SSF_CONNECT then
    		return false
    end

    self.m_Conn = assert(socket.connect(self.m_sIp, self.m_nPort))
    self.m_nState = 1
    self.m_Conn:settimeout(0)
    self.m_Conn:setoption("tcp-nodelay", true)
    self.m_Conn:setstats(1024*1024,1024*1024)
    print("连接成功，请输入信息")
    return true
end

function ClientSocket:Send(buf)
     self.m_Conn:send(buf)
end

function ClientSocket:Receive()
    local recvt, sendt, status
    recvt, sendt, status = socket.select({self.m_Conn}, nil, 1)
    while #recvt > 0 do
        local dat, receive_status = self.m_Conn:receive(1)
        if receive_status ~= "closed" then
            if dat then
                self:ReceivePacket(0, dat)
                recvt, sendt, status = socket.select({self.m_Conn}, nil, 1)
            end
        else
            break
        end
    end
end

CLIENTSOCKET = ClientSocket:new(nil, "127.0.0.1", 31700)