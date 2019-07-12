//Network.js
var BitStream = require("./BitStream");
var ws = require('ws')
var messagepb= require("./pb/message");
var clientpb= require("./pb/client");
var Packet = require("./Packet");

var Network  = (function(){
	var instance = null;
	var SERVER_VERSION = 102008000;
	var MAX_PACKET_SIZE	= 32*1024;
	var	Parent=null;
	
	function getNetworkInstance (){
		var networkInstance = {
			socket:null,
			isInit:false,
			initNetwork:function(){
				console.log('Network initSocket...');
				//this.host = "ws://192.168.1.122:21001";;
				//this.testhost = "ws://echo.websocket.org"
				this.socket = new ws(this.host, {
  						origin: 'http://localhost/'
				});

				this.socket.binaryType = "arraybuffer";
				this.socket.onopen = function(evt){
					console.log('Network onopen...');
					{			
						this.isInit = true;
					}
					LoginAccount();
				};

				this.socket.onmessage = function(evt){
					var aa = new Uint8Array(evt.data);
					Packet:ReceivePacket(aa);
				};
				
				this.socket.onerror = function(evt){
					console.log('Network onerror...');
				};
				
				this.socket.onclose = function(evt){
					console.log('Network onclose...');
					this.isInit = false;
				};

				this.socket.addEventListener("error", function(event) {
					console.log(event)
				})
			},
			//发送消息
			send:function(data){
				if (Parent.isInit){
					console.log('Network is not inited...');
				}else if(Parent.socket.readyState == ws.OPEN){
					Parent.socket.send(data);
				}else{
					console.log('Network WebSocket readState:'+Parent.socket.readyState);
				}
			},
			close:function(){
				if (this.socket){
					console.log("Network close...");
					this.socket.close();
					this.socket = null;
				}
			},
		};
		return networkInstance;
	};


	return {
		getInstance:function(){
			if(instance === null){
				instance = getNetworkInstance();
				Parent = instance;
			}
			return instance;
		}
	};
})();


function LoginAccount(){
	var id = 10003;
	var AccountName = "test10003";
	var packet1 =  clientpb.message.C_A_LoginRequest.create();
	packet1.PacketHead = BuildPacketHead(0, SERVICE_ACCOUNTSERVER);
	packet1.AccountName = AccountName;
	packet1.BuildNo = BUILD_NO;
	packet1.SocketId = 0;
	SendPacket("C_A_LoginRequest", packet1)
};

module.exports=Network;