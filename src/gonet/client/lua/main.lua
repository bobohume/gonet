require("isocket")
require("client_socket")
local socket = require("socket")

dofile("isocket.lua")
dofile("client_socket.lua")
dofile("crc.lua")
dofile("packet.lua")

function main()
    client = ClientSocket:new(nil, "127.0.0.1", 31700)
    client:Connect()
     recvt, sendt, status = socket.select({client.m_Conn}, nil, 1)
    while #recvt > 0 do
        local dat, receive_status = client.m_Conn:receive()
        if receive_status ~= "closed" then
            if dat then
                client:ReceivePacket(0, dat)
                print(dat)
                -- recvt, sendt, status = socket.select({client.m_Conn}, nil, 1)
            end
        else
            break
        end
    end
end

main()