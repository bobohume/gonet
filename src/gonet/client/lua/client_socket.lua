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
    print("连接成功，请输入信息")
    return true
end