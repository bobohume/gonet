//IPacket.js
var CRC32 = require('crc-32'); 
var Network = require("./network");
var BitStream = require("./BitStream");
var messagepb= require("./pb/message");
var clientpb= require("./pb/client");
STX = 0X27;
CKX = 0x72;

SERVICE_NONE          = 0;
SERVICE_CLIENT        = 1;
SERVICE_GATESERVER    = 2;
SERVICE_ACCOUNTSERVER = 3;
SERVICE_WORLDSERVER   = 4;
SERVICE_MONITORSERVER = 5;

Uint8Array.prototype.indexOfMulti = function(searchElements, fromIndex) {
    fromIndex = fromIndex || 0;

    var index = Array.prototype.indexOf.call(this, searchElements[0], fromIndex);
    if(searchElements.length === 1 || index === -1) {
        // Not found or no other elements to check
        return index;
    }

    for(var i = index, j = 0; j < searchElements.length && i < this.length; i++, j++) {
        if(this[i] !== searchElements[j]) {
            return this.indexOfMulti(searchElements, index + 1);
        }
    }

    return (i === index + searchElements.length) ? index : -1;
};

var Packet = (function(){
    //var Packet = {
        //整形转换成字节
        IntToBytes =  function(v){
            var b = new Uint8Array(4);
            b[0] = v
            b[1] = v >> 8
            b[2] = v >> 16
            b[3] = v >> 24
            return b
        },

        //字节转换成整形
        BytesToInt =   function(b){
            return b[0] | b[1]<<8 | b[2]<<16 | b[3]<<24;
        },


        //创建包头
        BuildPacketHead =  function(Id, DestServerType=0){
            var pHead =  new messagepb.message.Ipacket.create();
            pHead.Ckx = CKX;
            pHead.Stx = STX;
            pHead.Id = Id;
            pHead.DestServerType = DestServerType;
            return pHead
        },

        /*Uint8ArrayToString = function(fileData){
              var dataString = "";
              for (var i = 0; i < fileData.length; i++) {
                dataString += String.fromCharCode(fileData[i]);
              }
             
              return dataString
        },

        stringToUint8Array = function(str){
              var arr = [];
              for (var i = 0, j = str.length; i < j; ++i) {
                arr.push(str.charCodeAt(i));
              }
             
              var tmpUint8Array = new Uint8Array(arr);
              return tmpUint8Array
        },*/

        //包处理回调函数
        RegisterPacket =  function(packetName, func){
            var name = packetName.toLowerCase();
            var id = CRC32.str(name);
            if(m_callbacks[id] != null)
            {
                console.log("消息[%s]重复注册!!!!!!!!!!!!!",packetName);
                return;
            }

            m_callbacks[id] = func;
        },

        //--处理包函数
        HandlePacket  =  function(dat){
            var id = BytesToInt(dat.slice(0,4));
            var buf = dat.slice(4);
            var packetcreator = m_packets[id];
            console.log(id, packetcreator, m_packets, packet, m_callbacks);
            var packet = packetcreator().decode(buf);
            m_callbacks[id](packet);
            return true;
        },


        //--发送包函数
        SendPacket = function(name, dat){
            var id = CRC32.str(name.toLowerCase());
            var packetName = "message." + name;
            if(packetName != null)
            {
                var packetcreator = m_packets[id];
                var crc =IntToBytes(id);
                var packetcreator = m_packets[id];
                var packet = packetcreator().encode(dat).finish();
                var buff = new Uint8Array(crc.length + packet.length + TCP_END_LENGTH);
                buff.set(crc, 0);
                buff.set(packet, crc.length);
                buff.set(TCP_END, crc.length+packet.length);
                //var buff = Uint8ArrayToString(IntToBytes(id)) + dat + TCP_END;
                Network.getInstance().send(buff);
            }
            return false;
        },

        //解析包
        ReceivePacket =  function(dat){
            seekToTcpEnd = function(dat){
                var nCurLen = dat.indexOfMulti(TCP_END);
                if(nCurLen != -1){
                    return nCurLen + TCP_END_LENGTH;
                }
                return 0;
            }

            var buff = new Uint8Array(m_pInBuffer.length + dat.length);
            buff.set(m_pInBuffer, 0);
            buff.set(dat, m_pInBuffer.length);
            m_pInBuffer = new Uint8Array(0);
            var nCurSize = 0
            ParsePacekt = function(){
                var nPacketSize = 0
                var buff1 = buff.slice(nCurSize)
                var nBufferSize = buff1.length;
                nPacketSize = seekToTcpEnd(buff1)
                //console.log(nPacketSize, nBufferSize, TCP_END_LENGTH)
                if (nPacketSize != 0){
                    if(nBufferSize == nPacketSize){ //完整包
                        //print(string.sub(buff1, 0, nPacketSize - TCP_END_LENGTH))
                        HandlePacket(buff1.slice(0, nPacketSize - TCP_END_LENGTH))
                        nCurSize =  nCurSize + nPacketSize
                    }
                    else if (nBufferSize > nPacketSize){
                        //print(string.sub(buff1, 0, nPacketSize - TCP_END_LENGTH))
                        HandlePacket(buff1.slice(0, nPacketSize - TCP_END_LENGTH))
                        nCurSize =  nCurSize + nPacketSize
                        ParsePacekt();
                    }
                    else if(nBufferSize < 128 * 1024){
                        self.m_pInBuffer = buff.slice(nCurSize)
                    }
                    else{
                        console.log("超出最大包限制，丢弃该包")
                    }
                }
            },

            ParsePacekt();
        };
    //};
})();
//var TCP_END  = "💞♡";
var TCP_END = new Uint8Array(7);
TCP_END.set( [240,159,146,158,226,153,161],0);
var TCP_END_LENGTH  = TCP_END.length;
var m_pInBuffer = new Uint8Array();
var m_callbacks = new Array();
var m_packets = new Array();

//--protobuf-- 自动化解析
function RegisterPacketCreator(name, func){
    var id = CRC32.str(name.toLowerCase());
    m_packets[id] = func;
};

RegisterPacketCreator("W_C_SelectPlayerResponse", function(){
    return clientpb.message.W_C_SelectPlayerResponse;
})
RegisterPacketCreator("C_A_LoginRequest", function(){
    return clientpb.message.C_A_LoginRequest;
})

Network.getInstance().host = "ws://192.168.84.62:31700/ws";
Network.getInstance().initNetwork();
module.exports = Packet;