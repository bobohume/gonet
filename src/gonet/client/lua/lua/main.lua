require("isocket")
require("clientsocket")
require("account")
local socket = require("socket")


--加载lua文件
dofile("isocket.lua")
dofile("clientsocket.lua")
dofile("crc.lua")
dofile("packet.lua")
dofile("serpent.lua")
dofile("account.lua")

function main()
    CLIENTSOCKET:Connect()
    Account:LoginAccount()
    CLIENTSOCKET:Receive()
end

main()