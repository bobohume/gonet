require("crc")
local protobuf = require "protobuf"
m_PacketCreateMap = {}
m_PacketMap = {}

function RegisterPacket(packet, func)
    name = string.lower(proto.MessageName(packet))
    id = CRC32.hash("test")
    packetFunc = function()
    		packet = proto.Clone(packet)
    		return packet
    end
    m_PacketCreateMap[id] = packetFunc
    m_PacketMap[id] = func
end

function HandlePacket(dat)
    id = bytes_to_int(string.sub(dat, 0, 4))
    buff = string.sub(dat, 4)
    packet, bEx = m_PacketCreateMap[id]
    if bEx then
        proto.Unmarshal(buff, packet)
        m_PacketMap[id](packet)
    end
end


function bytes_to_int(str,endian,signed) -- use length of string to determine 8,16,32,64 bits
    local t={str:byte(1,-1)}
    if endian=="big" then --reverse bytes
        local tt={}
        for k=1,#t do
            tt[#t-k+1]=t[k]
        end
        t=tt
    end
    local n=0
    for k=1,#t do
        n=n+t[k]*2^((k-1)*8)
    end
    if signed then
        n = (n > 2^(#t-1) -1) and (n - 2^#t) or n -- if last bit set, negative.
    end
    return n
end

function int_to_bytes(num,endian,signed)
    if num<0 and not signed then num=-num print"warning, dropping sign from number converting to unsigned" end
    local res={}
    local n = math.ceil(select(2,math.frexp(num))/8) -- number of bytes to be used.
    if signed and num < 0 then
        num = num + 2^n
    end
    for k=n,1,-1 do -- 256 = 2^8 bits per char.
        local mul=2^(8*(k-1))
        res[k]=math.floor(num/mul)
        num=num-res[k]*mul
    end
    assert(num==0)
    if endian == "big" then
        local t={}
        for k=1,n do
            t[k]=res[n-k+1]
        end
        res=t
    end
    return string.char(unpack(res))
end