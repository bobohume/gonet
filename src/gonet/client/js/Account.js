//Account.js
//账号一些方法
var BitStream = require("./BitStream");
var Packet = require("./Packet");
var messagepb= require("./pb/message");
var clientpb= require("./pb/client");

BUILD_NO = "1,5,1,1";

(function () {

})();


Packet:RegisterPacket("W_C_SelectPlayerResponse", function(packet){
		console.log(packet = "TEST111");
	});

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