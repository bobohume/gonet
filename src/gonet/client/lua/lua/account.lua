require("packet")
Account={
    AccountId = 0,
    PlayerId = 0,
    AccountName = ""
}

function Account:LoginAccount()
	local id = 0
	self.AccountName = "test"..id
	packet1 = {PacketHead=BuildPacketHead(0, SERVICE_ACCOUNTSERVER),
		AccountName=self.AccountName, BuildNo=BUILD_NO, SocketId=0}
	SendPacket("C_A_LoginRequest",packet1)
end

--注册消息
RegisterPacket("C_A_LoginRequest", nil)
--登录回调
RegisterPacket("A_C_LoginRequest", function(packet)
    if packet.Error == ACCOUNT_NOEXIST then
    end
    print(require "serpent".block(packet))
end)

