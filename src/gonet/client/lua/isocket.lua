require ("packet")
Socket={
    m_Conn = {},
    m_nPort = 0,
    m_sIp = "",
    m_nState = 0,
    m_bShuttingDown = false,
    m_pInBuffer = ""
}

TCP_END = "💞♡"
TCP_END_LENGTH = #(TCP_END)

function Socket:new(o, ip, port)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    self.m_sIp = ip
    self.m_nPort = port
    self.m_nState = 0
    self.m_pInBuffer = ""
    return o
end

function Socket:Start()
    return true
end

function Socket:Stop()
    self.m_bShuttingDown = true
    return true
end

function Socket:Restart()
	return true
end

function Socket:Connect()
	return true
end

function Socket:Disconnect()
	return true
end

function Socket:OnNetFail()
	self.Stop()
end

function Socket:Send(buf)
	return  0
end

function Socket:ReceivePacket(Id, dat)
	--找包结束
	seekToTcpEnd = function(dat)
		nLen, nEnd = string.find(buff, TCP_END)
		if nLen ~= nil then
			return true, nEnd
		end
		return false, 0
	end

	buff = self.m_pInBuffer .. dat
	self.m_pInBuffer = ""
	nCurSize = 0
::ParsePacekt:: do
	nPacketSize = 0
	buff1 = string.sub(buff, nCurSize+1)
	nBufferSize = #(buff1)
	bFindFlag = false
	bFindFlag, nPacketSize = seekToTcpEnd(buff1)
	print(bFindFlag, nPacketSize, nBufferSize)
	if bFindFlag then
		if nBufferSize == nPacketSize then --完整包
		    print(string.sub(buff1, 0, nPacketSize - TCP_END_LENGTH))
			HandlePacket(string.sub(buff1, 0, nPacketSize - TCP_END_LENGTH))
			nCurSize =  nCurSize + nPacketSize
		elseif (nBufferSize > nPacketSize) then
		    print(string.sub(buff1, 0, nPacketSize - TCP_END_LENGTH))
			HandlePacket(string.sub(buff1, 0, nPacketSize - TCP_END_LENGTH))
			nCurSize =  nCurSize + nPacketSize
			goto ParsePacekt
		end
	elseif nBufferSize < 128 * 1024 then
		self.m_pInBuffer = buff[nCurSize]
	else
		fmt.Println("超出最大包限制，丢弃该包")
	end
end
end
