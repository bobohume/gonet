local socket = require("socket")
local host = "127.0.0.1"
local file = "/"
local sock = assert(socket.connect(host, 31700))

g_test = {}
i = 0
function GetStr()
        g_test[i] = "test"..100
        i = i +1
        print "world1"
        for i1, v1 in pairs(g_test)
        do
            print(i1, v1)
        end
end

GetStr()

